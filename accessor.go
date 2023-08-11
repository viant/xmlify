package xmlify

import (
	io2 "github.com/viant/sqlx/io"
	//io2 "github.com/viant/sqlx/io"
	"github.com/viant/xunsafe"
	"unsafe"
)

type (
	Accessor struct {
		_parent             *Accessor
		cache               map[unsafe.Pointer]*stringified
		emitedFirst         bool
		parentAccessorIndex int
		fields              []*Field // requested fields
		children            []*Accessor
		config              *Config
		path                string

		currSliceIndex int
		slicePtr       unsafe.Pointer
		slice          *xunsafe.Slice

		ptr      unsafe.Pointer // refering to a single object
		field    *xunsafe.Field // used to get value from parent pointer
		fieldTag *Tag
		holder   string
		xType    *xunsafe.Type
	}

	stringified struct {
		values            []string
		wasStrings        []bool
		types             []string
		dataRowFieldTypes []string
		shouldWrite       []bool
		fieldNames        []string
	}
)

func (a *Accessor) Reset() {
	a.emitedFirst = false
	a.cache = map[unsafe.Pointer]*stringified{}
	a.currSliceIndex = 0
}

func (a *Accessor) Has() bool {
	if a.emitedFirst {
		return a.prepare()
	}

	a.emitedFirst = true
	return a.ptr != nil
}

func (a *Accessor) prepare() bool {
	accessor, ok := a.next()
	if !ok {
		return false
	}

	parent, childIndex := sliceParentOf(accessor)
	for i := 0; i < childIndex; i++ {
		parent.children[i].Set(parent.ptr)
	}

	return true
}

func sliceParentOf(accessor *Accessor) (*Accessor, int) {
	for accessor != nil {
		if accessor._parent != nil && accessor._parent.slice != nil {
			return accessor._parent, accessor.parentAccessorIndex
		}

		accessor = accessor._parent
	}

	return nil, -1
}

func (a *Accessor) Set(pointer unsafe.Pointer) {
	a.ptr = pointer

	for _, child := range a.children {
		if pointer == nil {
			child.Set(nil)
		} else {
			if child.field == nil {
				for _, child2 := range child.children {
					child2.Set(pointer)
				}
			} else {
				aPtr, slicePtr := a.getChildValue(pointer, child)
				child.slicePtr = slicePtr
				child.Set(aPtr)
			}
		}

		child.currSliceIndex = 0
	}
}

func (a *Accessor) getChildValue(pointer unsafe.Pointer, child *Accessor) (valuePtr unsafe.Pointer, slicePtr unsafe.Pointer) {
	valuePointer := child.field.ValuePointer(pointer)
	if child.slice == nil {
		return valuePointer, nil
	}

	lenSlice := child.slice.Len(valuePointer)
	if lenSlice == 0 {
		return nil, valuePointer
	}

	at := child.slice.ValuePointerAt(valuePointer, 0)
	return xunsafe.AsPointer(at), valuePointer
}

func (a *Accessor) Headers() ([]string, []string) {

	headers := make([]string, len(a.fields))
	headerRowFieldTypes := make([]string, len(a.fields))
	var rowFieldType string
	var ok bool

	for i, field := range a.fields {
		//tag := field.xField.Tag // TODO MFI TAG
		//if tag != "" {
		//	fmt.Printf(" %s tag = %s\n", field.name, tag)
		//}
		//	tag := field.xField.Type. //ParseTag(field.Tag.Get(option.TagSqlx))
		if headers[i] = field.tag.Name; headers[i] == "" {
			headers[i] = field.name
		}

		rowFieldType, ok = a.config.headerRowFieldType[field.xField.Type.String()]
		if ok {
			headerRowFieldTypes[i] = rowFieldType
		} else {
			headerRowFieldTypes[i] = "TYPE_ERR"
		}
	}

	/***
		for _, child := range a.children {
			//if child.config == nil {
			//	childHeaders := child.Headers()
			//	headers = append(headers, childHeaders...)
			//} else {
			//	headers = append(headers, child.holder)
			//}
			//headers = append(headers, child.holder)

			headers = append(headers, child.field.Name)
		}
	***/
	return headers, headerRowFieldTypes
}

func (a *Accessor) stringifyFields(writer *writer) ([]string, []bool, []string, []string) {
	if value, ok := a.cache[a.ptr]; ok {
		return value.values, value.wasStrings, value.types, value.dataRowFieldTypes
	}

	result := make([]string, len(a.fields))
	wasStrings := make([]bool, len(a.fields))
	types := make([]string, len(a.fields))
	dataRowFieldTypes := make([]string, len(a.fields))

	if a.ptr == nil {
		strings := make([]string, len(a.fields))
		for i := range strings {
			strings[i] = writer.config.NullValue
		}

		return strings, make([]bool, len(a.fields)), make([]string, len(a.fields)), make([]string, len(a.fields))
	}

	var rowFieldType string
	var ok bool
	for i, field := range a.fields {
		result[i], wasStrings[i] = field.stringifier(a.ptr) //TODO how to get real type here?
		types[i] = field.xField.Type.String()

		rowFieldType, ok = a.config.dataRowFieldTypes[types[i]]
		if ok {
			dataRowFieldTypes[i] = rowFieldType
		} else {
			dataRowFieldTypes[i] = "TYPE_ERR"
		}

	}

	a.cache[a.ptr] = &stringified{
		values:            result,
		wasStrings:        wasStrings,
		types:             types,
		dataRowFieldTypes: dataRowFieldTypes,
	}

	return result, wasStrings, types, dataRowFieldTypes
}

// next returns true if record was not exhausted and first not ehausted Accessor
func (a *Accessor) next() (*Accessor, bool) {
	for _, child := range a.children {
		if child.config != nil {
			continue
		}

		if accessor, ok := child.next(); ok {
			return accessor, true
		}
	}

	if a.slice != nil && a.slicePtr != nil {
		sliceLen := a.slice.Len(a.slicePtr)
		if a.currSliceIndex < sliceLen-1 {
			a.currSliceIndex++
			value := a.slice.ValuePointerAt(a.slicePtr, a.currSliceIndex)
			a.Set(xunsafe.AsPointer(value))
			return a, true
		}
	}

	return nil, false
}

func (a *Accessor) ResetAllChildren() {
	for _, child := range a.children {
		child.ResetAllChildren()
	}

	a.Reset()
}

func (a *Accessor) values() (io2.ValueAccessor, int) {
	if a.ptr == nil {
		return func(index int) interface{} {
			return nil
		}, 0
	}

	if a.slice != nil {
		if a.slicePtr == nil {
			return func(index int) interface{} {
				return nil
			}, 0
		}

		return func(index int) interface{} {
			if index >= a.slice.Len(a.slicePtr) {
				return nil
			}

			return a.slice.ValuePointerAt(a.slicePtr, index)
		}, 1
	}

	interfacer := a.Interfacer()
	return func(index int) interface{} {
		return interfacer.Interface(a.ptr)
	}, 1
}

func (a *Accessor) Interfacer() *xunsafe.Type {
	if a.xType == nil {
		a.xType = xunsafe.NewType(a.field.Type)
	}

	return a.xType
}
