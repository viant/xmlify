package xmlify

import (
	"fmt"
	"github.com/viant/xunsafe"
	"reflect"
	"unsafe"
)

func (w *writer) writeRegularAllObjects(acc *Accessor, parentLevel bool) {

	if parentLevel {
		w.buffer.writeString(`<?xml version="1.0" encoding="UTF-8" ?>`) // TODO add to config?
	}

	if acc.fieldTag != nil && acc.fieldTag.Tabular {
		w.buffer.writeString(w.config.NewLineSeparator)
		w.writeTabularAllObjects(acc, parentLevel)
		return
	}

	var fieldKind reflect.Kind
	var fieldName string
	var rowFieldName string

	if acc.field != nil {
		fieldName = acc.field.Name
		fieldKind = acc.field.Kind()
	}

	// TODO add tag as value not pointer?
	if acc.fieldTag != nil && acc.fieldTag.Name != "" {
		fieldName = acc.fieldTag.Name
	}

	// TODO move
	rowFieldName = w.config.RegularRowTag
	if fieldKind == reflect.Slice && fieldName != "" {
		rowFieldName = fieldName
	}

	// TODO user RegularRootTag only if needed
	if parentLevel && fieldName == "" {
		fieldName = w.config.RegularRootTag
	}

	omitRootElement := false

	if parentLevel && acc.slice == nil && len(acc.fields) == 1 { //TODO rebuild and add test with struct = {nestedstruct01, nestedstruct02}
		f0Type := acc.fields[0].xField.Type
		f0Kind := f0Type.Kind()
		if f0Kind == reflect.Ptr {
			f0Type = f0Type.Elem()
			f0Kind = f0Type.Kind()
		}

		if f0Kind == reflect.Struct {
			omitRootElement = true
		}
	}

	areAttributes, attrNames, containAttr := acc.attributes()

	if !(parentLevel && omitRootElement) && fieldKind != reflect.Slice && fieldName != "" {
		if containAttr {
			w.buffer.writeString(w.config.NewLineSeparator + "<" + fieldName)
		} else {
			w.buffer.writeString(w.config.NewLineSeparator + "<" + fieldName + ">")
		}

		defer w.buffer.writeString(w.config.NewLineSeparator + "</" + fieldName + ">")
	}

	headers, _ /*hTypes*/ := acc.RegularHeaders() //TODO rename
	var xType *xunsafe.Type

	for i := 0; i < w.size; i++ {
		acc.ResetAllChildren()

		if parentLevel {
			if i != 0 {
				acc.Reset()
			}
			at := w.valueAt(i)
			if i == 0 {
				if reflect.TypeOf(at).Kind() == reflect.Ptr {
					xType = w.dereferencer
				}
			}
			if xType != nil {
				at = xType.Deref(at)
			}
			acc.Set(xunsafe.AsPointer(at))
		}

		for acc.Has() {

			result, wasStrings, types, _, shouldWrite, fieldNames := acc.stringifyRegularFields(w, &headers, &attrNames)

			//len(field)< len(result) // when one of field is a slice !!!

			if fieldKind == reflect.Slice || acc.slice != nil {
				w.buffer.writeString(w.config.NewLineSeparator + "<" + rowFieldName + ">")
			}

			// writing fields one by one in correct order (field can be a slice or attribute)
			start := -1
			end := -1

			upperBound := make([]int, len(acc.fields))
			lowerBound := make([]int, len(acc.fields))
			for k, field := range acc.fields {

				start = -1
				end = -1
				for j, _ := range fieldNames {
					if field.name == fieldNames[j] {
						if start == -1 {
							start = j
							end = j
						} else {
							end = j
						}
					}
				}

				end = end + 1
				lowerBound[k] = start
				upperBound[k] = end
			}

			start = 0
			end = 0

			// handle attributes first
			for j, _ := range acc.fields {
				start = lowerBound[j]
				end = upperBound[j]

				if areAttributes[start] && shouldWrite[start] { // TODO handle slices
					w.writeRegularObjectAttr(result[start:end], wasStrings[start:end], types[start:end], []string{}, attrNames[start:end], shouldWrite[start:end])

				}
			}
			if containAttr {
				w.buffer.writeString(">")
			}

			for j, field := range acc.fields {
				start = lowerBound[j]
				end = upperBound[j]

				if !areAttributes[start] && shouldWrite[start] { // whole result[start:end] represents one field (can be a slice)
					w.writeRegularElement(result[start:end], headers[start:end], shouldWrite[start:end])
				}

				for _, child := range acc.children {
					if child.field.Name != field.name {
						continue
					}
					_, childSize := child.values()

					if childSize > 0 {
						tmpSize := w.size
						w.size = childSize
						w.writeRegularAllObjects(child, false)
						w.size = tmpSize
					}
				}
			}

			if fieldKind == reflect.Slice || acc.slice != nil {
				w.buffer.writeString(w.config.NewLineSeparator + "</" + rowFieldName + ">")
			}

		} // ~ has loop
	} // ~ main for loop
}

