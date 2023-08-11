package xmlify

import (
	"fmt"
	"github.com/viant/sqlx/io"
)

const regularStyle = "regularStyle"
const cdataMode = "cdataMode"
const tabularStyle = "tabularStyle"

type (
	Config struct {
		style string
		// TODO MFI move below to another config
		rootTag                string
		headerTag              string
		headerRowTag           string
		headerRowFieldAttr     string
		headerRowFieldTypeAttr string
		headerRowFieldType     map[string]string
		dataTag                string
		dataRowTag             string
		dataRowFieldTag        string
		dataRowFieldTypes      map[string]string
		nullValueTODO          string
		newLine                string
		////
		regularRootTag string
		regularRowTag  string
		// TODO map of types for header and data
		FieldSeparator    string
		ObjectSeparator   string
		EncloseBy         string
		EscapeBy          string
		NullValue         string
		Stringify         StringifyConfig
		UniqueFields      []string
		References        []*Reference // parent -> children. Example02.ID -> Boo.FooId
		ExcludedPaths     []string
		StringifierConfig io.StringifierConfig
	}

	// StringifyConfig "extends" Config with ignore flags
	StringifyConfig struct {
		IgnoreFieldSeparator  bool
		IgnoreObjectSeparator bool
		IgnoreEncloseBy       bool
	}
)

func (c *Config) init() (map[string]bool, error) {
	if c.EncloseBy == "" {
		c.EncloseBy = `"`
	}

	if c.EscapeBy == "" {
		c.EscapeBy = `\`
	}

	if c.FieldSeparator == "" {
		c.FieldSeparator = `,`
	}

	if c.ObjectSeparator == "" {
		c.ObjectSeparator = "\n"
	}

	if c.NullValue == "" {
		//c.NullValue = "null"
		c.NullValue = "\u0000" //"\u001f" // or other sequence
	}

	if c.StringifierConfig.StringifierFloat32Config.Precision == "" {
		c.StringifierConfig.StringifierFloat32Config.Precision = "-1"
	}

	if c.StringifierConfig.StringifierFloat64Config.Precision == "" {
		c.StringifierConfig.StringifierFloat64Config.Precision = "-1"
	}

	excluded := map[string]bool{}
	for _, path := range c.ExcludedPaths {
		excluded[path] = true
	}

	if c.style == "" { //TODO MFI
		c.style = regularStyle
	}

	if c.style == tabularStyle {
		if err := c.initTabular(); err != nil {
			return nil, err
		}
	}
	return excluded, nil
}

func (c *Config) initTabular() error {

	// TODO move config to test or crate default
	///
	//if c.dataRowFieldTypes == nil {
	//	c.dataRowFieldTypes = make(map[string]string)
	//
	//	//TODO MFI move to init
	//	c.dataRowFieldTypes["int"] = "lg"
	//	c.dataRowFieldTypes["*int"] = "lg"
	//	c.dataRowFieldTypes["time.Time"] = "dt"
	//	c.dataRowFieldTypes["string"] = "string"
	//}
	//
	//if c.headerRowFieldType == nil {
	//	c.headerRowFieldType = make(map[string]string)
	//
	//	//TODO MFI move to init
	//	c.headerRowFieldType["int"] = "long"
	//	c.headerRowFieldType["time.Time"] = "date"
	//	c.headerRowFieldType["string"] = "string"
	//}
	///

	if c.dataRowFieldTypes == nil {
		return fmt.Errorf("data row fields types expected for %s", tabularStyle)
	}

	if c.headerRowFieldType == nil {
		return fmt.Errorf("header row fields types expected for %s", tabularStyle)
	}

	return nil
}