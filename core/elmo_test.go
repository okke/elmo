package elmo

import "testing"

func expectOneLine(t *testing.T) func(*node32) {
	return func(ast *node32) {
		if !ruleSlicesAreEqual(childrenRules(ast), []pegRule{ruleLine}) {
			t.Errorf("does not contain one line, found %v", childrenRules(ast))
		}
	}
}

func expectOneLineContaining(t *testing.T, testChildren func([]*node32)) func(*node32) {
	return func(ast *node32) {

		if !ruleSlicesAreEqual(childrenRules(ast), []pegRule{ruleLine}) {
			t.Error("does not contain one line")
		}

		testChildren(nodeChildren(ast.up))
	}
}

func expectTwoLines(t *testing.T) func(*node32) {
	return func(ast *node32) {
		if !ruleSlicesAreEqual(childrenRules(ast), []pegRule{ruleLine, ruleLine}) {
			t.Error("does not contain two lines")
		}
	}
}

func IdentifierFollowedByShortcutAndArgument(t *testing.T, cut pegRule, ruleType pegRule) func([]*node32) {
	return func(children []*node32) {
		if !ruleSlicesAreEqual(pegRules(children), []pegRule{ruleArgument, cut, ruleArgument}) {
			t.Errorf("expected <identifier> <argument>, found %v", children)
		}
	}
}

func IdentifierFollowedByOneArgument(t *testing.T, ruleType pegRule) func([]*node32) {
	return func(children []*node32) {
		if !ruleSlicesAreEqual(pegRules(children), []pegRule{ruleArgument, ruleArgument}) {
			t.Errorf("expected <identifier> <argument>, found %v", children)
		}
		if children[1].up.pegRule != ruleType {
			t.Errorf("unexpected ruletype of argument, found %v", children[1].up)
		}
	}
}

func IdentifierFollowedByMultipleArguments(t *testing.T, ruleTypes []pegRule) func([]*node32) {
	return func(children []*node32) {

		if children[0].pegRule != ruleArgument {
			t.Errorf("expected to start with an identifier, found %v", children[0])
		}

		if !ruleSlicesAreEqual(pegRulesFirstChild(children[1:]), ruleTypes) {
			t.Errorf("expected %v, found %v", ruleTypes, pegRulesFirstChild(children[1:]))
		}
	}
}

func IdentifierFollowedByBlock(t *testing.T, blockTestFunc func(*node32)) func([]*node32) {
	return func(children []*node32) {

		if !ruleSlicesAreEqual(pegRules(children), []pegRule{ruleArgument, ruleArgument}) {
			t.Errorf("expected <identifier> <block>, found %v", children)
		}

		blockTestFunc(children[1].up)
	}
}

func IdentifierFollowedbyPipe(t *testing.T) func([]*node32) {
	return func(children []*node32) {
		if !ruleSlicesAreEqual(pegRules(children), []pegRule{ruleArgument, rulePipedOutput}) {
			t.Errorf("expected <identifier> <pipe>, found %v", children)
		}
		if !ruleSlicesAreEqual(pegRules(nodeChildren(children[1])), []pegRule{rulePIPE, ruleLine}) {
			t.Errorf("expected <identifier> <pipe>, found %v", children[1])
		}
	}
}

func TestParseSimpleCommand(t *testing.T) {
	ParseAndTest(t, "chipotle", expectOneLine(t))
}

func TestParseSimpleCommandWithWhiteSpace(t *testing.T) {
	ParseAndTest(t, " chipotle ", expectOneLine(t))
}

func TestParseSimpleCommandWithNewLines(t *testing.T) {
	ParseAndTest(t, "\nchipotle ", expectOneLine(t))
}

func TestParseTwoSimpleCommands(t *testing.T) {
	ParseAndTest(t, "chipotle;chipotle", expectTwoLines(t))
}

func TestParseTwoSimpleCommandsWithWhiteSpace(t *testing.T) {
	ParseAndTest(t, " chipotle; chipotle ", expectTwoLines(t))
}

func TestParseTwoSimpleCommandsOnNewLines(t *testing.T) {
	ParseAndTest(t, "chipotle\nchipotle", expectTwoLines(t))
}

func TestParseTwoSimpleCommandsOnNewLinesWithWhiteSpace(t *testing.T) {
	ParseAndTest(t, " chipotle\n chipotle ", expectTwoLines(t))
}

func TestParseTwoSimpleCommandsOnNewLinesWithMoreNewLines(t *testing.T) {
	ParseAndTest(t, "chipotle\n\nchipotle", expectTwoLines(t))
}

func TestParseTwoSimpleCommandsOnNewLinesWithMoreNewLinesAndSpacing(t *testing.T) {
	ParseAndTest(t, "chipotle\n \n\n chipotle", expectTwoLines(t))
}

func TestParseCommandWithIdentifierAsParameter(t *testing.T) {
	ParseAndTest(t, "chipotle sauce", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleIdentifier)))
}

func TestParseCommandWithStringAsParameter(t *testing.T) {
	ParseAndTest(t, "chipotle \"sauce\"", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleStringLiteral)))
}

