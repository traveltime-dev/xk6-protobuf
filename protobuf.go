package protobuf

import (
	"context"
	"errors"

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

func (p *Protobuf) Load(protoFilePath string, lookupType string) (*ProtoFile, error) {
	return p.LoadWithPath(protoFilePath, []string{}, lookupType)
}

func (p *Protobuf) LoadWithPath(protoFilePath string, importPaths []string, lookupType string) (*ProtoFile, error) {
	compiler := protocompile.Compiler{
		Resolver: &protocompile.SourceResolver{
			ImportPaths: importPaths,
		},
	}

	files, err := compiler.Compile(context.Background(), protoFilePath)
	if err != nil {
		return nil, err
	}
	if files == nil {
		return nil, errors.New("No files were passed as arguments")
	}
	if len(files) == 0 {
		return nil, errors.New("Zero files were parsed")
	}

	return &ProtoFile{files[0].Messages().ByName(protoreflect.Name(lookupType))}, nil
}

func (p *ProtoFile) Encode(data string) ([]byte, error) {
	dynamicMessage := dynamicpb.NewMessage(p.messageDesc)

	err := protojson.Unmarshal([]byte(data), dynamicMessage)

	if err != nil {
		return nil, err
	}

	encodedBytes, err := proto.Marshal(dynamicMessage)
	if err != nil {
		return nil, err
	}

	return encodedBytes, nil
}

func (p *ProtoFile) Decode(decodedBytes []byte) (string, error) {

	decodedMessage := dynamicpb.NewMessage(p.messageDesc)

	err := proto.Unmarshal(decodedBytes, decodedMessage)
	if err != nil {
		return "", err
	}

	marshalOptions := protojson.MarshalOptions{
		UseProtoNames: true,
	}

	jsonString, err := marshalOptions.Marshal(decodedMessage)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
}
