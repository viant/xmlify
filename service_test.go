package xmlify

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/assertly"
	"reflect"
	"testing"
	"time"
)

func Test_RegularXML_Response_Marshal_Filter_Tabular(t *testing.T) {

	type Result struct {
		Id   int     "sqlx:\"name=ID\" velty:\"names=ID|Id\""
		Name *string "sqlx:\"name=NAME\" velty:\"names=NAME|Name\""
	}

	type IntFilter struct {
		Include []int `json:",omitempty" xmlify:"omitempty,path=@include-ids"`
		Exclude []int `json:",omitempty" xmlify:"omitempty,path=@exclude-ids"`
	}

	type StringsFilter struct {
		Include []string `json:",omitempty" xmlify:"omitempty,path=@include-ids"`
		Exclude []string `json:",omitempty" xmlify:"omitempty,path=@exclude-ids"`
	}

	type Filter struct {
		ID          *IntFilter     "json:\",omitempty\" "
		UserCreated *IntFilter     "json:\",omitempty\" "
		Name        *StringsFilter "json:\",omitempty\" "
		AccountID   *IntFilter     "json:\",omitempty\" "
	}

	type Response struct {
		Result []*Result `xmlify:"name=result,tabular"`
		Sql    string    `xmlify:"name=sql"`
		Filter *Filter   `xmlify:"name=filter"`
	}

	nameStr := "name 1"
	nameStr2 := "name 2"

	var testCases = []struct {
		description   string
		rType         reflect.Type
		input         interface{}
		expected      string
		config        *Config
		depthsConfigs []*Config
		useAssertPkg  bool
	}{
		{
			description: "01 response ver. 01",
			rType:       reflect.TypeOf(Response{}),
			input: Response{
				Result: []*Result{
					{
						Id:   1,
						Name: &nameStr,
					},
					{
						Id:   2,
						Name: &nameStr2,
					},
				},
				Sql: "abc\n",
				Filter: &Filter{
					ID: &IntFilter{
						Include: []int{1, 2},
						Exclude: []int{},
					},
					Name: &StringsFilter{
						Include: []string{"Kate", "Ann"},
						Exclude: []string{"Bob", "John"},
					},
					AccountID: &IntFilter{
						Exclude: []int{444},
					},
				},
			},

			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<result>
<columns>
<column id="Id" type="long"/>
<column id="Name" type="string"/>
</columns>
<rows>
<r>
<c lg="1"/>
<c>name 1</c>
</r>
<r>
<c lg="2"/>
<c>name 2</c>
</r>
</rows>
</result>
<sql>abc
</sql>
<filter>
<ID include-ids="1,2"/>
<Name include-ids="Kate,Ann" exclude-ids="Bob,John"/>
<AccountID exclude-ids="444"/>
</filter>
</root>`,
			config:        getMixedConfig(), //getRegularConfig(),
			depthsConfigs: []*Config{},
			useAssertPkg:  false,
		},
	}
	for _, testCase := range testCases {

		if testCase.rType == nil {
			fn, ok := (testCase.input).(func() interface{})
			assert.True(t, ok, testCase.description)

			testCase.input = fn()
			testCase.rType = reflect.TypeOf(testCase.input)
			if testCase.rType.Kind() == reflect.Slice {
				testCase.rType = testCase.rType.Elem()
			}
		}

		marshaller, err := NewMarshaller(testCase.rType, testCase.config)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		marshal, err := marshaller.Marshal(testCase.input, testCase.depthsConfigs)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		actual := string(marshal)

		if testCase.useAssertPkg {
			assert.EqualValues(t, testCase.expected, actual)
			continue
		}

		assertly.AssertValues(t, testCase.expected, actual, testCase.description)
	}
}

func Test_RegularXML_Response_Marshal_Filter(t *testing.T) {

	type Result struct {
		Id   int     "sqlx:\"name=ID\" velty:\"names=ID|Id\""
		Name *string "sqlx:\"name=NAME\" velty:\"names=NAME|Name\""
	}

	type IntFilter struct {
		Include []int `json:",omitempty" xmlify:"omitempty"`
		Exclude []int `json:",omitempty" xmlify:"omitempty"`
	}

	type StringsFilter struct {
		Include []string `json:",omitempty" xmlify:"omitempty"`
		Exclude []string `json:",omitempty" xmlify:"omitempty"`
	}

	type Filter struct {
		ID          *IntFilter     "json:\",omitempty\" "
		UserCreated *IntFilter     "json:\",omitempty\" "
		Name        *StringsFilter "json:\",omitempty\" "
		AccountID   *IntFilter     "json:\",omitempty\" "
	}

	type Response struct {
		Result []*Result `xmlify:"name=result"` //NOT tabular here
		Sql    string    `xmlify:"name=sql"`
		Filter *Filter   `xmlify:"name=filter"`
	}

	nameStr := "name 1"
	nameStr2 := "name 2"

	var testCases = []struct {
		description   string
		rType         reflect.Type
		input         interface{}
		expected      string
		config        *Config
		depthsConfigs []*Config
		useAssertPkg  bool
	}{
		{
			description: "01 response ver. 01",
			rType:       reflect.TypeOf(Response{}),
			input: Response{
				Result: []*Result{
					{
						Id:   1,
						Name: &nameStr,
					},
					{
						Id:   2,
						Name: &nameStr2,
					},
				},
				Sql: "abc\n",
				Filter: &Filter{
					ID: &IntFilter{
						Include: []int{1, 2},
						Exclude: []int{},
					},
					Name: &StringsFilter{
						Include: []string{"Kate", "Ann"},
						Exclude: []string{"Bob", "John"},
					},
					AccountID: &IntFilter{
						Exclude: []int{444},
					},
				},
			},

			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<result>
<Id>1</Id>
<Name>name 1</Name>
</result>
<result>
<Id>2</Id>
<Name>name 2</Name>
</result>
<sql>abc
</sql>
<filter>
<ID>
<Include>1</Include>
<Include>2</Include>
</ID>
<Name>
<Include>Kate</Include>
<Include>Ann</Include>
<Exclude>Bob</Exclude>
<Exclude>John</Exclude>
</Name>
<AccountID>
<Exclude>444</Exclude>
</AccountID>
</filter>
</root>`,
			config:        getMixedConfig(), //getRegularConfig(),
			depthsConfigs: []*Config{},
			useAssertPkg:  false,
		},
	}
	for _, testCase := range testCases {

		if testCase.rType == nil {
			fn, ok := (testCase.input).(func() interface{})
			assert.True(t, ok, testCase.description)

			testCase.input = fn()
			testCase.rType = reflect.TypeOf(testCase.input)
			if testCase.rType.Kind() == reflect.Slice {
				testCase.rType = testCase.rType.Elem()
			}
		}

		marshaller, err := NewMarshaller(testCase.rType, testCase.config)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		marshal, err := marshaller.Marshal(testCase.input, testCase.depthsConfigs)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		actual := string(marshal)

		if testCase.useAssertPkg {
			assert.EqualValues(t, testCase.expected, actual)
			continue
		}

		assertly.AssertValues(t, testCase.expected, actual, testCase.description)
	}
}

