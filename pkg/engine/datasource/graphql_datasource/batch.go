package graphql_datasource

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"

	"github.com/jensneuse/graphql-go-tools/pkg/engine/resolve"
)

const variableWrapper = "$$"

func representationPathByIdx(idx int) []string {
	return []string{"body", "variables", "representations", fmt.Sprintf("[%d]", idx)}
}

func ConfigureBatch(fetch *resolve.SingleFetch, representationsVars ...resolve.Variables) error {
	rawInput := []byte(fetch.Input)

	initialRepresentation, err := jsonparser.GetUnsafeString(rawInput, representationPathByIdx(0)...)
	if err != nil {
		return err
	}

	variablesCount := len(fetch.Variables)
	segments := strings.Split(initialRepresentation, variableWrapper)
	segmentIdxToVarIdx := make(map[int]int, len(representationsVars))

	for segIdx, seg := range segments {
		if segIdx%2 == 0 { // variable always has an even index
			if segmentIdxToVarIdx[segIdx], err = strconv.Atoi(seg); err != nil {
				return err
			}
		}
	}

	for i := range representationsVars {
		batchRepresentationIdx := i + 1 // the first element is representationsVars from initial fetch

		for segIdx, varIdx := range segmentIdxToVarIdx {
			segments[segIdx] = strconv.Itoa(varIdx + batchRepresentationIdx*variablesCount)
		}

		representation := []byte(strings.Join(segments, variableWrapper))
		representationPath := representationPathByIdx(batchRepresentationIdx)

		if rawInput, err = jsonparser.Set(rawInput, representation, representationPath...); err != nil {
			return err
		}

		varPathPrefix := fmt.Sprintf("[%d]", batchRepresentationIdx)
		for _, variable := range representationsVars[i] {
			if v, ok := variable.(*resolve.ObjectVariable); ok {
				v.Path = append([]string{varPathPrefix}, v.Path...)
			}

			fetch.Variables = append(fetch.Variables, variable)
		}
	}

	fetch.Input = string(rawInput)

	return nil
}
