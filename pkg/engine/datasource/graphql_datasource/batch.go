package graphql_datasource

import (
	"fmt"

	"github.com/buger/jsonparser"
)

var representationPath = []string{"body", "variables", "representations"}

func representationPathByIdx(idx int) []string {
	return []string{"body", "variables", "representations", fmt.Sprintf("[%d]", idx)}
}

func mergeFederationInputs(dst []byte, rest ...[]byte) error {
	if len(rest) == 0 {
		return nil
	}

	var nextIndex int

	_, err := jsonparser.ArrayEach(dst, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		nextIndex++
	}, representationPath...)
	if err != nil {
		return err
	}

	for i := range rest {
		_, err = jsonparser.ArrayEach(rest[i], func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			_, _ = jsonparser.Set(dst, value, representationPathByIdx(nextIndex)...)
			nextIndex++
		}, representationPath...)
		if err != nil {
			return err
		}
	}

	return nil
}
