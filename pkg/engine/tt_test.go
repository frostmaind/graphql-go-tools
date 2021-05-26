package engine

import (
	"fmt"
	"testing"

	"github.com/buger/jsonparser"
)

var representationPath = []string{"body", "variables", "representations"}


func TestWorkWithJSONParser(t *testing.T) {
	rawData := []byte(`{"body":"A highly effective form of birth control.","product":{"reviews":[{"author":{"id":"1234"}},{"author":{"id":"1234"}},{"author":{"id":"7777"}}],"upc":"top-1"}}`)


	rawData2, err := jsonparser.Set(rawData, []byte("test string"), "new_field")
	if err != nil {
		fmt.Println(">>>", err)
	}

	fmt.Println("rawData", string(rawData))
	fmt.Println("rawData2", string(rawData2))
}
