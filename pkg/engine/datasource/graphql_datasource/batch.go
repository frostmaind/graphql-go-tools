package graphql_datasource

import (
	"bytes"

	"github.com/buger/jsonparser"
)

var representationPath = []string{"body", "variables", "representations"}

func mergeFederationInputs(inputs ...[]byte) ([]byte, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	var variables [][]byte

	for i := range inputs {
		_, err := jsonparser.ArrayEach(inputs[i], func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			variables = append(variables, value)
		}, representationPath...)
		if err != nil {
			return nil, err
		}
	}

	representationJson := append([]byte("["), append(bytes.Join(variables, []byte(",")), []byte("]")...)...)

	result := make([]byte, len(inputs[0]))
	copy(result, inputs[0])

	result, err := jsonparser.Set(result, representationJson, representationPath...)
	if err != nil {
		return nil, err
	}

	return result, nil
}