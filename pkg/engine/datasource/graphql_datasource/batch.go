package graphql_datasource

import (
	"github.com/jensneuse/graphql-go-tools/pkg/engine/resolve"
)

func ComposeFederationEntityBatch(fetches ...*resolve.SingleFetch) *resolve.SingleFetch {
	var variables []resolve.Variables
	var representations []string

	for i, fetch := range fetches {
		for _, v := range fetch.Variables {
			switch v.VariableKind() {
			case resolve.VariableKindObject:
				v
			}
		}
	}

	return nil
}
