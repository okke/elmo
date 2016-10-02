package elmo

import "strings"

// Filter nodes
//
func Filter(nodes []*node32, f func(*node32) bool) []*node32 {
	result := []*node32{}
	for _, v := range nodes {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// Children returns an array of children of a node in the parsetree without
// spacing, newlines other non-semantic tokens
//
func Children(node *node32) []*node32 {
	result := []*node32{}

	cursor := node.up
	for cursor != nil {
		result = append(result, cursor)
		cursor = cursor.next
	}

	return Filter(result, func(child *node32) bool {
		return (child.pegRule != ruleSpacing &&
			child.pegRule != ruleNewLine &&
			child.pegRule != ruleEndOfLine &&
			child.pegRule != ruleLCURLY &&
			child.pegRule != ruleRCURLY &&
			child.pegRule != ruleLBRACKET &&
			child.pegRule != ruleRBRACKET &&
			child.pegRule != ruleLPAR &&
			child.pegRule != ruleRPAR &&
			child.pegRule != ruleDOLLAR)
	})
}

// PegRules returns an array of the peg rules of a node without Spacing
//
func PegRules(nodes []*node32) []pegRule {
	result := []pegRule{}

	for _, v := range nodes {
		result = append(result, v.pegRule)
	}

	return result
}

// PegRulesFirstChild returns an array of the peg rules of a nodes first child
//
func PegRulesFirstChild(nodes []*node32) []pegRule {
	result := []pegRule{}

	for _, v := range nodes {
		result = append(result, v.up.pegRule)
	}

	return result
}

// TestEqRules will test if two array of rules are the same
//
func TestEqRules(a, b []pegRule) bool {

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
func ChildrenRules(node *node32) []pegRule {
	return PegRules(Children(node))
}

// Text returns the textual representation of a node without any Spacing
//
func Text(node *node32, buf string) string {
	return strings.TrimSpace(buf[node.begin : node.begin+(node.end-node.begin)])
}
