package main

import (
	"fmt"
	"github.com/json-iterator/go"
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

//func (s *s0) Test()  {
//	//fmt.Println(s.Type)
//	fmt.Println(s)
//}

type s1 struct {
	s0
	//Type uint
	BBB string `cl:"bbb" dg:"b"`
}

func main() {
	var v0 s1
	sss := []byte(`{"a":"2333", "b":"hhhhhh", "c":{"d": 233, "e": "213"}}`)
	fmt.Println(DGraphJsonMarshaller.Get(sss, "c").ToString())
	err := DGraphJsonMarshaller.Unmarshal(sss, &v0)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(v0)

	fmt.Println(ClientJsonMarshaller.MarshalToString(v0))

	//(&s0{AAA: "hhh", Type:0}).Test()
	//(&s1{s0:s0{AAA: "233", Type:0}, Type:1}).Test()
}

