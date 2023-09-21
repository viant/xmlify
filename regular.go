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

	areAttributes, attrNames, containAttr, allAreAttr := acc.attributes()

	customMarshaling := false
	// check custom marshaling
	if acc.field != nil && acc._parent != nil && acc._parent.ptr != nil {
		value := acc.field.Value(acc.ptr)
		if _, ok := value.(XMLMarhsaler); ok {
			customMarshaling = true
		}
	}

	// start
	if !customMarshaling && !(parentLevel && omitRootElement) && fieldKind != reflect.Slice && fieldName != "" {
		if containAttr {
			w.buffer.writeString(w.config.NewLineSeparator + "<" + fieldName)
		} else {
			w.buffer.writeString(w.config.NewLineSeparator + "<" + fieldName + ">")
		}
	}

	// custom marshaling
	if customMarshaling {
		w.buffer.writeString(w.config.NewLineSeparator)
		value := acc.field.Value(acc._parent.ptr)
		custom, ok := value.(XMLMarhsaler)
		if ok {
			data, err := custom.MarshalXML()
			if err != nil {
				w.buffer.writeString(err.Error()) //TODO error handling
				fmt.Printf("custom marshaling error: %s\n", err.Error())
				return
			}
			w.buffer.writeString(string(data))
		}
		return
	}

	// finish
	defer func() {
		if !customMarshaling && !(parentLevel && omitRootElement) && fieldKind != reflect.Slice && fieldName != "" {
			if !allAreAttr {
				w.buffer.writeString(w.config.NewLineSeparator + "</" + fieldName + ">")
			}
		}
	}()

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

			result, wasStrings, types, _, shouldWrite, fieldNames := acc.stringifyRegularFields(w, &headers, &attrNames, &areAttributes)

			//len(field)< len(result) // when one of field is a slice !!!

			if fieldKind == reflect.Slice || acc.slice != nil {
				w.buffer.writeString(w.config.NewLineSeparator + "<" + rowFieldName)
				if !containAttr {
					w.buffer.writeString(">")
				}
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
			for j, field := range acc.fields {
				start = lowerBound[j]
				end = upperBound[j]

				if areAttributes[start] && shouldWrite[start] { // TODO handle slices ?
					w.writeRegularObjectAttr(result[start:end], wasStrings[start:end], types[start:end], []string{}, attrNames[start:end], shouldWrite[start:end], field)
					// TODO close tag when all fields are attributes?
				}
			}
			if containAttr {
				if allAreAttr {
					w.buffer.writeString("/>")
				} else {
					w.buffer.writeString(">")
				}
			}

			for j, field := range acc.fields {
				start = lowerBound[j]
				end = upperBound[j]

				if !areAttributes[start] && shouldWrite[start] { // whole result[start:end] represents one field (can be a slice)
					w.writeRegularElement(result[start:end], headers[start:end], shouldWrite[start:end], field)
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

			if (fieldKind == reflect.Slice || acc.slice != nil) && !allAreAttr {
				w.buffer.writeString(w.config.NewLineSeparator + "</" + rowFieldName + ">")
			}

		} // ~ has loop
	} // ~ main for loop
}

// TODO return ommited slice
func (w *writer) WriteRegularObjectAttr(values []string, wasString []bool, types []string, dataRowFieldTypes []string, headers []string, shouldWrite []bool, field *Field) []bool {
	omited := make([]bool, len(values))

	if len(values) == 0 {
		return omited
	}

	currentAttr := ""
	lastAttr := ""

	for j := 0; j < len(values); j++ {
		if !shouldWrite[j] {
			continue
		}

		asString := EscapeSpecialChars(values[j], w.config)

		if field.tag.OmitEmpty && asString == w.config.escapedNullValue {
			omited[j] = true
			continue
		}

		currentAttr = headers[j]

		if currentAttr != lastAttr && lastAttr != "" { // closing last attr
			w.buffer.writeString("\"")
		}

		if currentAttr != lastAttr { // opening new attr
			w.buffer.writeString(" " + currentAttr + "=" + "\"")
		}

		if currentAttr == lastAttr && lastAttr != "" { //separate another value for current attr
			w.buffer.writeString(",") // TODO ADD attr separator into config?
		}

		lastAttr = currentAttr

		if asString == w.config.escapedNullValue {
			asString = w.config.RegularNullValue
			w.buffer.writeString(asString)
			continue
		}

		w.buffer.writeString(asString)
	}

	//closing last attr
	if lastAttr != "" {
		w.buffer.writeString("\"")
	}

	//w.buffer.writeString("]")
	return omited
}

// TODO pass *Field
func (w *writer) writeRegularObjectAttr(data []string, wasStrings []bool, types, dataRowFieldTypes, headers []string, shouldWrite []bool, field *Field) {
	if w.writtenObject {
		//w.writeObjectSeparator()
	} else {
		w.buffer.writeString(w.beforeFirst) // TODO w.beforeFirst [[??
	}

	//WriteObject(w.buffer, w.config, data, wasStrings)
	w.WriteRegularObjectAttr(data, wasStrings, types, dataRowFieldTypes, headers, shouldWrite, field)
	w.writtenObject = true
}

