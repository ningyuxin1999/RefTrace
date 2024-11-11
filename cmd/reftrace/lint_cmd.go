package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"reft-go/nf"
	"reft-go/parser"
	"strconv"

	"github.com/antlr4-go/antlr/v4"
	"github.com/spf13/cobra"
	"go.starlark.net/starlark"
)

var (
	rulesFile string
	dir       string
	ruleToRun string
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint a directory using rules from a rules.py file",
	Run:   runLint,
}

func init() {
	rootCmd.AddCommand(lintCmd)

	lintCmd.Flags().StringVarP(&rulesFile, "rules", "r", "rules.py", "Path to the rules file")
	lintCmd.Flags().StringVarP(&dir, "directory", "d", ".", "Directory to lint")
	lintCmd.Flags().StringVarP(&ruleToRun, "name", "n", "", "Name of a single rule to run")
}

type StarlarkParamInfo struct {
	paramInfo nf.ParamInfo
}

func (p *StarlarkParamInfo) String() string {
	return fmt.Sprintf("Param(name=%s, line=%d)", p.paramInfo.Name, p.paramInfo.LineNumber)
}

func (p *StarlarkParamInfo) Type() string {
	return "Param"
}

func (p *StarlarkParamInfo) Freeze() {}

func (p *StarlarkParamInfo) Truth() starlark.Bool {
	return starlark.Bool(true)
}

func (p *StarlarkParamInfo) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(p.paramInfo.Name))
	h.Write([]byte(strconv.Itoa(p.paramInfo.LineNumber)))
	return h.Sum32(), nil
}

func (p *StarlarkParamInfo) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(p.paramInfo.Name), nil
	case "line":
		return starlark.MakeInt(p.paramInfo.LineNumber), nil
	default:
		return nil, nil
	}
}

func (p *StarlarkParamInfo) AttrNames() []string {
	return []string{"name", "line"}
}

func NewStarlarkParamInfo(p nf.ParamInfo) *StarlarkParamInfo {
	return &StarlarkParamInfo{paramInfo: p}
}

func parse(filePath string) ([]nf.ParamInfo, []nf.IncludeStatement) {
	//debug.SetGCPercent(-1)
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		log.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := parser.NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	//tokens := lexer.GetAllTokens()
	//tokenStream := NewPreloadedTokenStream(tokens, lexer)
	stream.Fill()
	groovyParser := parser.NewGroovyParser(stream)
	//parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

	// Parse the file
	tree := groovyParser.CompilationUnit()
	builder := parser.NewASTBuilder(filePath)
	ast := builder.Visit(tree).(*parser.ModuleNode)
	paramVisitor := nf.NewParamVisitor()
	paramVisitor.VisitBlockStatement(ast.StatementBlock)
	params := paramVisitor.GetSortedParams()

	includeVisitor := nf.NewIncludeVisitor()
	includeVisitor.VisitBlockStatement(ast.StatementBlock)
	includes := includeVisitor.GetSortedIncludes()

	return params, includes
}

func runLint(cmd *cobra.Command, args []string) {
	config := nf.LintConfig{
		RulesFile: rulesFile,
		Directory: dir,
		RuleToRun: ruleToRun,
	}
	err := nf.RunLintWithConfig(config, os.Stdout)
	if err != nil {
		log.Fatalf("Linting failed: %v", err)
	}
}