func TestParseCommandWithIntegerAsParameter(t *testing.T) {
	ParseAndTest(t, "chipotle 138", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleNumber)))
	ParseAndTest(t, "chipotle 0", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleNumber)))
	ParseAndTest(t, "chipotle -6", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleNumber)))
}

func TestParseCommandWithFunctionCallAsParameter(t *testing.T) {
	ParseAndTest(t, "chipotle (sauce)", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleFunctionCall)))
	ParseAndTest(t, "chipotle $sauce", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleFunctionCall)))
}

func TestParseCommandWithEmptyBlockAsParameter(t *testing.T) {
	ParseAndTest(t, "chipotle {}", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleBlock)))
}

func TestParseCommandWithListAsParameter(t *testing.T) {
	ParseAndTest(t, "chipotle [1 2 3]", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleList)))
	ParseAndTest(t, `chipotle [1
		2 3]`, expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleList)))
	ParseAndTest(t, `chipotle [
		1
		2
		3
	]`, expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleList)))
}

func TestParseCommandWithMultipleParameters(t *testing.T) {
	ParseAndTest(t, "chipotle sauce in_a_jar", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleIdentifier, ruleIdentifier})))

	ParseAndTest(t, "chipotle sauce 128", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleIdentifier, ruleNumber})))

	ParseAndTest(t, "chipotle (sauce 128) (jar 136)", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleFunctionCall, ruleFunctionCall})))

	ParseAndTest(t, "chipotle \"sauce\" {}", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleStringLiteral, ruleBlock})))

	ParseAndTest(t, "chipotle {} jar {}", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleBlock, ruleIdentifier, ruleBlock})))

	ParseAndTest(t, "jar \"chipotle\" \"chipotle\"", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleStringLiteral, ruleStringLiteral})))
}

func TestParseCommandWithMultipleCommaSeparatedParameters(t *testing.T) {
	ParseAndTest(t, "chipotle sauce, in_a_jar", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleIdentifier, ruleIdentifier})))

	ParseAndTest(t, "chipotle sauce, in_a_jar;", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleIdentifier, ruleIdentifier})))

	ParseAndTest(t, "chipotle sauce, 128, 132", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleIdentifier, ruleNumber, ruleNumber})))

	ParseAndTest(t, "chipotle sauce,\n in_a_jar", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleIdentifier, ruleIdentifier})))

	ParseAndTest(t, "chipotle sauce,\n \n in_a_jar", expectOneLineContaining(t, IdentifierFollowedByMultipleArguments(t,
		[]pegRule{ruleIdentifier, ruleIdentifier})))
}

func TestParseCommandWithBlockContainingNewlinesAsArguments(t *testing.T) {
	ParseAndTest(t, "chipotle {\n}", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleBlock)))
	ParseAndTest(t, "chipotle {\n\n}", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleBlock)))
	ParseAndTest(t, "chipotle {\n \n}", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleBlock)))

}

func TestParseCommandWithLeadingNewlines(t *testing.T) {
	ParseAndTest(t, "\nchipotle {\n \n}", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleBlock)))
	ParseAndTest(t, "\n\nchipotle {\n \n}", expectOneLineContaining(t, IdentifierFollowedByOneArgument(t, ruleBlock)))
}

func TestParseCommandWithBlockWithCalls(t *testing.T) {
	ParseAndTest(t, "\nchipotle {sauce 33}", expectOneLineContaining(t, IdentifierFollowedByBlock(t, expectOneLine(t))))
	ParseAndTest(t, "\nchipotle {sauce 33; sauce 34}", expectOneLineContaining(t, IdentifierFollowedByBlock(t, expectTwoLines(t))))
	ParseAndTest(t, "\nchipotle {sauce 33\n sauce 34}", expectOneLineContaining(t, IdentifierFollowedByBlock(t, expectTwoLines(t))))
	ParseAndTest(t, "\nchipotle {sauce {}}", expectOneLineContaining(t, IdentifierFollowedByBlock(t, expectOneLine(t))))
	ParseAndTest(t, "\nchipotle {sauce {\n}}", expectOneLineContaining(t, IdentifierFollowedByBlock(t, expectOneLine(t))))
}

func TestParseCommandWithShortcutAsParameter(t *testing.T) {
	ParseAndTest(t, "chipotle : sauce", expectOneLineContaining(t, IdentifierFollowedByShortcutAndArgument(t, ruleCOLON, ruleIdentifier)))
	ParseAndTest(t, "chipotle: sauce", expectOneLineContaining(t, IdentifierFollowedByShortcutAndArgument(t, ruleCOLON, ruleIdentifier)))
}

func TestParseCommandWithPipedOutput(t *testing.T) {
	ParseAndTest(t, "chipotle | sauce", expectOneLineContaining(t, IdentifierFollowedbyPipe(t)))
	ParseAndTest(t, "chipotle | sauce | jar", expectOneLineContaining(t, IdentifierFollowedbyPipe(t)))
	ParseAndTest(t, "chipotle | sauce 33", expectOneLineContaining(t, IdentifierFollowedbyPipe(t)))
	ParseAndTest(t, "chipotle | sauce 33 34 | jar 28", expectOneLineContaining(t, IdentifierFollowedbyPipe(t)))
}
