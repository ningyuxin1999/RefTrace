package nf

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	pb "reft-go/nf/proto"
)

type ResolvedInclude struct {
	ModulePath string
	Includes   []string
}

type UnresolvedInclude struct {
	ModulePath string
	Includes   []string
}

func (inc *ResolvedInclude) ToProto() *pb.ResolvedInclude {
	return &pb.ResolvedInclude{
		ModulePath: inc.ModulePath,
		Includes:   inc.Includes,
	}
}

func (inc *UnresolvedInclude) ToProto() *pb.UnresolvedInclude {
	return &pb.UnresolvedInclude{
		ModulePath: inc.ModulePath,
		Includes:   inc.Includes,
	}
}

func makeAbs(modulePath, includePath string) string {
	// If the include path is absolute (doesn't start with . or ..), return as is
	if !strings.HasPrefix(includePath, ".") && !strings.HasPrefix(includePath, "..") {
		return includePath
	}

	// Get the directory containing the module
	moduleDir := filepath.Dir(modulePath)

	// Join the module directory with the relative include path and clean it
	return filepath.Clean(filepath.Join(moduleDir, includePath))
}

func canonicalize(modulePath, includePath string) string {
	abs := makeAbs(modulePath, includePath)

	// If path ends with .nf, return as is
	if strings.HasSuffix(abs, ".nf") {
		return abs
	}

	// If path ends with "main", append ".nf"
	if strings.HasSuffix(abs, "main") {
		return abs + ".nf"
	}

	// Special case for workflows directory - just append .nf
	if filepath.Base(filepath.Dir(abs)) == "workflows" {
		return abs + ".nf"
	}

	// Otherwise append "/main.nf"
	return abs + "/main.nf"
}

func ResolveIncludes(modules []*Module) ([]*ResolvedInclude, []*UnresolvedInclude) {
	moduleNames := make(map[string]struct{})
	for _, module := range modules {
		moduleNames[module.Path] = struct{}{}
	}
	resolvedIncludes := []*ResolvedInclude{}
	unresolvedIncludes := []*UnresolvedInclude{}
	for _, module := range modules {
		resolvedSet := make(map[string]struct{})
		unresolvedSet := make(map[string]struct{})
		for _, include := range module.Includes {
			canonicalPath := canonicalize(module.Path, include.ModulePath)
			if _, ok := moduleNames[canonicalPath]; ok {
				resolvedSet[canonicalPath] = struct{}{}
			} else {
				unresolvedSet[canonicalPath] = struct{}{}
			}
		}
		// Convert sets to slices
		resolved := make([]string, 0, len(resolvedSet))
		for path := range resolvedSet {
			resolved = append(resolved, path)
		}

		unresolved := make([]string, 0, len(unresolvedSet))
		for path := range unresolvedSet {
			unresolved = append(unresolved, path)
		}
		if len(resolved) > 0 {
			resolvedIncludes = append(resolvedIncludes, &ResolvedInclude{
				ModulePath: module.Path,
				Includes:   resolved,
			})
		}
		if len(unresolved) > 0 {
			unresolvedIncludes = append(unresolvedIncludes, &UnresolvedInclude{
				ModulePath: module.Path,
				Includes:   unresolved,
			})
		}
	}
	return resolvedIncludes, unresolvedIncludes
}

func ProcessDirectory(dir string) ([]*Module, error) {
	var modules []*Module
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && filepath.Ext(path) == ".nf" {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				module, err, _ := BuildModule(path)
				if err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("%s: %v", path, err))
					mu.Unlock()
					return
				}
				mu.Lock()
				modules = append(modules, module)
				mu.Unlock()
			}(path)
		}
		return nil
	})

	wg.Wait()

	if err != nil {
		return nil, err
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("encountered %d errors: %v", len(errors), errors)
	}

	return modules, nil
}