func WriteRegularObjectAttr(writer *Buffer, config *Config, values []string, wasString []bool, types []string, dataRowFieldTypes []string, headers []string, shouldWrite []bool) {
	if len(values) == 0 {
		return
	}

	escapedNullValue := EscapeSpecialChars(config.NullValue, config)

	currentAttr := ""
	lastAttr := ""

	for j := 0; j < len(values); j++ {
		if !shouldWrite[j] { // TODO CHECK IF NEEDED
			continue
		}

		currentAttr = headers[j]

		if currentAttr != lastAttr && lastAttr != "" { // closing last attr
			writer.writeString("\"")
		}

		if currentAttr != lastAttr { // opening new attr
			writer.writeString(" " + currentAttr + "=" + "\"")
		}

		if currentAttr == lastAttr && lastAttr != "" { //separate another value for current attr
			writer.writeString(",") // TODO ADD attr separator into config?
		}

		asString := EscapeSpecialChars(values[j], config) //TODO MFI escaping

		lastAttr = currentAttr

		if asString == escapedNullValue {
			asString = config.RegularNullValue
			writer.writeString(asString)
			continue
		}

		writer.writeString(asString)
	}

	//closing last attr
	if lastAttr != "" {
		writer.writeString("\"")
	}

	//writer.writeString("]")
}

func (w *writer) writeRegularObjectAttr(data []string, wasStrings []bool, types, dataRowFieldTypes []string, headers []string, shouldWrite []bool) {
	if w.writtenObject {
		//w.writeObjectSeparator()
	} else {
		w.buffer.writeString(w.beforeFirst) // TODO w.beforeFirst [[??
	}

	//WriteObject(w.buffer, w.config, data, wasStrings)
	WriteRegularObjectAttr(w.buffer, w.config, data, wasStrings, types, dataRowFieldTypes, headers, shouldWrite)
	w.writtenObject = true
}

func (w *writer) writeRegularElement(values []string, headers []string, shouldWrite []bool) {
	if len(values) == 0 {
		return
	}

	for j := 0; j < len(values); j++ {
		if !shouldWrite[j] {
			continue
		}

		asString := EscapeSpecialChars(values[j], w.config)

		if asString == w.config.escapedNullValue {
			if w.config.RegularNullValue == "" {
				asString = w.config.NewLineSeparator + "<" + headers[j] + "/>"
			} else {
				asString = w.config.NewLineSeparator + "<" + headers[j] +
					" " + w.config.RegularNullValue +
					"/>"
			}
			w.buffer.writeString(asString)
			continue
		}

		asString = w.config.NewLineSeparator + "<" + headers[j] + ">" +
			asString +
			"</" + headers[j] + ">"

		w.buffer.writeString(asString)
	}
	w.writtenObject = true
}

func (a *Accessor) RegularHeaders() ([]string, []string) {

	headers := make([]string, len(a.fields))
	headerRowFieldTypes := make([]string, len(a.fields))

	for i, field := range a.fields {
		if headers[i] = field.tag.Name; headers[i] == "" {
			headers[i] = field.name
		}
	}

	return headers, headerRowFieldTypes
}

