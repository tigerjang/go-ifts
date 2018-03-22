package dgraph_ldap

//import "time"
//
//func (n *nodeBase) InitBeforeUnmarshal() error {
//	n.Children = make([]Node, 0)
//	n.CreateTime = time.Now()
//	n.ModifyTime = time.Now()
//	return nil
//}
//
func (n *nodeBase) GetChildren() error {
	return nil  // TODO
}

//
//
//func (n *assetInfoNode) InitBeforeUnmarshal() error {
//	// TODO
//	n.nodeBase.InitBeforeUnmarshal()
//	return nil
//}

//func (n *Document) InitBeforeUnmarshal() error {
//	n.assetInfoNode.InitBeforeUnmarshal()
//	n.Data = &DataSet{}
//	return nil
//}
