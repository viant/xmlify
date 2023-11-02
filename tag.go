package xmlify

import (
	"fmt"
	"github.com/viant/structology/format"
	"github.com/viant/structology/tags"
	"reflect"
	"strings"
)

// Tag represent field tag
type Tag struct {
	format.Tag
	Path         string
	Tabular      bool
	OmitTagName  bool
	Cdata        bool
	NullifyEmpty bool
}

// Parse Tag parses tag
func ParseTag(rTag reflect.StructTag) (*Tag, error) {
	fTag, err := format.Parse(rTag, TagName)
	if err != nil {
		return nil, err
	}
	ret := &Tag{Tag: *fTag}
	tagSgtring := rTag.Get(TagName)
	if tagSgtring == "" {
		return ret, nil
	}

	values := tags.Values(tagSgtring)

	err = values.MatchPairs(func(key, value string) error {
		//fmt.Printf("### KEY %s VALUE %s\n", key, value)
		switch key {

		case "path":
			ret.Path = strings.TrimSpace(value)
		case "tabular":
			ret.Tabular = strings.TrimSpace(value) == "true" || value == ""
		case "cdata":
			ret.Cdata = strings.TrimSpace(value) == "true" || value == ""
		case "omittagname":
			ret.OmitTagName = strings.TrimSpace(value) == "true" || value == ""
		case "nullifyempty":
			ret.NullifyEmpty = strings.TrimSpace(value) == "true" || value == ""
		default:
			if !format.IsValidTagKey(key) {
				return fmt.Errorf("unsupportedxmlfy option:%s", key)
			}
		}

		return nil
	})

	return ret, err
}
