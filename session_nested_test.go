package xmlify

import (
	"encoding/xml"
	"io"
	"reflect"
	"strings"
	"testing"
)

type nestedCustomXML struct{ Value string }

func (v *nestedCustomXML) MarshalXML() ([]byte, error) {
	return []byte("<custom>" + v.Value + "</custom>"), nil
}

func TestNestedCustomMarshallerUsesParentFieldPointer(t *testing.T) {
	type row struct {
		Prefix, Padding int64
		Payload         *nestedCustomXML
	}
	marshaller, err := NewMarshaller(reflect.TypeOf(row{}), &Config{RegularRootTag: "Records", RegularRowTag: "Record"})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := marshaller.Marshal([]row{{Prefix: 1, Padding: 2, Payload: &nestedCustomXML{Value: "Ada"}}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(encoded), "<custom>Ada</custom>") {
		t.Fatalf("custom nested marshaller not called: %s", encoded)
	}
}

func TestMarshallerNestedSession(t *testing.T) {
	type profile struct{ Name string }
	type valueRow struct {
		ID      int
		Profile profile
	}
	type pointerRow struct {
		ID      int
		Profile *profile
	}
	for _, test := range []struct {
		name      string
		row, rows any
	}{
		{"value", valueRow{}, []valueRow{{ID: 1, Profile: profile{Name: "Ada"}}}},
		{"pointer", pointerRow{}, []pointerRow{{ID: 1, Profile: &profile{Name: "Ada"}}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			marshaller, err := NewMarshaller(reflect.TypeOf(test.row), &Config{RegularRootTag: "Records", RegularRowTag: "Record"})
			if err != nil {
				t.Fatal(err)
			}
			data, err := marshaller.Marshal(test.rows)
			if err != nil {
				t.Fatal(err)
			}
			decoder := xml.NewDecoder(strings.NewReader(string(data)))
			found := false
			for {
				token, err := decoder.Token()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatalf("invalid XML %s: %v", data, err)
				}
				if value, ok := token.(xml.CharData); ok && string(value) == "Ada" {
					found = true
				}
			}
			if !found {
				t.Fatalf("nested value absent: %s", data)
			}
		})
	}
}

func TestNestedSessionAppenderOwnsSliceStorage(t *testing.T) {
	type profile struct{ Name string }
	for _, parent := range []reflect.Type{reflect.TypeOf(profile{}), reflect.TypeOf(&profile{})} {
		field := &Field{path: "Profile", parentType: parent}
		marshal := &UnmarshalSession{}
		if dest, appender := marshal.destWithAppender(field); dest != nil || appender != nil {
			t.Fatal("marshal allocated unmarshal state")
		}
		var roots []profile
		session := &UnmarshalSession{dest: &roots}
		dest, appender := session.destWithAppender(field)
		if reflect.TypeOf(dest) != reflect.PointerTo(reflect.SliceOf(parent)) {
			t.Fatalf("nested destination type %T", dest)
		}
		appender.Append(&profile{Name: "Ada"})
		if reflect.ValueOf(dest).Elem().Len() != 1 {
			t.Fatal("nested append failed")
		}
	}
}

func TestNestedSlicesAndRepeatedSiblingTypes(t *testing.T) {
	type filter struct{ Values []int }
	type siblings struct{ First, Second *filter }
	type response struct{ Filters *siblings }
	for _, input := range []any{
		response{Filters: &siblings{First: &filter{Values: []int{11, 12}}, Second: &filter{Values: []int{21, 22}}}},
		&response{Filters: &siblings{First: &filter{Values: []int{11, 12}}, Second: &filter{Values: []int{21, 22}}}},
	} {
		marshaller, err := NewMarshaller(reflect.TypeOf(response{}), &Config{RegularRootTag: "Response", RegularRowTag: "Row"})
		if err != nil {
			t.Fatal(err)
		}
		data, err := marshaller.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		for _, value := range []string{"<First>", "<Second>", ">11<", ">12<", ">21<", ">22<"} {
			if !strings.Contains(string(data), value) {
				t.Fatalf("missing %s in %s", value, data)
			}
		}
		decoder := xml.NewDecoder(strings.NewReader(string(data)))
		for {
			_, err := decoder.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("invalid XML %s: %v", data, err)
			}
		}
	}
}
