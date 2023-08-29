package xmlify

import (
	"github.com/viant/xunsafe"
	"reflect"
)

func (w *writer) writeTabularAllObjects(acc *Accessor, parentLevel bool) {
	var xType *xunsafe.Type

	w.buffer.writeString("<" + acc.config.RootTag + ">")
	defer w.buffer.writeString(acc.config.NewLineSeparator + "</" + acc.config.RootTag + ">")

	w.writeTabularHeaderElement(acc.Headers())

	w.buffer.writeString(acc.config.NewLineSeparator + "<" + acc.config.DataTag + ">")
	defer w.buffer.writeString(acc.config.NewLineSeparator + "</" + acc.config.DataTag + ">")

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
			result, _, _, dataRowFieldTypes := acc.stringifyFields(w)
			w.writeTabularElement(result, dataRowFieldTypes)
		}
	}
}

func (w *writer) writeTabularHeaderElement(data []string, headerRowFieldTypes []string) {
	w.buffer.writeString(w.config.NewLineSeparator + "<" + w.config.HeaderTag + ">")
	w.writeTabularHeaderValue(data, headerRowFieldTypes)
	w.buffer.writeString(w.config.NewLineSeparator + "</" + w.config.HeaderTag + ">")
	w.writtenObject = true
}

func (w *writer) writeTabularHeaderValue(values []string, headerRowFieldTypes []string) {
	if len(values) == 0 {
		return
	}

	for j := 0; j < len(values); j++ {
		asString := EscapeSpecialChars(values[j], w.config)
		asString = w.config.NewLineSeparator + "<" + w.config.HeaderRowTag + " " + w.config.HeaderRowFieldAttr + "=" +
			w.config.EncloseBy + asString + w.config.EncloseBy +
			" " + w.config.HeaderRowFieldTypeAttr + "=" + "\"" + headerRowFieldTypes[j] + "\"" +
			"/>"

		w.buffer.writeString(asString)
	}
}

func (w *writer) writeTabularElement(data []string, dataRowFieldTypes []string) {
	if len(data) == 0 {
		return
	}

	w.buffer.writeString(w.config.NewLineSeparator + "<" + w.config.DataRowTag + ">")
	w.writeTabularValue(data, dataRowFieldTypes)
	w.buffer.writeString(w.config.NewLineSeparator + "</" + w.config.DataRowTag + ">")
	w.writtenObject = true
}

func (w *writer) writeTabularValue(data []string, dataRowFieldTypes []string) {
	for j := 0; j < len(data); j++ {
		asString := EscapeSpecialChars(data[j], w.config)

		if asString == w.config.escapedNullValue {
			if w.config.TabularNullValue == "" {
				asString = w.config.NewLineSeparator + "<" + w.config.DataRowFieldTag + "/>"
			} else {
				asString = w.config.NewLineSeparator + "<" + w.config.DataRowFieldTag +
					" " + w.config.TabularNullValue +
					"/>"
			}

			w.buffer.writeString(asString)
			continue
		}

		dataType := dataRowFieldTypes[j]

		if dataType != "string" {
			asString = w.config.NewLineSeparator + "<" + w.config.DataRowFieldTag +
				" " + dataType + "=" +
				w.config.EncloseBy + asString + w.config.EncloseBy +
				"/>"
		} else {
			asString = w.config.NewLineSeparator + "<" + w.config.DataRowFieldTag + ">" +
				asString +
				"</" + w.config.DataRowFieldTag + ">"
		}
		w.buffer.writeString(asString)
	}
}
