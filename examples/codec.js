import protobuf from 'k6/x/protobuf';

const data = open("example.json");
const protoFile = protobuf.load("examples/example.proto", "CountryList") // todo: Fix paths

export default function () {
    console.log(protoFile.decode(protoFile.encode(data)))
}

