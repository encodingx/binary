package structs

import (
	"reflect"

	"github.com/joel-ling/go-bitfields/pkg/formats"
)

type FormatStructParser interface {
	ParseFormatStruct(interface{}) (formats.Format, error)
}

type formatStructParser struct {
	formatCache map[reflect.Type]formats.Format
}

func NewFormatStructParser() (parser *formatStructParser) {
	// Return a default implementation of interface FormatStructParser.

	parser = &formatStructParser{
		formatCache: make(map[reflect.Type]formats.Format),
	}

	return
}

func (p *formatStructParser) ParseFormatStruct(v interface{}) (
	format formats.Format, e error,
) {
	// Implement interface FormatStructParser.

	var (
		formatInCache  bool
		typeReflection reflect.Type
	)

	typeReflection = reflect.TypeOf(v)

	format, formatInCache = p.formatCache[typeReflection]

	if formatInCache {
		return
	}

	format, e = formats.NewFormatFromType(typeReflection)
	if e != nil {
		return
	}

	p.formatCache[typeReflection] = format

	return
}
