package dgraph_ldap

//import "dgraph_ldap"

import (
	"testing"
	"time"
)

func TestUnmarshalNode(t *testing.T) {
	t.Log("Testing unmarshal document node ...")
	data := []byte(`{
			"uid": 1,
			"gn": "aaa",
			"t": 0,
			"n": "DocNode"
		}`)
	node, err := UnmarshalNode(data)
	doc := node.(*Document)
	t.Logf("node: %v", node)
	t.Logf("node: %v", doc.Name)
	t.Logf("err: %v", err)
}

func TestMarshalNode(t *testing.T) {
	t.Log("Testing marshal document node ...")
	node := &Document{
		userNode: userNode{
			nodeBase: nodeBase{
				UID: 0x1,
				GlobalName: "aaa",
				Type: TypeDocument,
				Name: "DocNode",
				TreeName: "DocNode",
				CreateTime: time.Now(),
				ModifyTime: time.Now(),
			},
		},
	}
	data, err := MarshalNode(node)
	t.Logf("json: %s", data)
	t.Logf("err: %v", err)
}
