package xmlify

import "strings"

// Tag represent field tag
type Tag struct {
	Name      string
	Path      string
	Transient bool //TODO implement
	Tabular   bool
	OmitEmpty bool
}

// Parse Tag parses tag
func ParseTag(tagString string) *Tag {
	tag := &Tag{}
	if tagString == "-" {
		tag.Transient = true
		return tag
	}
	elements := strings.Split(tagString, ",")
	if len(elements) == 0 {
		return tag
	}
	for _, element := range elements {
		nv := strings.Split(element, "=")
		switch len(nv) {
		case 2:
			switch strings.ToLower(strings.TrimSpace(nv[0])) {
			case "name":
				tag.Name = strings.TrimSpace(nv[1])
			case "path":
				tag.Path = strings.TrimSpace(nv[1])
			case "transient":
				tag.Transient = strings.TrimSpace(nv[1]) == "true"
			case "tabular":
				tag.Tabular = strings.TrimSpace(nv[1]) == "true"
			case "omitempty":
				tag.OmitEmpty = strings.TrimSpace(nv[1]) == "true"
			}
			continue
		case 1:
			switch strings.ToLower(element) {
			case "tabular":
				tag.Tabular = true
			case "omitempty":
				tag.OmitEmpty = true
			case "-":
				tag.Transient = true
			}
		}

	}

	return tag
}
