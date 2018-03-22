package dgraph_ldap

import (
	"time"
	"github.com/json-iterator/go"
	"fmt"
)

type NodeType uint

const (
	TypeDocument NodeType = iota
	TypeFile
	TypeUser
	TypeGroup
	TypeTag
)

type Node interface {
	GetChildren() error
}

type nodeBase struct {
	// DGraph uid
	UID uint64 `db:"uid" client:"uid,omitempty"`

	// global name, xid?
	GlobalName string `db:"gn,omitempty" client:"global_name,omitempty" schema:"string @index(hash)"`

	// node type
	Type NodeType `db:"t" client:"type" schema:"int @index(int)"`

	// name:
	Name string `db:"n" client:"name" schema:"string @index(hash, exact, fulltext) @count"`

	// DGraph facets: ch|n  TODO 能不能和name合并?
	TreeName string `db:"ch|n,omitempty"`

	// children:
	Children []Node `db:"ch,omitempty" client:"children,omitempty" schema:"uid @reverse"`

	// create time:
	CreateTime time.Time `db:"ct" client:"create_time" schema:"dateTime @index(day)"`

	// modify time:
	ModifyTime time.Time `db:"mt" client:"modify_time" schema:"dateTime @index(hour)"`
}

type assetInfoNode struct {
	nodeBase

	// user creator:
	Creator uint64 `db:"uc" client:"creator" schema:"uid"`
	// owner:
	Owner uint64 `db:"uo" client:"owner" schema:"uid"`
	// group owner:
	Group uint64 `db:"go" client:"group" schema:"uid"`
	// permissions:
	// roles = owner, group, user, tmp_user, guest
	// permissions:
	// r: read data, w: write data,     // data: data, tags, d-xxx
	// l: list children, c: append/remove children,
	// m: modify name/parrent, d: delete,
	// p: modify permission, owner, group
	//
	// excute !!!!!!!!!!!!!
	//	EX|210clwr|43210clwr|765pdm43210clwr|765pdm43210clwr|765pdm43210clwr    // 63 bits -> int64
	// E: has extra perm
	// +/- no/ has extra perm ?  补码 XXXX deprecated !!!!!!!!
	// guest can't be owner !!!!!!!!!!!!!!
	// np.binary_repr(np.bitwise_and(np.int8(-0b1010110), np.int(0b1111111)))  wrong !!!
	// np.binary_repr(np.bitwise_and(np.abs(np.int8(-0b1010110)), np.int(0b1111111)))  right
	Permissions int64 `db:"pm" schema:"int"`
}

type DataTypeType uint

const (
	DataTypeJson DataTypeType = iota
	DataTypeText
	DataTypeHTML
)

var DateTypeNameMap = map[DataTypeType]string {
	DataTypeJson: "json",
	DataTypeText: "text",
	DataTypeHTML: "html",
}

type DataNode struct {
	Field string `db:"d|f"`
	Data interface{} `db:"-" client:"data"`
	RawData string `db:"rd" client:"-" schema:"string"`  // TODO: string?
	Type DataTypeType `db:"dt" client:"-" schema:"int"`
	Language string // FIXME 应该在这一层， 还是Data那一层?
	Version uint
}

func (dn *DataNode) UnmarshalFromDB(valueData *DecodeValueData) error {
	valueData.DecodeDefault()
	api := valueData.GetAPI()
	if dn.Type == DataTypeJson {
		dn.Data = api.Get([]byte(dn.RawData))
	} else {
		dn.Data = &dn.RawData
	}
	return nil
}

type DataSet struct{
	dataMap map[string]int
	data []*DataNode
}

func (ds *DataSet) Get(f string) (*DataNode, bool) {
	if idx, ok := ds.dataMap[f]; ok {
		return ds.data[idx], true
	} else {
		return nil, false
	}
}

func (ds *DataSet) Set(f string, v *DataNode) (err error) {
	err = nil
	if idx, ok := ds.dataMap[f]; ok {
		ds.data[idx] = v
	} else {
		ds.data = append(ds.data, v)
		ds.dataMap[f] = len(ds.data) - 1
	}
	return
}

func (ds *DataSet) UnmarshalFromDB(valueData *DecodeValueData) error {
	iter := valueData.GetIterator()
	//ds.data = make([]*DataNode, 0)
	iter.ReadVal(&ds.data)
	ds.dataMap = make(map[string]int)
	for i, d := range ds.data {
		ds.dataMap[d.Field] = i
	}
	return nil
}

type Document struct {
	assetInfoNode
	Data *DataSet `db:"d" client:"data" schema:"uid"`
}

type User struct {
	nodeBase
}

func UnmarshalNode(b []byte) (node Node, err error) {
	if nt := DGraphJsonMarshaller.Get(b, "t"); nt.ValueType() == jsoniter.NumberValue {
		nodeType := NodeType(nt.ToUint())
		switch nodeType {
		case TypeDocument:
			node = &Document{}
			err = DGraphJsonMarshaller.Unmarshal(b, node)
		default:
			node = nil
			err = NewError(ErrorInvalidGraphNodeType, fmt.Sprintf("No Such Node Type code: %d.", nodeType))
		}
	} else {
		node = nil
		err = NewError(ErrorInvalidGraphNodeType, fmt.Sprintf("Can't find node type field."))
	}
	return
}

func MarshalNode(n Node) (b []byte, err error) {
	b, err = DGraphJsonMarshaller.Marshal(n)
	return
}



