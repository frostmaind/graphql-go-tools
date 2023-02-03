package ast_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wundergraph/graphql-go-tools/internal/pkg/unsafeparser"
	"github.com/wundergraph/graphql-go-tools/pkg/astprinter"
)

func mustString(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}

func TestDocument_Clone(t *testing.T) {
	run := func(operation string) func(*testing.T) {
		return func(t *testing.T) {
			operationDocument := unsafeparser.ParseGraphqlDocumentString(operation)
			clonedOperationDocument := operationDocument.Clone()
			got := mustString(astprinter.PrintStringIndent(clonedOperationDocument, nil, " "))
			want := mustString(astprinter.PrintStringIndent(&operationDocument, nil, " "))
			assert.Equal(t, want, got)
		}
	}

	t.Run("complex query", run(`
		query NestedQuery ($firstArg: String, $secondArg: Boolean, $thirdArg: Int, $fourthArg: Float){
			serviceOne(serviceOneArg: $firstArg) {
				fieldOne
				countries {
					name
				}
			}
			serviceTwo(serviceTwoArg: $secondArg){
				fieldTwo
				serviceOneResponse {
					fieldOne
				}
			}
			anotherServiceOne(anotherServiceOneArg: $thirdArg){
				fieldOne
			}
			secondServiceTwo(secondServiceTwoArg: $fourthArg){
				fieldTwo
				serviceOneField
			}
			reusingServiceOne(reusingServiceOneArg: $firstArg){
				fieldOne
			}
		}`))

	t.Run("mutation", run(`
		mutation AddTask($title: String!, $completed: Boolean!, $name: String! @fromClaim(name: "sub")) {
		  addTask(input: [{title: $title, completed: $completed, user: {name: $name}}]){
			task {
			  id
			  title
			  completed
			}
		  }
		}
	`))
}
