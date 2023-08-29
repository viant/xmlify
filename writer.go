package xmlify

import (
	io2 "github.com/viant/sqlx/io"
	"github.com/viant/xunsafe"
)

type writer struct {
	beforeFirst   string
	writtenObject bool
	dereferencer  *xunsafe.Type
	buffer        *Buffer
	config        *Config
	accessor      *Accessor
	valueAt       io2.ValueAccessor
	size          int
}

func newWriter(accessor *Accessor, config *Config, buffer *Buffer, dereferencer *xunsafe.Type, valueAt io2.ValueAccessor, size int, beforeFirst string) *writer {
	return &writer{
		dereferencer: dereferencer,
		buffer:       buffer,
		config:       config,
		accessor:     accessor,
		valueAt:      valueAt,
		size:         size,
		beforeFirst:  beforeFirst,
	}
}

func (w *writer) writeObjectSeparator() {
	w.buffer.writeString(w.config.ObjectSeparator)
}