func Test_RegularXML_Attributes_Marshal(t *testing.T) {

	type Example01 struct {
		FooIntPtr *int `xmlify:"path=@idPtr,omitempty"`
		Id        int  `xmlify:"path=@id"`
		Name      string
		Flag01    string `xmlify:"path=@flag_01"`
		Desc      string
		Flag02    string `xmlify:"path=@flag_02"`
	}

	type Example02B struct {
		Ex Example01
	}

	type ProviderTaxonomy struct {
		IncludeIds []int `xmlify:"path=@include-ids"`
	}

	type Filter struct {
		ProviderTaxonomy *ProviderTaxonomy `xmlify:"name=provider_taxonomy"`
	}

	type Response struct {
		Filter *Filter `xmlify:"name=filter"`
	}

	var testCases = []struct {
		description   string
		rType         reflect.Type
		input         interface{}
		expected      string
		config        *Config
		depthsConfigs []*Config
		useAssertPkg  bool
	}{
		{
			description: "01 simple attribute",
			rType:       reflect.TypeOf(Example02B{}),
			input: Example02B{
				Ex: Example01{
					Id:     1,
					Name:   "name for 1",
					Flag01: "f1",
					Desc:   "description 1",
					Flag02: "f2",
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<Ex id="1" flag_01="f1" flag_02="f2">
<Name>name for 1</Name>
<Desc>description 1</Desc>
</Ex>
`,
			config: getMixedConfig(),
		},
		{
			description: "02 slice as attribute",
			rType:       reflect.TypeOf(Response{}),
			input: Response{
				Filter: &Filter{
					ProviderTaxonomy: &ProviderTaxonomy{
						IncludeIds: []int{1, 2}},
				},
			},

			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<filter>
<provider_taxonomy include-ids="1,2"/>
</filter>`,
			config:        getMixedConfig(), //getRegularConfig(),
			depthsConfigs: []*Config{getTabularConfig(), getTabularConfig(), getTabularConfig()},
			useAssertPkg:  false,
		},
	}
	for _, testCase := range testCases[1:2] {

		if testCase.rType == nil {
			fn, ok := (testCase.input).(func() interface{})
			assert.True(t, ok, testCase.description)

			testCase.input = fn()
			testCase.rType = reflect.TypeOf(testCase.input)
			if testCase.rType.Kind() == reflect.Slice {
				testCase.rType = testCase.rType.Elem()
			}
		}

		marshaller, err := NewMarshaller(testCase.rType, testCase.config)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		marshal, err := marshaller.Marshal(testCase.input, testCase.depthsConfigs)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		actual := string(marshal)

		if testCase.useAssertPkg {
			assert.EqualValues(t, testCase.expected, actual)
			continue
		}

		assertly.AssertValues(t, testCase.expected, actual, testCase.description)
	}
}

func Test_RegularXML_Response_Marshal(t *testing.T) {

	type Request struct {
		QueryString string `xmlify:"name=query_string"`
		Timestamp   string `xmlify:"name=timestamp"`
		ViewId      string `xmlify:"name=viewId"`
	}

	type Result struct {
		Avails         int     `xmlify:"name=avails"`
		ClearingPrice  float64 `xmlify:"name=clearingPrice"`
		FinalHhUniqsV1 int     `xmlify:"name=finalHhUniqsV1"`
		Uniqs          int     `xmlify:"name=uniqs"`
	}

	type ProviderTaxonomy struct {
		IncludeIds []int `xmlify:"name=include-ids"`
	}

	type Filter struct {
		ProviderTaxonomy *ProviderTaxonomy `xmlify:"name=providerTaxonomy"`
	}

	type Response struct {
		Request *Request `xmlify:"name=request"`
		Result  *Result  `xmlify:"name=result,tabular"`
		Sql     string   `xmlify:"name=sql"`
		Filter  *Filter  `xmlify:"name=filter"`
	}

	var testCases = []struct {
		description   string
		rType         reflect.Type
		input         interface{}
		expected      string
		config        *Config
		depthsConfigs []*Config
		useAssertPkg  bool
	}{
		{
			description: "01 response ver. 01",
			rType:       reflect.TypeOf(Response{}),
			input: Response{
				Request: &Request{
					QueryString: "views=TOTAL&amp;from=2023-08-06&amp;to",
					Timestamp:   "2023-08023",
					ViewId:      "total",
				},
				Result: &Result{
					Avails:         2476852,
					ClearingPrice:  0.43943723015873004,
					FinalHhUniqsV1: 37500,
					Uniqs:          520000,
				},
				Sql: "werwerewrew\n",
				Filter: &Filter{
					ProviderTaxonomy: &ProviderTaxonomy{
						IncludeIds: []int{1, 2}},
				},
			},

			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<request>
<query_string>views=TOTAL&amp;amp;from=2023-08-06&amp;amp;to</query_string>
<timestamp>2023-08023</timestamp>
<viewId>total</viewId>
</request>
<result>
<columns>
<column id="avails" type="long"/>
<column id="clearingPrice" type="double"/>
<column id="finalHhUniqsV1" type="long"/>
<column id="uniqs" type="long"/>
</columns>
<rows>
<r>
<c lg="2476852"/>
<c db="0.43943723015873004"/>
<c lg="37500"/>
<c lg="520000"/>
</r>
</rows>
</result>
<sql>werwerewrew
</sql>
<filter>
<providerTaxonomy>
<include-ids>1</include-ids>
<include-ids>2</include-ids>
</providerTaxonomy>
</filter>
</root>`,
			config:        getMixedConfig(), //getRegularConfig(),
			depthsConfigs: []*Config{getTabularConfig(), getTabularConfig(), getTabularConfig()},
			useAssertPkg:  false,
		},
	}
	for _, testCase := range testCases {

		if testCase.rType == nil {
			fn, ok := (testCase.input).(func() interface{})
			assert.True(t, ok, testCase.description)

			testCase.input = fn()
			testCase.rType = reflect.TypeOf(testCase.input)
			if testCase.rType.Kind() == reflect.Slice {
				testCase.rType = testCase.rType.Elem()
			}
		}

		marshaller, err := NewMarshaller(testCase.rType, testCase.config)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		marshal, err := marshaller.Marshal(testCase.input, testCase.depthsConfigs)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		actual := string(marshal)

		if testCase.useAssertPkg {
			assert.EqualValues(t, testCase.expected, actual)
			continue
		}

		assertly.AssertValues(t, testCase.expected, actual, testCase.description)
	}
}

func Test_RegularXML_Marshal(t *testing.T) {

	type Entry struct {
		Id   int    `xmlify:"name=id"`
		Name string `xmlify:"name=name"`
	}

	type Example01 struct {
		Id   int
		Name string
	}

	type Example02 struct {
		Id int
		Ex Example01
	}

	type Example02B struct {
		Ex Example01
	}

	type Example03 struct {
		Entries []*Entry `xmlify:"name=entries"`
	}

	type Example04 struct {
		ProviderTaxonomy []string `xmlify:"name=provider_taxonomy"`
	}

	type Foo struct {
		FooId   int
		Entries []*Entry `xmlify:"name=entries"`
	}

	type Race struct {
		Entries []*Entry `xmlify:"name=entries"`
	}

	type Example06 struct {
		Foo *Foo
	}

	type Example07 struct {
		Race *Race
	}

	type RaceWrapper struct {
		Race *Race `xmlify:"name=race"`
	}

	type Example08 struct {
		RaceWrapper *RaceWrapper `xmlify:"name=race_wrapper"`
	}

	type Escape struct {
		Ex []string `xmlify:"name=ex"`
	}

	var testCases = []struct {
		description   string
		rType         reflect.Type
		input         interface{}
		expected      string
		config        *Config
		depthsConfigs []*Config
		useAssertPkg  bool
	}{
		{
			description: "01 - simple single struct - root element",
			rType:       reflect.TypeOf(Example01{}),
			input: Example01{
				Id:   1,
				Name: "name 1",
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<Id>1</Id>
<Name>name 1</Name>
</root>`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "02 - slice with one simple struct - root and row elements ",
			rType:       reflect.TypeOf(Example01{}),
			input: []Example01{
				{
					Id:   1,
					Name: "name 1",
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<row>
<Id>1</Id>
<Name>name 1</Name>
</row>
</root>`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "03 - slice with three simple struct - root and row elements",
			rType:       reflect.TypeOf(Example01{}),
			input: []Example01{
				{
					Id:   1,
					Name: "name 1",
				},
				{
					Id:   2,
					Name: "name 2",
				},
				{
					Id:   3,
					Name: "name 3",
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<row>
<Id>1</Id>
<Name>name 1</Name>
</row>
<row>
<Id>2</Id>
<Name>name 2</Name>
</row>
<row>
<Id>3</Id>
<Name>name 3</Name>
</row>
</root>`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "04 struct with nested struct - root element",
			rType:       reflect.TypeOf(Example02{}),
			input: Example02{
				Id: 2,
				Ex: Example01{
					Id:   1,
					Name: "name for 1",
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<Id>2</Id>
<Ex>
<Id>1</Id>
<Name>name for 1</Name>
</Ex>
</root>`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "05 slice of struct with nested struct - root and row elements",
			rType:       reflect.TypeOf(Example02{}),
			input: []Example02{
				{
					Id: 2,
					Ex: Example01{
						Id:   1,
						Name: "name for 1",
					},
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<row>
<Id>2</Id>
<Ex>
<Id>1</Id>
<Name>name for 1</Name>
</Ex>
</row>
</root>`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "06 struct with nested struct - no root element",
			rType:       reflect.TypeOf(Example02B{}),
			input: Example02B{
				Ex: Example01{
					Id:   1,
					Name: "name for 1",
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<Ex>
<Id>1</Id>
<Name>name for 1</Name>
</Ex>
`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "07 slice of struct with nested struct - root and row",
			rType:       reflect.TypeOf(Example02B{}),
			input: []Example02B{
				{
					Ex: Example01{
						Id:   1,
						Name: "name for 1",
					},
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<row>
<Ex>
<Id>1</Id>
<Name>name for 1</Name>
</Ex>
</row>
</root>
`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "08 nested named slice - root element, no row element",
			rType:       reflect.TypeOf(Example03{}),
			input: Example03{
				Entries: []*Entry{
					{Id: 11, Name: "Johnson"},
					{Id: 12, Name: "Tacher"},
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<entries>
<id>11</id>
<name>Johnson</name>
</entries>
<entries>
<id>12</id>
<name>Tacher</name>
</entries>
</root>`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "09 slice of nested named slice - root element, row element",
			rType:       reflect.TypeOf(Example03{}),
			input: []Example03{
				{
					Entries: []*Entry{
						{Id: 11, Name: "Johnson"},
						{Id: 12, Name: "Tacher"},
					},
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<row>
<entries>
<id>11</id>
<name>Johnson</name>
</entries>
<entries>
<id>12</id>
<name>Tacher</name>
</entries>
</row>
</root>`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		//		{
		//			description: "10 slice of strings",
		//			rType:       reflect.TypeOf(""),
		//			input:       []string{"a", "b", ""},
		//
		//			expected: `<?xml version="1.0" encoding="UTF-8" ?>
		//<root>
		//<row>a</row>
		//<row>b</row>
		//<row></row>
		//</root>`,
		//			config:       getRegularConfig(),
		//			useAssertPkg: false,
		//		},
		{
			description: "11 field with slice type - root element",
			rType:       reflect.TypeOf(Example04{}),
			input: Example04{
				ProviderTaxonomy: []string{"a", "b", ""},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<provider_taxonomy>a</provider_taxonomy>
<provider_taxonomy>b</provider_taxonomy>
<provider_taxonomy></provider_taxonomy>
</root>`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "12 nested struct - no root element",
			rType:       reflect.TypeOf(Example06{}),
			input: Example06{
				Foo: &Foo{
					FooId: 1,
					Entries: []*Entry{
						{Id: 11, Name: "Johnson"},
						{Id: 12, Name: "Tacher"},
					},
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<Foo>
<FooId>1</FooId>
<entries>
<id>11</id>
<name>Johnson</name>
</entries>
<entries>
<id>12</id>
<name>Tacher</name>
</entries>
</Foo>
`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "13 nested struct - no root element",
			rType:       reflect.TypeOf(Example07{}),
			input: Example07{
				Race: &Race{
					Entries: []*Entry{
						{Id: 11, Name: "Johnson"},
						{Id: 12, Name: "Tacher"},
					},
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<Race>
<entries>
<id>11</id>
<name>Johnson</name>
</entries>
<entries>
<id>12</id>
<name>Tacher</name>
</entries>
</Race>
`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "14 wrapped nested struct - no root element",
			rType:       reflect.TypeOf(Example08{}),
			input: Example08{
				RaceWrapper: &RaceWrapper{
					Race: &Race{
						Entries: []*Entry{
							{Id: 11, Name: "Johnson"},
							{Id: 12, Name: "Tacher"},
						},
					},
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<race_wrapper>
<race>
<entries>
<id>11</id>
<name>Johnson</name>
</entries>
<entries>
<id>12</id>
<name>Tacher</name>
</entries>
</race>
</race_wrapper>
`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "14 escaping characters",
			rType:       reflect.TypeOf(Escape{}),
			input: Escape{
				Ex: []string{
					"ampersand (&)",
					"double quotes (\")",
					"single quotes (')",
					"less than (<)",
					"greater than (>)",
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<root>
<ex>ampersand (&amp;)</ex>
<ex>double quotes (&quot;)</ex>
<ex>single quotes (&#x27;)</ex>
<ex>less than (&lt;)</ex>
<ex>greater than (&gt;)</ex>
</root>`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
		{
			description: "15 nested struct with nil slice - no root element",
			rType:       reflect.TypeOf(Example06{}),
			input: Example06{
				Foo: &Foo{
					FooId:   1,
					Entries: nil,
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?>
<Foo>
<FooId>1</FooId>
</Foo>
`,
			config:       getRegularConfig(),
			useAssertPkg: false,
		},
	}
	for _, testCase := range testCases {
		//for _, testCase := range testCases[len(testCases)-1:] {
		//for _, testCase := range testCases[0:7] {

		//for _, testCase := range testCases[0:1] {
		//for _, testCase := range testCases[1:2] {
		//for _, testCase := range testCases[2:3] { //*
		//for _, testCase := range testCases[3:4] {
		//for _, testCase := range testCases[4:5] { // * ---> 5
		//for _, testCase := range testCases[5:6] { // *
		//for _, testCase := range testCases[6:7] { // *
		//	for i, testCase := range testCases {
		//fmt.Println("====", i, " ", testCase.description)

		if testCase.rType == nil {
			fn, ok := (testCase.input).(func() interface{})
			assert.True(t, ok, testCase.description)

			testCase.input = fn()
			testCase.rType = reflect.TypeOf(testCase.input)
			if testCase.rType.Kind() == reflect.Slice {
				testCase.rType = testCase.rType.Elem()
			}
		}

		marshaller, err := NewMarshaller(testCase.rType, testCase.config)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		marshal, err := marshaller.Marshal(testCase.input, testCase.depthsConfigs)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		actual := string(marshal)

		if testCase.useAssertPkg {
			assert.EqualValues(t, testCase.expected, actual)
			continue
		}

		assertly.AssertValues(t, testCase.expected, actual, testCase.description)
	}
}

// TODO add test for all types
// TODO passing null value as string: config.NullValue = "\u0000"
func Test_TabularXML_Marshal(t *testing.T) {
	type Example01 struct {
		Day        time.Time `xmlify:"name=day" timeLayout:"2006-01-02T15:04:05.000Z07:00"`
		Inventory  int       `xmlify:"name=inventory"`
		Uniques    int       `xmlify:"name=uniques"`
		Hh_uniques int       `xmlify:"name=hh_uniques"`
	}

	type Example02 struct {
		Metroarea string
		Inventory int  `xmlify:"name=inventory"`
		Uniques   *int `xmlify:"name=uniques"`
	}

	type Example03 struct {
		Inventory     int     `xmlify:"name=inventory"`
		Uniques       int     `xmlify:"name=uniques"`
		Hhuniques     int     `xmlify:"name=hh_uniques"`
		ClearingPrice float64 `xmlify:"name=clearing_price"`
	}

	type Example04Empties struct {
		Name   string
		Name2  *string
		Count  int
		Count2 *int
	}

	type Result struct {
		Avails         int     `xmlify:"name=avails"`
		ClearingPrice  float64 `xmlify:"name=clearingPrice"`
		FinalHhUniqsV1 int     `xmlify:"name=finalHhUniqsV1"`
		Uniqs          int     `xmlify:"name=uniqs"`
	}

	//RFC3339 = "2006-01-02T15:04:05Z07:00"
	//RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	//DateTime   = "2006-01-02 15:04:05"

	var testCases = []struct {
		description   string
		rType         reflect.Type
		input         interface{}
		expected      string
		config        *Config
		depthsConfigs []*Config
		useAssertPkg  bool
	}{
		{
			description: "example 01",
			rType:       reflect.TypeOf(Example01{}),
			input: []Example01{
				{
					//Day:        parseTime("2006-01-02T15:04:05.000Z07:00", "2023-07-01T00:00:00.123+02:00"),
					//Day:        parseTime("2006-01-02T15:04:05.000Z07:00", "2023-07-01T00:00:00.123Z"),
					Day:        parseTime("2006-01-02T15:04:05.000Z07:00", "2023-07-01T00:00:00.000Z"),
					Inventory:  1000,
					Uniques:    0,
					Hh_uniques: 0,
				},
				{
					Day:        parseTime(time.DateTime, "2023-07-02 00:00:00"),
					Inventory:  1100,
					Uniques:    0,
					Hh_uniques: 0,
				},
			},
			expected: `<result>
<columns>
<column id="day" type="date"/>
<column id="inventory" type="long"/>
<column id="uniques" type="long"/>
<column id="hh_uniques" type="long"/>
</columns>
<rows>
<r>
<c dt="2023-07-01T00:00:00.000Z"/>
<c lg="1000"/>
<c lg="0"/>
<c lg="0"/>
</r>
<r>
<c dt="2023-07-02T00:00:00.000Z"/>
<c lg="1100"/>
<c lg="0"/>
<c lg="0"/>
</r>
</rows>
</result>`,
			config:       getTabularConfig(),
			useAssertPkg: false,
		},
		////
		{
			description: "example 02",
			rType:       reflect.TypeOf(Example02{}),
			input: []Example02{
				{
					Metroarea: "Los Angeles CA",
					Inventory: 4100,
					Uniques:   nil,
				},
			},
			expected: `<result>
<columns>
<column id="Metroarea" type="string"/>
<column id="inventory" type="long"/>
<column id="uniques" type="long"/>
</columns>
<rows>
<r>
<c>Los Angeles CA</c>
<c lg="4100"/>
<c nil="true"/>
</r>
</rows>
</result>`,
			config: getTabularConfig(),
		},
		{
			description: "example 03", //TOOD
			rType:       reflect.TypeOf(Example03{}),
			input: []Example03{
				{
					Inventory:     116369911269,
					Uniques:       8876062366,
					Hhuniques:     98040000,
					ClearingPrice: 5.6903252496592245,
				},
			},
			expected: `<result>
<columns>
<column id="inventory" type="long"/>
<column id="uniques" type="long"/>
<column id="hh_uniques" type="long"/>
<column id="clearing_price" type="double"/>
</columns>
<rows>
<r>
<c lg="116369911269"/>
<c lg="8876062366"/>
<c lg="98040000"/>
<c db="5.6903252496592245"/>
</r>
</rows>
</result>`,
			config: getTabularConfig(),
		},
		////
		{
			description: "example 04",
			rType:       reflect.TypeOf(Example04Empties{}),
			input: []Example04Empties{
				{
					Name:   "",
					Name2:  nil,
					Count:  0,
					Count2: nil,
				},
			},
			expected: `<result>
<columns>
<column id="Name" type="string"/>
<column id="Name2" type="string"/>
<column id="Count" type="long"/>
<column id="Count2" type="long"/>
</columns>
<rows>
<r>
<c></c>
<c nil="true"/>
<c lg="0"/>
<c nil="true"/>
</r>
</rows>
</result>`,
			config: getTabularConfig(),
		},
		{
			description: "example 05",
			rType:       reflect.TypeOf(Result{}),
			input: &Result{
				Avails:         2476852,
				ClearingPrice:  0.43943723015873004,
				FinalHhUniqsV1: 37500,
				Uniqs:          520000,
			},
			expected: `<result>
<columns>
<column id="avails" type="long"/>
<column id="clearingPrice" type="double"/>
<column id="finalHhUniqsV1" type="long"/>
<column id="uniqs" type="long"/>
</columns>
<rows>
<r>
<c lg="2476852"/>
<c db="0.43943723015873004"/>
<c lg="37500"/>
<c lg="520000"/>
</r>
</rows>
</result>`,
			config: getTabularConfig(),
		},
	}
	for _, testCase := range testCases {
		//for _, testCase := range testCases[0:1] {
		//for _, testCase := range testCases[1:2] {
		//for _, testCase := range testCases[2:3] {
		//for _, testCase := range testCases[3:4] {
		//	for i, testCase := range testCases {
		//fmt.Println("====", i, " ", testCase.description)

		if testCase.rType == nil {
			fn, ok := (testCase.input).(func() interface{})
			assert.True(t, ok, testCase.description)

			testCase.input = fn()
			testCase.rType = reflect.TypeOf(testCase.input)
			if testCase.rType.Kind() == reflect.Slice {
				testCase.rType = testCase.rType.Elem()
			}
		}

		marshaller, err := NewMarshaller(testCase.rType, testCase.config)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		marshal, err := marshaller.Marshal(testCase.input, testCase.depthsConfigs)
		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		if !assert.Nil(t, err, testCase.description) {
			continue
		}

		actual := string(marshal)

		if testCase.useAssertPkg {
			assert.EqualValues(t, testCase.expected, actual)
			continue
		}

		assertly.AssertValues(t, testCase.expected, actual, testCase.description)
	}
}

func parseTime(format, timeStr string) time.Time {
	t, err := time.Parse(format, timeStr)
	if err != nil {
		fmt.Println(err.Error())
		return time.Time{}
	}
	return t
}

func getTabularConfig() *Config {
	return &Config{
		Style:                  "tabularStyle", // Style
		RootTag:                "result",
		HeaderTag:              "columns",
		HeaderRowTag:           "column",
		HeaderRowFieldAttr:     "id",
		HeaderRowFieldTypeAttr: "type",
		DataTag:                "rows",
		DataRowTag:             "r",
		DataRowFieldTag:        "c",
		NewLineSeparator:       "\n",
		DataRowFieldTypes: map[string]string{
			"uint":    "lg",
			"uint8":   "lg",
			"uint16":  "lg",
			"uint32":  "lg",
			"uint64":  "lg",
			"int":     "lg",
			"int8":    "lg",
			"int16":   "lg",
			"int32":   "lg",
			"int64":   "lg",
			"*uint":   "lg",
			"*uint8":  "lg",
			"*uint16": "lg",
			"*uint32": "lg",
			"*uint64": "lg",
			"*int":    "lg",
			"*int8":   "lg",
			"*int16":  "lg",
			"*int32":  "lg",
			"*int64":  "lg",
			/////
			"float32": "db",
			"float64": "db",
			/////
			"string":  "string",
			"*string": "string",
			//////
			"time.Time":  "dt",
			"*time.Time": "dt",
		},
		HeaderRowFieldType: map[string]string{
			"uint":    "long",
			"uint8":   "long",
			"uint16":  "long",
			"uint32":  "long",
			"uint64":  "long",
			"int":     "long",
			"int8":    "long",
			"int16":   "long",
			"int32":   "long",
			"int64":   "long",
			"*uint":   "long",
			"*uint8":  "long",
			"*uint16": "long",
			"*uint32": "long",
			"*uint64": "long",
			"*int":    "long",
			"*int8":   "long",
			"*int16":  "long",
			"*int32":  "long",
			"*int64":  "long",
			/////
			"float32": "double",
			"float64": "double",
			/////
			"string":  "string",
			"*string": "string",
			//////
			"time.Time":  "date",
			"*time.Time": "date",
		},
		TabularNullValue: "nil=\"true\"", //TODO MFI
	}
}

func getRegularConfig() *Config {
	return &Config{
		Style: "regularStyle", // Style
		//RootTag:                "result",
		//HeaderTag:              "columns",
		//HeaderRowTag:           "column",
		//HeaderRowFieldAttr:     "id",
		//HeaderRowFieldTypeAttr: "type",
		//DataTag:                "rows",
		//DataRowTag:             "r",
		//DataRowFieldTag:        "c",
		NewLineSeparator: "\n",
		DataRowFieldTypes: map[string]string{
			"uint":    "lg",
			"uint8":   "lg",
			"uint16":  "lg",
			"uint32":  "lg",
			"uint64":  "lg",
			"int":     "lg",
			"int8":    "lg",
			"int16":   "lg",
			"int32":   "lg",
			"int64":   "lg",
			"*uint":   "lg",
			"*uint8":  "lg",
			"*uint16": "lg",
			"*uint32": "lg",
			"*uint64": "lg",
			"*int":    "lg",
			"*int8":   "lg",
			"*int16":  "lg",
			"*int32":  "lg",
			"*int64":  "lg",
			/////
			"float32": "db",
			"float64": "db",
			/////
			"string":  "string",
			"*string": "string",
			//////
			"time.Time":  "dt",
			"*time.Time": "dt",
		},
		HeaderRowFieldType: map[string]string{
			"uint":    "long",
			"uint8":   "long",
			"uint16":  "long",
			"uint32":  "long",
			"uint64":  "long",
			"int":     "long",
			"int8":    "long",
			"int16":   "long",
			"int32":   "long",
			"int64":   "long",
			"*uint":   "long",
			"*uint8":  "long",
			"*uint16": "long",
			"*uint32": "long",
			"*uint64": "long",
			"*int":    "long",
			"*int8":   "long",
			"*int16":  "long",
			"*int32":  "long",
			"*int64":  "long",
			/////
			"float32": "double",
			"float64": "double",
			/////
			"string":  "string",
			"*string": "string",
			//////
			"time.Time":  "date",
			"*time.Time": "date",
		},
		TabularNullValue: "nil=\"true\"", //TODO MFI
		RegularRootTag:   "root",
		RegularRowTag:    "row",
	}
}

func getMixedConfig() *Config {
	return &Config{
		Style:                  "regularStyle", // Style
		RootTag:                "result",
		HeaderTag:              "columns",
		HeaderRowTag:           "column",
		HeaderRowFieldAttr:     "id",
		HeaderRowFieldTypeAttr: "type",
		DataTag:                "rows",
		DataRowTag:             "r",
		DataRowFieldTag:        "c",
		NewLineSeparator:       "\n",
		DataRowFieldTypes: map[string]string{
			"uint":    "lg",
			"uint8":   "lg",
			"uint16":  "lg",
			"uint32":  "lg",
			"uint64":  "lg",
			"int":     "lg",
			"int8":    "lg",
			"int16":   "lg",
			"int32":   "lg",
			"int64":   "lg",
			"*uint":   "lg",
			"*uint8":  "lg",
			"*uint16": "lg",
			"*uint32": "lg",
			"*uint64": "lg",
			"*int":    "lg",
			"*int8":   "lg",
			"*int16":  "lg",
			"*int32":  "lg",
			"*int64":  "lg",
			/////
			"float32": "db",
			"float64": "db",
			/////
			"string":  "string",
			"*string": "string",
			//////
			"time.Time":  "dt",
			"*time.Time": "dt",
		},
		HeaderRowFieldType: map[string]string{
			"uint":    "long",
			"uint8":   "long",
			"uint16":  "long",
			"uint32":  "long",
			"uint64":  "long",
			"int":     "long",
			"int8":    "long",
			"int16":   "long",
			"int32":   "long",
			"int64":   "long",
			"*uint":   "long",
			"*uint8":  "long",
			"*uint16": "long",
			"*uint32": "long",
			"*uint64": "long",
			"*int":    "long",
			"*int8":   "long",
			"*int16":  "long",
			"*int32":  "long",
			"*int64":  "long",
			/////
			"float32": "double",
			"float64": "double",
			/////
			"string":  "string",
			"*string": "string",
			//////
			"time.Time":  "date",
			"*time.Time": "date",
		},
		TabularNullValue: "nil=\"true\"", //TODO MFI
		RegularRootTag:   "root",
		RegularRowTag:    "row",
		RegularNullValue: "",
	}
}
