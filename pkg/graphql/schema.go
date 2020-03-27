package graphql

import (
	"io"
	"io/ioutil"

	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/astparser"
	"github.com/jensneuse/graphql-go-tools/pkg/execution/boilerplate"
)

type Schema struct {
	rawInput []byte
	document ast.Document
}

func NewSchemaFromReader(reader io.Reader) (*Schema, error) {
	schemaContent, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return createSchema(schemaContent)
}

func NewSchemaFromString(schema string) (*Schema, error) {
	schemaContent := []byte(schema)

	return createSchema(schemaContent)
}

func (s *Schema) Document() []byte {
	return s.rawInput
}

func (s *Schema) Validate() (valid bool, errors SchemaValidationErrors) {
	// TODO: Needs to be implemented in the core of the library
	return true, nil
}

func createSchema(schemaContent []byte) (*Schema, error) {
	schemaContentWithBoilerplate := boilerplate.NewSchemaBytesWithBoilerplate(schemaContent)
	document, report := astparser.ParseGraphqlDocumentBytes(schemaContentWithBoilerplate)
	if report.HasErrors() {
		return nil, report
	}

	return &Schema{
		rawInput: schemaContent,
		document: document,
	}, nil
}
