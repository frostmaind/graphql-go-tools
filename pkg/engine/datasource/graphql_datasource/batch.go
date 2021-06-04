package graphql_datasource

import (
	"bytes"

	"github.com/buger/jsonparser"
	"github.com/cespare/xxhash"

	"github.com/jensneuse/graphql-go-tools/pkg/engine/resolve"
)

var representationPath = []string{"body", "variables", "representations"}

func mergeFederationInputs(inputs ...[]byte) (*resolve.BatchInput, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	var variables [][]byte
	var currOutPosition int
	var err error

	outToInPositions := make(map[int][]int, len(inputs))
	hashToOutPositions := make(map[uint64]int, len(inputs))

	hasher := xxhash.New()

	for i := range inputs {
		inputVariables, _, _, err := jsonparser.Get(inputs[i], representationPath...)
		if err != nil {
			return nil, err
		}

		if _, err = hasher.Write(inputVariables); err != nil {
			return nil, err
		}
		// deduplicate inputs, do not send the same representation inputVariables
		inputHash := hasher.Sum64()
		hasher.Reset()

		if outPosition, ok := hashToOutPositions[inputHash]; ok {
			outToInPositions[outPosition] = append(outToInPositions[outPosition], i)
			continue
		}

		hashToOutPositions[inputHash] = currOutPosition
		outToInPositions[currOutPosition] = []int{i}
		currOutPosition++

		_, err = jsonparser.ArrayEach(inputVariables, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			variables = append(variables, value)
		})
		if err != nil {
			return nil, err
		}
	}

	batchInput := &resolve.BatchInput{OutToInPositions: outToInPositions}

	representationJson := append([]byte("["), append(bytes.Join(variables, []byte(",")), []byte("]")...)...)
	if batchInput.Input, err = jsonparser.Set(inputs[0], representationJson, representationPath...); err != nil {
		return nil, err
	}

	return batchInput, nil
}
