package xmlify

import (
	"github.com/viant/xunsafe"
	"reflect"
)

func (w *writer) writeTabularAllObjects(acc *Accessor, parentLevel bool) {

	if w.config.Style != tabularStyle {
		//TODO MFI
	}
	/*
		w.buffer.writeString("<" + w.config.RootTag + ">") // TODO MFI
		defer w.buffer.writeString(w.config.NewLine + "</" + w.config.RootTag + ">")
	*/
	//TODO use acc config not writer config (can be the same but not always)
	w.buffer.writeString("<" + acc.config.RootTag + ">") // TODO MFI
	defer w.buffer.writeString(acc.config.NewLine + "</" + acc.config.RootTag + ">")

	w.writeHeadersIfNeeded(acc.Headers())

	var xType *xunsafe.Type

	//fmt.Printf("\n*** SIZE = %d ***\n", w.size) // TODO check sizes

	/*
		w.buffer.writeString(w.config.NewLine + "<" + w.config.DataTag + ">")
		defer w.buffer.writeString(w.config.NewLine + "</" + w.config.DataTag + ">")
	*/
	//TODO use acc config not writer config
	w.buffer.writeString(acc.config.NewLine + "<" + acc.config.DataTag + ">")
	defer w.buffer.writeString(acc.config.NewLine + "</" + acc.config.DataTag + ">")

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
			//w.buffer.writeString(w.config.NewLine) // TODO MFI use for test formatting only

			//w.buffer.writeString("[")
			//w.buffer.writeString(w.config.NewLine + "<" + w.config.DataRowTag + ">")

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
			//w.buffer.writeString("</" + w.config.DataRowTag + ">")
		}
	}
}

func WriteTabularObject(writer *Buffer, config *Config, values []string, wasString []bool, types []string, dataRowFieldTypes []string) {
	if len(values) == 0 {
		return
	}

	writer.writeString(config.NewLine + "<" + config.DataRowTag + ">")
	escapedNullValue := EscapeSpecialChars(config.NullValue, config)

	for j := 0; j < len(values); j++ {

		//if j != 0 {
		//	writer.writeString(config.FieldSeparator) //TODO MFI field separtator always
		//}

		asString := EscapeSpecialChars(values[j], config) //TODO MFI escaping

		if asString == escapedNullValue {
			asString = config.NullValueTODO
			asString = config.NewLine + "<" + config.DataRowFieldTag +
				" " + asString +
				"/>"
			writer.writeString(asString)
			continue
		}

		// if wasString[j] { every value has to be string

		dataType := dataRowFieldTypes[j]

		///////////
		if dataType != "string" {
			asString = config.NewLine + "<" + config.DataRowFieldTag +
				" " + dataType + "=" +
				config.EncloseBy + asString + config.EncloseBy +
				"/>"
		} else {
			asString = config.NewLine + "<" + config.DataRowFieldTag + ">" +
				asString +
				"</" + config.DataRowFieldTag + ">"
		}
		///////////

		writer.writeString(asString)
	}

	writer.writeString(config.NewLine + "</" + config.DataRowTag + ">")
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
			asString = config.NewLine + "<" + config.HeaderRowTag + " " + config.HeaderRowFieldAttr + "=" +
				config.EncloseBy + asString + config.EncloseBy +
				" " + config.HeaderRowFieldTypeAttr + "=" + "\"" + headerRowFieldTypes[j] + "\"" +
				"/" + /*config.HeaderRowTag +*/ ">"
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

	w.buffer.writeString(w.config.NewLine + "<" + w.config.HeaderTag + ">")
	WriteTabularHeaderObject(w.buffer, w.config, data, wasStrings, headerRowFieldTypes)
	w.buffer.writeString(w.config.NewLine + "</" + w.config.HeaderTag + ">")
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