func (a *Accessor) stringifyRegularFields(writer *writer, headers *[]string, attrNames *[]string) ([]string, []bool, []string, []string, []bool, []string) {
	if value, ok := a.cache[a.ptr]; ok {
		return value.values, value.wasStrings, value.types, value.dataRowFieldTypes, value.shouldWrite, value.fieldNames
	}

	shouldWrite := make([]bool, len(a.fields))
	result := make([]string, len(a.fields))
	wasStrings := make([]bool, len(a.fields))
	types := make([]string, len(a.fields))
	dataRowFieldTypes := make([]string, len(a.fields))
	fieldNames := make([]string, len(a.fields))

	if a.ptr == nil {
		strings := make([]string, len(a.fields))
		for i := range strings {
			strings[i] = writer.config.NullValue
		}

		return strings, make([]bool, len(a.fields)), make([]string, len(a.fields)), make([]string, len(a.fields)), make([]bool, len(a.fields)), make([]string, len(a.fields))
	}

	//var rowFieldType string
	//var ok bool
	sliceOffset := 0
	var sizeElem uintptr
	var feType reflect.Type
	var feKind reflect.Kind

	for i, field := range a.fields {
		fieldNames[i+sliceOffset] = field.name
		fType := field.xField.Type
		fKind := fType.Kind()

		if fKind == reflect.Ptr {
			fType = fType.Elem()
			fKind = fType.Kind()
		}

		if fKind == reflect.Struct {
			continue
		}

		if fKind == reflect.Slice {
			feType = fType.Elem()
			feKind = feType.Kind()

			switch feKind {
			case reflect.Bool:
				sizeElem = unsafe.Sizeof(*new(bool))
			case reflect.Int:
				sizeElem = unsafe.Sizeof(*new(int))
			case reflect.Int8:
				sizeElem = unsafe.Sizeof(*new(int8))
			case reflect.Int16:
				sizeElem = unsafe.Sizeof(*new(int16))
			case reflect.Int32:
				sizeElem = unsafe.Sizeof(*new(int32))
			case reflect.Int64:
				sizeElem = unsafe.Sizeof(*new(int64))
			case reflect.Uint:
				sizeElem = unsafe.Sizeof(*new(uint))
			case reflect.Uint8:
				sizeElem = unsafe.Sizeof(*new(uint8))
			case reflect.Uint16:
				sizeElem = unsafe.Sizeof(*new(uint16))
			case reflect.Uint32:
				sizeElem = unsafe.Sizeof(*new(uint32))
			case reflect.Uint64:
				sizeElem = unsafe.Sizeof(*new(uint64))
			case reflect.Float32:
				sizeElem = unsafe.Sizeof(*new(float32))
			case reflect.Float64:
				sizeElem = unsafe.Sizeof(*new(float64))
			case reflect.String:
				sizeElem = unsafe.Sizeof(*new(string))
			default:
				continue
			}

			sHdr := (*reflect.SliceHeader)(a.ptr /*unsafe.Pointer(&s)*/)
			//fmt.Printf("Len = %d, Cap = %d", sHdr.Len, sHdr.Cap)
			result = append(result, make([]string, sHdr.Len-1)...)
			wasStrings = append(wasStrings, make([]bool, sHdr.Len-1)...)
			types = append(types, make([]string, sHdr.Len-1)...)
			*headers = append(*headers, make([]string, sHdr.Len-1)...)
			*attrNames = append(*attrNames, make([]string, sHdr.Len-1)...)
			shouldWrite = append(shouldWrite, make([]bool, sHdr.Len-1)...)
			fieldNames = append(fieldNames, make([]string, sHdr.Len-1)...)

			//fmt.Printf("sizElem Z = %d\n", sizeElem)

			for j := 0; j < sHdr.Len; j++ {
				zPtr := unsafe.Pointer(sHdr.Data + uintptr(j)*sizeElem)
				var z string
				switch feKind {
				case reflect.Bool:
					z = fmt.Sprintf("%t", *(*bool)(zPtr))
				case reflect.Int:
					z = fmt.Sprintf("%d", *(*int)(zPtr))
				case reflect.Int8:
					z = fmt.Sprintf("%d", *(*int8)(zPtr))
				case reflect.Int16:
					z = fmt.Sprintf("%d", *(*int16)(zPtr))
				case reflect.Int32:
					z = fmt.Sprintf("%d", *(*int32)(zPtr))
				case reflect.Int64:
					z = fmt.Sprintf("%d", *(*int64)(zPtr))
				case reflect.Uint:
					z = fmt.Sprintf("%d", *(*uint)(zPtr))
				case reflect.Uint8:
					z = fmt.Sprintf("%d", *(*uint8)(zPtr))
				case reflect.Uint16:
					z = fmt.Sprintf("%d", *(*uint16)(zPtr))
				case reflect.Uint32:
					z = fmt.Sprintf("%d", *(*uint32)(zPtr))
				case reflect.Uint64:
					z = fmt.Sprintf("%d", *(*uint64)(zPtr))
				case reflect.Float32:
					z = fmt.Sprintf("%f", *(*float32)(zPtr))
				case reflect.Float64:
					z = fmt.Sprintf("%f", *(*float64)(zPtr))
				case reflect.String:
					z = *(*string)(zPtr)
				//	sizeElem = unsafe.Sizeof(*(new(string)))
				default:
					continue
				}

				result[i+sliceOffset], wasStrings[i+sliceOffset] = z, true
				types[i+sliceOffset] = "string" //TODO delete
				shouldWrite[i+sliceOffset] = true
				fieldNames[i+sliceOffset] = field.name
				if sliceOffset > 0 {
					(*headers)[i+sliceOffset] = (*headers)[i+sliceOffset-1]
					(*attrNames)[i+sliceOffset] = (*attrNames)[i+sliceOffset-1]
				}
				sliceOffset++

			}
			fmt.Println()
			//runtime.KeepAlive(s)
		} else {
			result[i+sliceOffset], wasStrings[i+sliceOffset] = field.stringifier(a.ptr) //TODO how to get real type here?
			types[i+sliceOffset] = field.xField.Type.String()
			shouldWrite[i+sliceOffset] = true
			fieldNames[i+sliceOffset] = field.name
		}

		//rowFieldType, ok = a.config.DataRowFieldTypes[types[i]]
		//if ok {
		//	DataRowFieldTypes[i] = rowFieldType
		//} else {
		//	DataRowFieldTypes[i] = "TYPE_ERR"
		//}

	}

	a.cache[a.ptr] = &stringified{
		values:            result,
		wasStrings:        wasStrings,
		types:             types,
		dataRowFieldTypes: dataRowFieldTypes,
		shouldWrite:       shouldWrite,
		fieldNames:        fieldNames,
	}

	return result, wasStrings, types, dataRowFieldTypes, shouldWrite, fieldNames
}
