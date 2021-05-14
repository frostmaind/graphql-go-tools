package graphql_datasource

import (
	"fmt"

	"github.com/jensneuse/graphql-go-tools/pkg/engine/resolve"
)

func ComposeFederationEntityBatch(fetches ...*resolve.SingleFetch) *resolve.SingleFetch {
	var variables []resolve.Variable
	var representations []string
	var lastVariableIndex

	for i, fetch := range fetches {
		idx := fmt.Sprintf("[%d]", i)

		for _, v := range fetch.Variables {
			switch v := v.(type) {
			case *resolve.ObjectVariable:
				v.Path = append([]string{idx}, v.Path...)
			}

			variables = append(variables, v)
		}
	}

	return nil
}
