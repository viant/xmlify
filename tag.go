package xmlify

import "strings"

// TODO MFI delete unused
// Tag represent field tag
type Tag struct {
	Column           string
	Autoincrement    bool
	PrimaryKey       bool
	Sequence         string
	Transient        bool
	Ns               string
	Generator        string
	IsUnique         bool
	Db               string
	Table            string
	RefDb            string
	RefTable         string
	RefColumn        string
	Required         bool
	NullifyEmpty     bool
	ErrorMgs         string
	PresenceProvider bool
	Bit              bool
	Encoding         string
}

// TODO MFI delete unused
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
	for i, element := range elements {
		nv := strings.Split(element, "=")
		switch len(nv) {
		case 2:
			switch strings.ToLower(strings.TrimSpace(nv[0])) {
			case "name":
				tag.Column = strings.TrimSpace(nv[1])
			case "ns":
				tag.Ns = strings.TrimSpace(nv[1])
			case "sequence":
				tag.Sequence = strings.TrimSpace(nv[1])
			case "presence":
				tag.PresenceProvider = true
				tag.Transient = true
			case "primarykey":
				tag.PrimaryKey = strings.TrimSpace(nv[1]) == "true"
			case "autoincrement":
				tag.Autoincrement = true
			case "unique":
				tag.IsUnique = strings.TrimSpace(nv[1]) == "true"
			case "db":
				tag.Db = nv[1]
			case "table":
				tag.Table = nv[1]
			case "refdb":
				tag.RefDb = nv[1]
			case "reftable":
				tag.RefTable = nv[1]
			case "refcolumn":
				tag.RefColumn = nv[1]
			case "transient":
				tag.Transient = strings.TrimSpace(nv[1]) == "true"
			case "bit":
				tag.Bit = strings.TrimSpace(nv[1]) == "true"
			case "required":
				tag.Required = strings.TrimSpace(nv[1]) == "true"
			case "errormsg":
				tag.ErrorMgs = strings.ReplaceAll(nv[1], "$coma", ",")
			case "generator":
				generatorStrat := strings.TrimSpace(nv[1])
				tag.Generator = generatorStrat
				if generatorStrat == "autoincrement" {
					tag.Autoincrement = true
					tag.Generator = ""
				}
			case "nullifyempty":
				nullifyEmpty := strings.TrimSpace(nv[1])
				tag.NullifyEmpty = nullifyEmpty == "true" || nullifyEmpty == ""
			case "enc":
				tag.Encoding = nv[1]
			}
			continue
		case 1:
			if i == 0 {
				tag.Column = strings.TrimSpace(element)
				continue
			}
			switch strings.ToLower(element) {
			case "autoincrement":
				tag.PrimaryKey = true
			case "bit":
				tag.Bit = true
			case "primarykey":
				tag.PrimaryKey = true
			case "unique":
				tag.IsUnique = true
			case "nullifyempty":
				tag.NullifyEmpty = true
			case "required":
				tag.Required = true
			case "-":
				tag.Transient = true
			case "presence":
				tag.PresenceProvider = true
				tag.Transient = true
			}
		}

	}
	tag.PrimaryKey = tag.PrimaryKey || tag.Autoincrement
	return tag
}
