package graphql

import (
	"bytes"

	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/astnormalization"
	"github.com/jensneuse/graphql-go-tools/pkg/astparser"
	"github.com/jensneuse/graphql-go-tools/pkg/astprinter"
	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

type NormalizationResult struct {
	Successful bool
	Errors     Errors
}

func (r *Request) Normalize(schema *Schema) (result NormalizationResult, err error) {
	if schema == nil {
		return NormalizationResult{Successful: false, Errors: nil}, ErrNilSchema
	}

	report := r.parseQueryOnce()
	if report.HasErrors() {
		return normalizationResultFromReport(report)
	}

	r.document.Input.Variables = r.Variables

	normalizer := astnormalization.NewWithOpts(
		astnormalization.WithExtractVariables(),
		astnormalization.WithRemoveFragmentDefinitions(),
		astnormalization.WithRemoveUnusedVariables(),
	)

	if r.OperationName != "" {
		normalizer.NormalizeNamedOperation(&r.document, &schema.document, []byte(r.OperationName), &report)
	} else {
		normalizer.NormalizeOperation(&r.document, &schema.document, &report)
	}

	if report.HasErrors() {
		return normalizationResultFromReport(report)
	}

	if r.normalizedASTCache != nil {
		r.cacheNormalizedDocument()
	}

	r.isNormalized = true
	r.Variables = r.document.Input.Variables
	if err := normalizeDuplicatedFieldRefs(&r.document); err != nil {
		return NormalizationResult{}, err
	}

	return NormalizationResult{Successful: true, Errors: nil}, nil
}

func normalizeDuplicatedFieldRefs(operation *ast.Document) error {
	buf := &bytes.Buffer{}

	if err := astprinter.Print(operation, nil, buf); err != nil {
		return err
	}

	rawQuery := buf.Bytes()
	variables := make([]byte, len(operation.Input.Variables))
	copy(variables, operation.Input.Variables)

	operation.Reset()
	operation.Input.ResetInputBytes(rawQuery)
	operation.Input.Variables = variables
	report := &operationreport.Report{}
	parser := astparser.NewParser()
	parser.Parse(operation, report)

	if report.HasErrors() {
		return report
	}

	return nil
}

func normalizationResultFromReport(report operationreport.Report) (NormalizationResult, error) {
	result := NormalizationResult{
		Successful: false,
		Errors:     nil,
	}

	if !report.HasErrors() {
		result.Successful = true
		return result, nil
	}

	result.Errors = RequestErrorsFromOperationReport(report)

	var err error
	if len(report.InternalErrors) > 0 {
		err = report.InternalErrors[0]
	}

	return result, err
}
