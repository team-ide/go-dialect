package dialect

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMapping(t *testing.T) {
	mappingSql, err := ParseMapping(mappingMySql)
	if err != nil {
		panic(err)
	}

	for key, value := range mappingSql.SqlTemplates {
		fmt.Println(*key, "-content", value.Content)
		bs, _ := json.Marshal(value.Root.Statements)
		fmt.Println(*key, "-node", string(bs))
	}
}
