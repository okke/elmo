package elmo

// filterNodes filters an array of nodes
//
func filterNodes(nodes []*node32, f func(*node32) bool) []*node32 {
	result := []*node32{}
	for _, v := range nodes {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// nodeChildren returns an array of children of a node in the parsetree without
// spacing, newlines other non-semantic tokens
//
func nodeChildren(node *node32) []*node32 {
	result := []*node32{}

	cursor := node.up
	for cursor != nil {
		result = append(result, cursor)
		cursor = cursor.next
	}

	return filterNodes(result, func(child *node32) bool {
		return (child.pegRule != ruleSpacing &&
			child.pegRule != ruleNewLine &&
			child.pegRule != ruleEndOfLine &&
			child.pegRule != ruleLCURLY &&
			child.pegRule != ruleRCURLY &&
			child.pegRule != ruleLBRACKET &&
			child.pegRule != ruleRBRACKET &&
			child.pegRule != ruleLPAR &&
			child.pegRule != ruleRPAR &&
			child.pegRule != ruleDOLLAR &&
			child.pegRule != ruleCOMMA)
	})
}

// pegRules returns an array of the peg rules of a node without Spacing
//
func pegRules(nodes []*node32) []pegRule {
	result := []pegRule{}

	for _, v := range nodes {
		result = append(result, v.pegRule)
	}

	return result
}

// pegRulesFirstChild returns an array of the peg rules of a nodes first child
//
func pegRulesFirstChild(nodes []*node32) []pegRule {
	result := []pegRule{}

	for _, v := range nodes {
		result = append(result, v.up.pegRule)
	}

	return result
}

// TestEqRules will test if two array of rules are the same
//
func ruleSlicesAreEqual(a, b []pegRule) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// ChildrenRules returns an array of the children rules
//
func childrenRules(node *node32) []pegRule {
	return pegRules(nodeChildren(node))
}

// Text returns the textual representation of a node without any Spacing
//
func nodeText(node *node32, buf string) string {

	begin := node.begin
	end := node.end
	cursor := node.up
	if cursor != nil {
		begin = cursor.begin

		for cursor != nil {
			if cursor.pegRule == ruleSpacing {
				break
			}
			end = cursor.end
			cursor = cursor.next
		}
	}

	return buf[begin : begin+(end-begin)]

}
