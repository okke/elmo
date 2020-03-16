package elmo

type astNode struct {
	meta    ScriptMetaData
	node    *node32
	endNode *node32
	//begin uint32
	//end   uint32
}

func (astNode *astNode) Meta() ScriptMetaData {
	return astNode.meta
}
func (astNode *astNode) BeginsAt() uint32 {
	if astNode == nil {
		return 0
	}
	return astNode.node.begin
}

func (astNode *astNode) EndsAt() uint32 {
	if astNode.node == nil {
		return 0
	}
	if astNode.endNode == nil {
		return endOfNode(astNode.node)
	}
	return endOfNode(astNode.endNode)
}
