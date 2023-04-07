package protobuf

import (
	"encoding/base64"
	"github.com/jhump/protoreflect/desc/protoparse"
	"google.golang.org/protobuf/encoding/protojson"
	"log"

	"go.k6.io/k6/js/modules"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func init() {
	modules.Register("k6/x/protobuf", new(Protobuf))
}

type Protobuf struct {
	messageDesc protoreflect.MessageDescriptor
}

func NewProtobuf(protoFilePath string) *Protobuf {
	// Read the .proto file directly
	parser := protoparse.Parser{}
	fileDesc, err := parser.ParseFiles(protoFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the *desc.FileDescriptor to *descriptorpb.FileDescriptorProto
	schema := fileDesc[0].AsFileDescriptorProto()

	// Convert the FileDescriptorProto to a protoreflect.FileDescriptor
	fd, err := protodesc.NewFile(schema, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get the message descriptor
	messageDesc := fd.Messages().ByName("Example")

	return &Protobuf{messageDesc: messageDesc}
}

func (p *Protobuf) Encode(data string) string {
	dynamicMessage := dynamicpb.NewMessage(p.messageDesc)

	err := protojson.Unmarshal([]byte(data), dynamicMessage)

	if err != nil {
		log.Fatal(err)
	}

	encodedBytes, err := proto.Marshal(dynamicMessage)
	if err != nil {
		log.Fatal(err)
	}

	encodedString := base64.StdEncoding.EncodeToString(encodedBytes)
	return encodedString
}

func (p *Protobuf) Decode(data string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Fatal(err)
	}

	decodedMessage := dynamicpb.NewMessage(p.messageDesc)

	err = proto.Unmarshal(decodedBytes, decodedMessage)
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
