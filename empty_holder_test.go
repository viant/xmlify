package xmlify

import (
	"reflect"
	"strings"
	"testing"
)

func TestPreserveEmptyTypedHolders(t *testing.T) {
	type child struct{ Name string }
	type row struct {
		ID       int
		Children []child
		Single   *child
	}
	for _, empty := range []bool{false, true} {
		value := row{}
		if empty {
			value.Children = []child{}
		}
		m, err := NewMarshaller(reflect.TypeOf(value), &Config{RegularRootTag: "result", RegularRowTag: "row", RegularNullValue: `nil="true"`, PreserveEmptyHolders: true})
		if err != nil {
			t.Fatal(err)
		}
		got, err := m.Marshal([]row{value})
		if err != nil {
			t.Fatal(err)
		}
		wanted := `<Children nil="true"/>`
		if empty {
			wanted = `<Children/>`
		}
		if !strings.Contains(string(got), wanted) || !strings.Contains(string(got), `<Single nil="true"/>`) {
			t.Fatalf("%s", got)
		}
	}
}
