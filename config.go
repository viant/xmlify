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
		Style string
		// TODO MFI move below to another config
		RootTag                string
		HeaderTag              string
		HeaderRowTag           string
		HeaderRowFieldAttr     string
		HeaderRowFieldTypeAttr string
		HeaderRowFieldType     map[string]string
		DataTag                string
		DataRowTag             string
		DataRowFieldTag        string
		DataRowFieldTypes      map[string]string
		TabularNullValue       string
		NewLine                string
		////
		RegularRootTag   string
		RegularRowTag    string
		RegularNullValue string
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

	if c.Style == "" { //TODO MFI
		c.Style = regularStyle
	}

	if c.Style == tabularStyle {
		if err := c.initTabular(); err != nil {
			return nil, err
		}
	}
	return excluded, nil
}

func (c *Config) initTabular() error {

	// TODO move config to test or crate default
	///
	//if c.DataRowFieldTypes == nil {
	//	c.DataRowFieldTypes = make(map[string]string)
	//
	//	//TODO MFI move to init
	//	c.DataRowFieldTypes["int"] = "lg"
	//	c.DataRowFieldTypes["*int"] = "lg"
	//	c.DataRowFieldTypes["time.Time"] = "dt"
	//	c.DataRowFieldTypes["string"] = "string"
	//}
	//
	//if c.HeaderRowFieldType == nil {
	//	c.HeaderRowFieldType = make(map[string]string)
	//
	//	//TODO MFI move to init
	//	c.HeaderRowFieldType["int"] = "long"
	//	c.HeaderRowFieldType["time.Time"] = "date"
	//	c.HeaderRowFieldType["string"] = "string"
	//}
	///

	if c.DataRowFieldTypes == nil {
		return fmt.Errorf("data row fields types expected for %s", tabularStyle)
	}

	if c.HeaderRowFieldType == nil {
		return fmt.Errorf("header row fields types expected for %s", tabularStyle)
	}

	return nil
}
