package astnormalization

import "testing"

func TestDirectiveIncludeVisitor(t *testing.T) {
	t.Run("remove static include true on inline fragment", func(t *testing.T) {
		run(directiveIncludeSkip, testDefinition, `
				{
					dog {
						name: nickname
						... @include(if: true) {
							includeName: name @include(if: true)
							notIncludeName: name @include(if: false)
							notSkipName: name @skip(if: false)
							skipName: name @skip(if: true)
						}
					}
					notInclude: dog @include(if: false) {
						name
					}
					skip: dog @skip(if: true) {
						name
					}
				}`, `
				{
					dog {
						name: nickname
						... {
							includeName: name
							notSkipName: name
						}
					}
				}`)
	})
	t.Run("apply `include` directive that contains variable", func(t *testing.T) {
		runVariables(directiveIncludeSkip, testDefinition, `
				query ApplyInclude($includeVar: Boolean!, $notInclude: Boolean!){
					dog {
						name: nickname
						includeName: name @include(if: $includeVar)
						notIncludeName: name @include(if: $notInclude)
					}
				}`,
			`{"includeVar": true, "notInclude": false}`,`
				query ApplyInclude($includeVar: Boolean!, $notInclude: Boolean!) {
					dog {
						name: nickname
						includeName: name
					}
				}`)
	})
	t.Run("apply `skip` directive that contains variable", func(t *testing.T) {
		runVariables(directiveIncludeSkip, testDefinition, `
				query ApplyInclude($skipVar: Boolean!, $notSkipVar: Boolean!){
					dog {
						name: nickname
						skipName: name @skip(if: $skipVar)
						notSkipName: name @skip(if: $notSkipVar)
					}
				}`,
			`{"skipVar": true, "notSkipVar": false}`,`
				query ApplyInclude($skipVar: Boolean!, $notSkipVar: Boolean!) {
					dog {
						name: nickname
						notSkipName: name
					}
				}`)
	})
}
