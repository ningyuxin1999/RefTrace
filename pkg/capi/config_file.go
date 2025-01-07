package main

// #include <stdlib.h>
import "C"
import (
	"encoding/base64"
	"reft-go/nf/configlint"
	pb "reft-go/nf/proto"
	"reft-go/parser"
	"unsafe"

	"google.golang.org/protobuf/proto"
)

type ConfigFile struct {
	Path          string
	ProcessScopes []configlint.ProcessScope
}

func (c *ConfigFile) ToProto() *pb.ConfigFile {
	protoConfig := &pb.ConfigFile{
		Path: c.Path,
	}

	for _, scope := range c.ProcessScopes {
		protoConfig.ProcessScopes = append(protoConfig.ProcessScopes, scope.ToProto())
	}

	return protoConfig
}

//export ConfigFile_New
func ConfigFile_New(filePath *C.char) *C.char {
	goPath := C.GoString(filePath)

	// Parse the AST
	ast, err := parser.BuildAST(goPath)
	if err != nil {
		parseError := &pb.ParseError{}
		likelyRtBug := false
		if _, ok := err.(*parser.SyntaxException); ok {
			likelyRtBug = true
		}
		parseError.LikelyRtBug = likelyRtBug
		parseError.Error = err.Error()
		return serializeResult(&pb.ConfigFileResult{
			Result: &pb.ConfigFileResult_Error{
				Error: parseError,
			},
		})
	}

	// Parse the config
	processScopes := configlint.ParseConfig(ast.StatementBlock)

	config := &ConfigFile{
		Path:          goPath,
		ProcessScopes: processScopes,
	}

	return serializeResult(&pb.ConfigFileResult{
		Result: &pb.ConfigFileResult_ConfigFile{
			ConfigFile: config.ToProto(),
		},
	})
}

//export ConfigFile_Free
func ConfigFile_Free(ptr *C.char) {
	C.free(unsafe.Pointer(ptr))
}

func serializeResult(result *pb.ConfigFileResult) *C.char {
	bytes, err := proto.Marshal(result)
	if err != nil {
		panic("serialization error: " + err.Error())
	}
	return C.CString(base64.StdEncoding.EncodeToString(bytes))
}
