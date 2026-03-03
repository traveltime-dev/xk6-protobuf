package protobuf

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

func TestLoadSimpleProto(t *testing.T) {
	p := &Protobuf{}
	pf := p.Load("example/v1/country.proto", "Country")
	if pf.messageDesc == nil {
		t.Fatal("expected messageDesc to be non-nil")
	}
}

func TestLoadProtoWithImport(t *testing.T) {
	p := &Protobuf{}
	pf := p.Load("example/v1/example.proto", "CountryList")
	if pf.messageDesc == nil {
		t.Fatal("expected messageDesc to be non-nil")
	}
}

func TestLoadWithAbsolutePath(t *testing.T) {
	p := &Protobuf{}
	absPath, err := filepath.Abs("example/v1/country.proto")
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}
	pf := p.Load(absPath, "Country")
	if pf.messageDesc == nil {
		t.Fatal("expected messageDesc to be non-nil")
	}
}

func TestLoadWithExplicitImportPaths(t *testing.T) {
	p := &Protobuf{}
	// Load using just the filename, relying on the explicit import path
	// to resolve it. Without "example/v1", protocompile can't find "country.proto".
	pf := p.Load("country.proto", "Country", "example/v1")
	if pf.messageDesc == nil {
		t.Fatal("expected messageDesc to be non-nil")
	}
}

func TestLoadWithFQN(t *testing.T) {
	p := &Protobuf{}
	pf := p.Load("example/v1/country.proto", "example.v1.Country")
	if pf.messageDesc == nil {
		t.Fatal("expected messageDesc to be non-nil")
	}
}

func TestEncodeAndDecode(t *testing.T) {
	p := &Protobuf{}
	pf := p.Load("example/v1/country.proto", "Country")

	input := `{"name":"Poland","area":312696,"population":38000000,"capital":"Warsaw"}`
	encoded := pf.Encode(input)
	if len(encoded) == 0 {
		t.Fatal("encoded output is empty")
	}

	decoded := pf.Decode([]byte(encoded))
	// Verify the decoded JSON contains the same data
	var original, result map[string]interface{}
	if err := json.Unmarshal([]byte(input), &original); err != nil {
		t.Fatalf("failed to parse input JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(decoded), &result); err != nil {
		t.Fatalf("failed to parse decoded JSON: %v", err)
	}
	if result["name"] != original["name"] {
		t.Errorf("name mismatch: got %v, want %v", result["name"], original["name"])
	}
	if result["capital"] != original["capital"] {
		t.Errorf("capital mismatch: got %v, want %v", result["capital"], original["capital"])
	}
}
