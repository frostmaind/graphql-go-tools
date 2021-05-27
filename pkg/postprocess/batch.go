package postprocess

import (
	"fmt"

	"github.com/jensneuse/graphql-go-tools/pkg/engine/plan"
	"github.com/jensneuse/graphql-go-tools/pkg/engine/resolve"
)

type ProcessBatch struct{}

func (b *ProcessBatch) Process(pre plan.Plan) plan.Plan {
	switch t := pre.(type) {
	case *plan.SynchronousResponsePlan:
		b.traverseRootNode(t.Response.Data)
	case *plan.StreamingResponsePlan:
		b.traverseRootNode(t.Response.InitialResponse.Data)
		for i := range t.Response.Patches {
			b.traverseRootNode(t.Response.Patches[i].Value)
		}
	case *plan.SubscriptionResponsePlan:
		b.traverseRootNode(t.Response.Response.Data)
	}

	return pre
}

func (b *ProcessBatch) traverseRootNode(node resolve.Node) {
	b.traverseNode(node, false, 0)
}

func (b *ProcessBatch) traverseNode(node resolve.Node, hasBufferID bool, bufferID int) {
	var path *[]string

	switch n := node.(type) {
	case *resolve.Object:
		path = &n.Path

		for _, field := range n.Fields {
			b.traverseNode(field.Value, field.HasBuffer, field.BufferID)
			// field will be resolved from data arguments
			field.HasBuffer = false
			field.BufferID = 0
		}

	case *resolve.Array:
		path = &n.Path
	case *resolve.String:
		path = &n.Path
	case *resolve.Boolean:
		path = &n.Path
	case *resolve.Integer:
		path = &n.Path
	case *resolve.Float:
		path = &n.Path
	default:
		return
	}

	if hasBufferID {
		rootKey := fmt.Sprintf("fetch_%d", bufferID)
		*path = append([]string{rootKey}, *path...)
	}
}
