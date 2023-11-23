package protobuf

import (
	"log"

	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/protobuf/encoding/protojson"

	"go.k6.io/k6/js/modules"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

func init() {
	modules.Register("k6/x/protobuf", new(Protobuf))
}

type Protobuf struct{}

type ProtoFile struct {
	messageDesc protoreflect.MessageDescriptor
}

func (p *Protobuf) Load(protoFilePath, lookupType string) ProtoFile {
	// Read the .proto file directly
	parser := protoparse.Parser{}
	fileDesc, err := parser.ParseFiles(protoFilePath)

	if err != nil {

		log.Fatal(err)
	}

	// Convert the *desc.FileDescriptor to *descriptorpb.FileDescriptorProto
	schema := fileDesc[0].AsFileDescriptorProto()
	// Convert the FileDescriptorProto to a protoreflect.FileDescriptor
	fd, err := protodesc.NewFile(schema, protoregistry.GlobalFiles)
	if err != nil {
		log.Fatal(err)
	}

	// Get the message descriptor
	return ProtoFile{fd.Messages().ByName(protoreflect.Name(lookupType))}
}

func (p *ProtoFile) Encode(data string) []byte {
	dynamicMessage := dynamicpb.NewMessage(p.messageDesc)

	err := protojson.Unmarshal([]byte(data), dynamicMessage)

	if err != nil {
		log.Fatal(err)
	}

	encodedBytes, err := proto.Marshal(dynamicMessage)
	if err != nil {
		log.Fatal(err)
	}

	return encodedBytes
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
