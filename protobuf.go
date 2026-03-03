package protobuf

import (
	"context"
	"log"
	"path/filepath"
	"strings"

	"github.com/bufbuild/protocompile"
	"google.golang.org/protobuf/encoding/protojson"

	"go.k6.io/k6/js/modules"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func init() {
	modules.Register("k6/x/protobuf", new(Protobuf))
}

type Protobuf struct{}

type ProtoFile struct {
	messageDesc protoreflect.MessageDescriptor
}

func (p *Protobuf) Load(protoFilePath, lookupType string, importPaths ...string) ProtoFile {
	// Default import paths if none provided.
	// Include CWD (matches the original behavior without ImportPaths) and
	// the proto file's directory (for imports relative to the same dir).
	if len(importPaths) == 0 {
		cwd, err := filepath.Abs(".")
		if err != nil {
			log.Fatalf("Failed to resolve current working directory: %v", err)
		}
		protoDir := filepath.Dir(protoFilePath)
		absProtoDir, err := filepath.Abs(protoDir)
		if err != nil {
			log.Fatalf("Failed to resolve absolute path for proto directory %s: %v", protoDir, err)
		}
		importPaths = []string{cwd}
		if absProtoDir != cwd {
			importPaths = append(importPaths, absProtoDir)
		}
	} else {
		// Convert all provided import paths to absolute paths
		absImportPaths := make([]string, len(importPaths))
		for i, path := range importPaths {
			absPath, err := filepath.Abs(path)
			if err != nil {
				log.Fatalf("Failed to resolve absolute path for import path %s: %v", path, err)
			}
			absImportPaths[i] = absPath
		}
		importPaths = absImportPaths
	}

	compiler := protocompile.Compiler{
		Resolver: &protocompile.SourceResolver{
			ImportPaths: importPaths,
		},
	}

	// Make proto file path relative to an import path for protocompile.
	// protocompile joins each import path with the file path, so the file
	// path must be relative to at least one of the import paths.
	absProtoFile, err := filepath.Abs(protoFilePath)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path for proto file %s: %v", protoFilePath, err)
	}
	for _, ip := range importPaths {
		relPath, err := filepath.Rel(ip, absProtoFile)
		if err == nil && !strings.HasPrefix(relPath, "..") {
			protoFilePath = relPath
			break
		}
	}
	if filepath.IsAbs(protoFilePath) {
		log.Fatalf("Proto file %s is not under any of the provided import paths", protoFilePath)
	}

	files, err := compiler.Compile(context.Background(), protoFilePath)
	if err != nil {
		log.Fatal(err)
	}
	if files == nil {
		log.Fatal("No files were passed as arguments")
	}
	if len(files) == 0 {
		log.Fatal("Zero files were parsed")
	}

	// Extract simple name from fully qualified name (e.g., "iot.test_messages.Ping" -> "Ping")
	simpleName := lookupType
	if idx := strings.LastIndex(lookupType, "."); idx >= 0 {
		simpleName = lookupType[idx+1:]
	}
	if simpleName == "" {
		log.Fatalf("Invalid message type name: %s", lookupType)
	}

	messageDesc := files[0].Messages().ByName(protoreflect.Name(simpleName))
	if messageDesc == nil {
		log.Fatalf("Message type %s not found", lookupType)
	}

	return ProtoFile{messageDesc}
}

func (p *ProtoFile) Encode(data string) string {
	dynamicMessage := dynamicpb.NewMessage(p.messageDesc)

	err := protojson.Unmarshal([]byte(data), dynamicMessage)

	if err != nil {
		log.Fatal(err)
	}

	encodedBytes, err := proto.Marshal(dynamicMessage)
	if err != nil {
		log.Fatal(err)
	}

	return string(encodedBytes)
}

func (p *ProtoFile) Decode(decodedBytes []byte) string {

	decodedMessage := dynamicpb.NewMessage(p.messageDesc)

	err := proto.Unmarshal(decodedBytes, decodedMessage)
	if err != nil {
		log.Fatal(err)
	}

	marshalOptions := protojson.MarshalOptions{
		UseProtoNames: true,
	}

	jsonString, err := marshalOptions.Marshal(decodedMessage)
	if err != nil {
		log.Fatal(err)
	}

	return string(jsonString)
}
