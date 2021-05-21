package graphql_datasource

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"

	"github.com/jensneuse/graphql-go-tools/pkg/engine/resolve"
)

const variableWrapper = "$$"

var representationPath = []string{"body", "variables", "representations", "[0]"}

func ComposeFederationEntityBatch(fetch *resolve.SingleFetch, variables ...resolve.Variables) (*resolve.SingleFetch, error) {
	initialRepresentation, err := jsonparser.GetUnsafeString([]byte(fetch.Input), representationPath...)
	if err != nil {
		return nil, err
	}

	batchRepresentations := make([]string, 0, len(variables)+1)
	batchVariables := make(resolve.Variables, 0, len(variables)+1)
	// add the representation and variables from SingleFetch
	batchRepresentations = append(batchRepresentations, initialRepresentation)
	batchVariables = append(batchVariables, fetch.Variables...)

	variablesCount := len(fetch.Variables)
	segments := strings.Split(initialRepresentation, variableWrapper)
	segmentIdxToVarIdx := make(map[int]int, len(variables))

	for segIdx, seg := range segments {
		if segIdx%2 == 0 { // variable always has an even index
			if segmentIdxToVarIdx[segIdx], err = strconv.Atoi(seg); err != nil {
				return nil, err
			}
		}
	}

	for i := range variables {
		totalVarOrderNum := i + 1 // the first element is variables from initial fetch

		for segIdx, varIdx := range segmentIdxToVarIdx {
			segments[segIdx] = strconv.Itoa(varIdx + totalVarOrderNum*variablesCount)
		}

		batchRepresentations = append(batchRepresentations, strings.Join(segments, variableWrapper))

		varPathPrefix := fmt.Sprintf("[%d]", totalVarOrderNum)
		for _, variable := range variables[i] {
			if v, ok := variable.(*resolve.ObjectVariable); ok {
				v.Path = append([]string{varPathPrefix}, v.Path...)
			}

			batchVariables = append(batchVariables, variable)
		}
	}

	return nil, nil
}