func (w *writer) writeRegularElement(values []string, headers []string, shouldWrite []bool, field *Field) {
	if len(values) == 0 {
		return
	}

	for j := 0; j < len(values); j++ {
		if !shouldWrite[j] {
			continue
		}

		asString := EscapeSpecialChars(values[j], w.config)
		tagStart := "<" + headers[j] + ">"
		tagEnd := "</" + headers[j] + ">"

		if asString == w.config.escapedNullValue {
			if field.tag.OmitEmpty {
				continue
			}

			if w.config.RegularNullValue == "" {
				tagStart = "<" + headers[j]
				tagEnd = "/>"
				asString = w.buildElement(tagStart, "", tagEnd, field.tag.OmitTagName)
			} else {
				tagStart = "<" + headers[j] + " "
				tagEnd = "/>"
				asString = w.buildElement(tagStart, w.config.RegularNullValue, tagEnd, field.tag.OmitTagName)
			}
			w.buffer.writeString(asString)
			continue
		}

		asString = w.buildElement(tagStart, asString, tagEnd, field.tag.OmitTagName)
		w.buffer.writeString(asString)
	}
	w.writtenObject = true
}

func (w *writer) buildElement(start, value, end string, omitTagName bool) string {
	if omitTagName {
		return value
	}

	return w.config.NewLineSeparator + start + value + end
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

func (a *Accessor) stringifyRegularFields(writer *writer, headers *[]string, attrNames *[]string, areAttributes *[]bool) ([]string, []bool, []string, []string, []bool, []string) {
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

	currentCounter := 0
	var sizeElem uintptr
	var feType reflect.Type
	var feKind reflect.Kind

	// TODO MFI
	var newHeaders = make([]string, len(*headers))
	var newAttrNames = make([]string, len(*attrNames))
	var newAreAttributes = make([]bool, len(*areAttributes))

	for fieldNumber, field := range a.fields {
		if currentCounter >= len(fieldNames) {
			// TODO MFI
		}

		fieldNames[currentCounter] = field.name
		newHeaders[currentCounter] = (*headers)[fieldNumber]
		newAttrNames[currentCounter] = (*attrNames)[fieldNumber]
		newAreAttributes[currentCounter] = (*areAttributes)[fieldNumber]

		fType := field.xField.Type
		fKind := fType.Kind()

		if fKind == reflect.Ptr {
			fType = fType.Elem()
			fKind = fType.Kind()
		}

		if fKind == reflect.Struct {
			currentCounter++
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
				currentCounter++
				continue
			}

			sHdr := (*reflect.SliceHeader)(unsafe.Add(a.ptr, field.xField.Offset) /*unsafe.Pointer(&s)*/)

			if sHdr.Len == 0 {
				result[currentCounter] = a.config.NullValue
				wasStrings[currentCounter] = true
				types[currentCounter] = field.xField.Type.String()
				shouldWrite[currentCounter] = true
				fieldNames[currentCounter] = field.name
				currentCounter++
				continue
			}

			if sHdr.Len > 1 {
				result = append(result, make([]string, sHdr.Len-1)...)
				wasStrings = append(wasStrings, make([]bool, sHdr.Len-1)...)
				types = append(types, make([]string, sHdr.Len-1)...)
				newHeaders = append(newHeaders, make([]string, sHdr.Len-1)...)
				newAttrNames = append(newAttrNames, make([]string, sHdr.Len-1)...)
				newAreAttributes = append(newAreAttributes, make([]bool, sHdr.Len-1)...)
				shouldWrite = append(shouldWrite, make([]bool, sHdr.Len-1)...)
				fieldNames = append(fieldNames, make([]string, sHdr.Len-1)...)
			}

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
				default:
					continue
				}

				fieldNames[currentCounter] = field.name
				result[currentCounter], wasStrings[currentCounter] = z, true
				types[currentCounter] = "string" //TODO delete
				shouldWrite[currentCounter] = true

				if j > 0 {
					newHeaders[currentCounter] = (*headers)[fieldNumber]
					newAttrNames[currentCounter] = (*attrNames)[fieldNumber]
					newAreAttributes[currentCounter] = (*areAttributes)[fieldNumber]
				}
				currentCounter++

			}
		} else {
			//			result[currentCounter], wasStrings[currentCounter] = field.stringifier(a.slicePtr)
			result[currentCounter], wasStrings[currentCounter] = field.stringifier(a.ptr)
			types[currentCounter] = field.xField.Type.String()
			shouldWrite[currentCounter] = true
			fieldNames[currentCounter] = field.name
			currentCounter++
		}
	}

	a.cache[a.ptr] = &stringified{
		values:            result,
		wasStrings:        wasStrings,
		types:             types,
		dataRowFieldTypes: dataRowFieldTypes,
		shouldWrite:       shouldWrite,
		fieldNames:        fieldNames,
	}

	if len(*headers) != len(newHeaders) {
		*headers = append(*headers, newHeaders[len(*headers):]...)
		*attrNames = append(*attrNames, newAttrNames[len(*attrNames):]...)
		*areAttributes = append(*areAttributes, newAreAttributes[len(*areAttributes):]...)
		copy(*headers, newHeaders)
		copy(*attrNames, newAttrNames)
		copy(*areAttributes, newAreAttributes)
	}

	return result, wasStrings, types, dataRowFieldTypes, shouldWrite, fieldNames
}
