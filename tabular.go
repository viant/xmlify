package xmlify

import (
	"github.com/viant/xunsafe"
	"reflect"
)

func (w *writer) writeTabularAllObjects(acc *Accessor, parentLevel bool) {

	if w.config.style != tabularStyle {
		//TODO MFI
	}
	/*
		w.buffer.writeString("<" + w.config.rootTag + ">") // TODO MFI
		defer w.buffer.writeString(w.config.newLine + "</" + w.config.rootTag + ">")
	*/
	//TODO use acc config not writer config (can be the same but not always)
	w.buffer.writeString("<" + acc.config.rootTag + ">") // TODO MFI
	defer w.buffer.writeString(acc.config.newLine + "</" + acc.config.rootTag + ">")

	w.writeHeadersIfNeeded(acc.Headers())

	var xType *xunsafe.Type

	//fmt.Printf("\n*** SIZE = %d ***\n", w.size) // TODO check sizes

	/*
		w.buffer.writeString(w.config.newLine + "<" + w.config.dataTag + ">")
		defer w.buffer.writeString(w.config.newLine + "</" + w.config.dataTag + ">")
	*/
	//TODO use acc config not writer config
	w.buffer.writeString(acc.config.newLine + "<" + acc.config.dataTag + ">")
	defer w.buffer.writeString(acc.config.newLine + "</" + acc.config.dataTag + ">")

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
			//w.buffer.writeString(w.config.FieldSeparator)
			//w.buffer.writeString("#")
			//w.buffer.writeString(w.config.newLine) // TODO MFI use for test formatting only

			//w.buffer.writeString("[")
			//w.buffer.writeString(w.config.newLine + "<" + w.config.dataRowTag + ">")

			result, wasStrings, types, dataRowFieldTypes := acc.stringifyFields(w)
			w.writeTabularObject(result, wasStrings, types, dataRowFieldTypes)
			//w.writeObject(result, wasStrings)

			//for _, child := range acc.children {
			//	w.buffer.writeString(w.config.FieldSeparator)
			//	//w.buffer.writeString("$")
			//
			//	_, childSize := child.values()
			//
			//	if childSize > 0 {
			//		tmpSize := w.size
			//		w.size = childSize
			//		w.writeTabularAllObjects(child, false)
			//		w.size = tmpSize
			//	}
			//
			//	if childSize == 0 {
			//		w.buffer.writeString("null")
			//	}
			//}

			//w.buffer.writeString("]")
			//w.buffer.writeString("</" + w.config.dataRowTag + ">")
		}
	}
}

func WriteTabularObject(writer *Buffer, config *Config, values []string, wasString []bool, types []string, dataRowFieldTypes []string) {
	if len(values) == 0 {
		return
	}

	writer.writeString(config.newLine + "<" + config.dataRowTag + ">")
	escapedNullValue := EscapeSpecialChars(config.NullValue, config)

	for j := 0; j < len(values); j++ {

		//if j != 0 {
		//	writer.writeString(config.FieldSeparator) //TODO MFI field separtator always
		//}

		asString := EscapeSpecialChars(values[j], config) //TODO MFI escaping

		if asString == escapedNullValue {
			asString = config.nullValueTODO
			asString = config.newLine + "<" + config.dataRowFieldTag +
				" " + asString +
				"/>"
			writer.writeString(asString)
			continue
		}

		// if wasString[j] { every value has to be string

		dataType := dataRowFieldTypes[j]

		///////////
		if dataType != "string" {
			asString = config.newLine + "<" + config.dataRowFieldTag +
				" " + dataType + "=" +
				config.EncloseBy + asString + config.EncloseBy +
				"/>"
		} else {
			asString = config.newLine + "<" + config.dataRowFieldTag + ">" +
				asString +
				"</" + config.dataRowFieldTag + ">"
		}
		///////////

		writer.writeString(asString)
	}

	writer.writeString(config.newLine + "</" + config.dataRowTag + ">")
	//writer.writeString("]")

}

func WriteTabularHeaderObject(writer *Buffer, config *Config, values []string, wasString []bool, headerRowFieldTypes []string) {
	if len(values) == 0 {
		return
	}

	//writer.writeString("[")

	for j := 0; j < len(values); j++ {
		//if j != 0 {
		//	writer.writeString(config.FieldSeparator) //TODO MFI field separtator always
		//}

		asString := EscapeSpecialChars(values[j], config) //TODO MFI escaping
		if wasString[j] {
			asString = config.newLine + "<" + config.headerRowTag + " " + config.headerRowFieldAttr + "=" +
				config.EncloseBy + asString + config.EncloseBy +
				" " + config.headerRowFieldTypeAttr + "=" + "\"" + headerRowFieldTypes[j] + "\"" +
				"/" + /*config.headerRowTag +*/ ">"
			//TODO additional config
		}

		writer.writeString(asString)
	}

	//writer.writeString("]")

}

func (w *writer) writeTabularHeaderObject(data []string, wasStrings []bool, headerRowFieldTypes []string) {
	if w.writtenObject {
		//w.writeObjectSeparator()
	} else {
		w.buffer.writeString(w.beforeFirst) // TODO w.beforeFirst [[??
	}

	w.buffer.writeString(w.config.newLine + "<" + w.config.headerTag + ">")
	WriteTabularHeaderObject(w.buffer, w.config, data, wasStrings, headerRowFieldTypes)
	w.buffer.writeString(w.config.newLine + "</" + w.config.headerTag + ">")
	w.writtenObject = true
}

func (w *writer) writeTabularObject(data []string, wasStrings []bool, types, dataRowFieldTypes []string) {
	if w.writtenObject {
		//w.writeObjectSeparator()
	} else {
		w.buffer.writeString(w.beforeFirst) // TODO w.beforeFirst [[??
	}

	//WriteObject(w.buffer, w.config, data, wasStrings)
	WriteTabularObject(w.buffer, w.config, data, wasStrings, types, dataRowFieldTypes)
	w.writtenObject = true
}
