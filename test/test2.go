package main

import (
	"fmt"
	"github.com/json-iterator/go"
	//"github.com/v2pro/plz/reflect2"
	"unsafe"
	//"encoding/json"
	"github.com/v2pro/plz/reflect2"
	//"reflect"
	//"strconv"
	//"encoding/json"
)

var DGraphJsonMarshaller = jsoniter.Config{
	//IndentionStep:                 4,
	//MarshalFloatWith6Digits:       false,
	EscapeHTML:                    false,
	SortMapKeys:                   false,
	//UseNumber:                     true,
	//DisallowUnknownFields:         true,
	TagKey:                        "dg",
	//OnlyTaggedField:               true,
	//ValidateJsonRawMessage:        true,
	//ObjectFieldMustBeSimpleString: true,
}.Froze()

var ClientJsonMarshaller = jsoniter.Config{
	EscapeHTML:                    false,
	SortMapKeys:                   false,
	TagKey:                        "cl",
}.Froze()

type s0 struct {
	AAA string `cl:"aaa" dg:"a"`
	//Type uint
}

func (s *s0) Test()  {
	//fmt.Println(s.Type)
	fmt.Println("s0")
}

func (s *s1) UnmarshalJSON(b []byte, api jsoniter.API) error  {
	//fmt.Println(s.Type)
	fmt.Println(api == DGraphJsonMarshaller)
	fmt.Println(api == ClientJsonMarshaller)
	fmt.Printf("%s\n", b)
	return nil
}

type s1 struct {
	s0
	//Type uint
	BBB string `cl:"bbb" dg:"b"`
}

func (s *s1) Test()  {
	//fmt.Println(s.Type)
	fmt.Println("s1")
	s.s0.Test()
}

type CustomJsoniterType interface {
	//MarshalJSON(api jsoniter.API) ([]byte, error)
	UnmarshalJSON(b []byte, api jsoniter.API) error
}

type CustomJsoniterTypeDecoder struct {
	valType reflect2.Type
}

func (decoder *CustomJsoniterTypeDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	if api, ok := iter.Pool().(jsoniter.API); ok {
		valType := decoder.valType
		obj := valType.UnsafeIndirect(ptr)
		unmarshaler := obj.(CustomJsoniterType)
		//unmarshaler := *((*CustomJsoniterType)(ptr))
		bytes := iter.SkipAndReturnBytes()
		err := unmarshaler.UnmarshalJSON(bytes, api)
		if err != nil {
			iter.ReportError("CustomJsoniterTypeDecoder", err.Error())
		}
	} else {
		iter.ReportError("CustomJsoniterTypeDecoder", "Can't find jsoniter API.")
	}
}

type testMapKeyExtension struct {
	jsoniter.DummyExtension
}

var unmarshalerType = reflect2.TypeOfPtr((*CustomJsoniterType)(nil)).Elem()

func (extension *testMapKeyExtension) CreateDecoder(typ reflect2.Type) jsoniter.ValDecoder {
	ptrType := reflect2.PtrTo(typ)
	if ptrType.Implements(unmarshalerType) {
		return &CustomJsoniterTypeDecoder{ptrType}

	}
	return nil
}

func init() {
	//DGraphJsonMarshaller.RegisterExtension(&testMapKeyExtension{})
}

func main() {
	var v0 s1
	sss := []byte(`{"i": 213213, "a":"2333", "b":"hhhhhh", "c":{"d": 233, "e": "213"}}`)
	iii := DGraphJsonMarshaller.Get(sss, "i")
	fmt.Println(iii.ValueType())
	fmt.Println(iii.ToUint())
	err := DGraphJsonMarshaller.Unmarshal(sss, &v0)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(v0)

	fmt.Println(ClientJsonMarshaller.MarshalToString(v0))

	v0.Test()

	var str0 []byte = []byte("aaabbbccc")
	str1 := str0
	fmt.Println(str0)
	str1[1] = 0
	fmt.Println(str0)

	str2 := "aaa"
	var str3 *string = &str2
	rrr0, _ := DGraphJsonMarshaller.MarshalToString(str3)
	fmt.Println("rrr0", rrr0)
	rrr1, _ := DGraphJsonMarshaller.MarshalToString(str1)
	fmt.Println("rrr1", rrr1)

	//(&s0{AAA: "hhh", Type:0}).Test()
	//(&s1{s0:s0{AAA: "233", Type:0}, Type:1}).Test()
}

