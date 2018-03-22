package dgraph_ldap

//import "dgraph_ldap"

import (
	"testing"
	"github.com/json-iterator/go"
)

func TestUnmarshalNode(t *testing.T) {
	t.Log("Testing unmarshal document node ...")
	data := []byte(`{
			"uid": 1,
			"gn": "aaa",
			"t": 0,
			"n": "DocNode",
			"d": [
				{
					"d|f": "fff",
					"rd": "{\"aaa\": \"bbb\"}",
					"dt": 0
				}
			]
		}`)
	node, err := UnmarshalNode(data)
	doc := node.(*Document)
	t.Logf("node: %v", node)
	t.Logf("node: %v", doc.Name)
	ddd, _ := doc.Data.Get("fff")
	t.Logf("node: %v", ddd.Data.(jsoniter.Any).Get("aaa"))
	t.Logf("err: %v", err)
}

//func TestMarshalNode(t *testing.T) {
//	t.Log("Testing marshal document node ...")
//	node := &Document{
//		userNode: userNode{
//			nodeBase: nodeBase{
//				UID: 0x1,
//				GlobalName: "aaa",
//				Type: TypeDocument,
//				Name: "DocNode",
//				TreeName: "DocNode",
//				CreateTime: time.Now(),
//				ModifyTime: time.Now(),
//			},
//		},
//	}
//	data, err := MarshalNode(node)
//	t.Logf("json: %s", data)
//	t.Logf("err: %v", err)
//}
