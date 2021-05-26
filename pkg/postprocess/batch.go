package postprocess

import (
	"github.com/jensneuse/graphql-go-tools/pkg/engine/plan"
	"github.com/jensneuse/graphql-go-tools/pkg/engine/resolve"
)

type ProcessBatch struct {}

func (b *ProcessBatch) Process(pre plan.Plan) plan.Plan {
	switch t := pre.(type) {
	case *plan.SynchronousResponsePlan:
		b.traverseNode(t.Response.Data)
	case *plan.StreamingResponsePlan:
		b.traverseNode(t.Response.InitialResponse.Data)
		for i := range t.Response.Patches {

		}
	}
}

func (b *ProcessBatch) traverseNode(node resolve.Node) {

}
