package dgraph_ldap

import (
	"github.com/json-iterator/go"
	"github.com/v2pro/plz/reflect2"
	"unsafe"
)

var DGraphJsonMarshaller = jsoniter.Config{
	//IndentionStep:                 4,
	//MarshalFloatWith6Digits:       false,
	EscapeHTML:                    false,
	SortMapKeys:                   false,
	//UseNumber:                     true,
	//DisallowUnknownFields:         true,
	TagKey:                        "db",
	//OnlyTaggedField:               true,
	//ValidateJsonRawMessage:        true,
	//ObjectFieldMustBeSimpleString: true,
}.Froze()

var ClientJsonMarshaller = jsoniter.Config{
	EscapeHTML:                    false,
	SortMapKeys:                   false,
	TagKey:                        "client",
}.Froze()

type DecodeValueData struct {
	ptr unsafe.Pointer
	iter *jsoniter.Iterator
	defaultDecoder jsoniter.ValDecoder
	consumed bool
}

const dataConsumedErrorStr string = "Data is already consumed, you can only use one of [GetIterator, GetBytes, DecodeDefault]."

func (d *DecodeValueData) GetIterator() *jsoniter.Iterator {
	if d.consumed {
		d.iter.ReportError("DecodeValueData.GetIterator", dataConsumedErrorStr)
		return nil
	}
	d.consumed = true
	return d.iter
}

func (d *DecodeValueData) GetBytes() []byte {
	if d.consumed {
		d.iter.ReportError("DecodeValueData.GetIterator", dataConsumedErrorStr)
		return nil
	}
	d.consumed = true
	return d.iter.SkipAndReturnBytes()
}

func (d *DecodeValueData) GetAPI() jsoniter.API {
	api, _ := d.iter.Pool().(jsoniter.API)
	return api
}

func (d *DecodeValueData) DecodeDefault() {
	if d.consumed {
		d.iter.ReportError("DecodeValueData.GetIterator", dataConsumedErrorStr)
		return
	}
	d.consumed = true
	d.defaultDecoder.Decode(d.ptr, d.iter)
}

type LazyInitStruct interface {
	InitBeforeUnmarshal() error
}

type LazyInitStructDecorator struct {
	valType reflect2.Type
	decoder jsoniter.ValDecoder
}

func (d *LazyInitStructDecorator) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valType := d.valType
	obj := valType.UnsafeIndirect(unsafe.Pointer(&ptr))  // TODO: Why unsafe.Pointer(&ptr) ???
	unmarshaler := obj.(LazyInitStruct)
	err := unmarshaler.InitBeforeUnmarshal()
	if err != nil {
		iter.ReportError("LazyInitStructDecorator", err.Error())
	}
	d.decoder.Decode(ptr, iter)
}

type CustomDBUnmarshal interface {
	UnmarshalFromDB(valueData *DecodeValueData) error
}

type CustomDBUnmarshalDecoderDecorator struct {
	valType reflect2.Type
	defaultDecoder jsoniter.ValDecoder
}

func (d *CustomDBUnmarshalDecoderDecorator) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	valType := d.valType
	obj := valType.UnsafeIndirect(unsafe.Pointer(&ptr))  // TODO: Why unsafe.Pointer(&ptr) ???
	unmarshaler := obj.(CustomDBUnmarshal)
	err := unmarshaler.UnmarshalFromDB(&DecodeValueData{ptr, iter, d.defaultDecoder, false})
	if err != nil {
		iter.ReportError("CustomDBUnmarshalDecoderDecorator", err.Error())
	}
}

type CustomDBMarshal interface {
	MarshalToDB(api jsoniter.API) ([]byte, error)
}

type CustomClientUnmarshal interface {
	UnmarshalFromClient(b []byte, api jsoniter.API) error
}

type CustomClientMarshal interface {
	MarshalToClient(api jsoniter.API) ([]byte, error)
}

type dGraphJsonMarshallerExtension struct {
	jsoniter.DummyExtension
}

var customDBUnmarshalType = reflect2.TypeOfPtr((*CustomDBUnmarshal)(nil)).Elem()
var lazyInitStructType = reflect2.TypeOfPtr((*LazyInitStruct)(nil)).Elem()

func (extension *dGraphJsonMarshallerExtension) DecorateDecoder(typ reflect2.Type, decoder jsoniter.ValDecoder) jsoniter.ValDecoder {
	ptrType := reflect2.PtrTo(typ)
	dec := decoder
	if ptrType.Implements(customDBUnmarshalType) {
		dec = &CustomDBUnmarshalDecoderDecorator{ptrType, dec}
	}
	if ptrType.Implements(lazyInitStructType) {
		dec = &LazyInitStructDecorator{ptrType, dec}
	}
	return dec
}

func init() {
	DGraphJsonMarshaller.RegisterExtension(&dGraphJsonMarshallerExtension{})
}

