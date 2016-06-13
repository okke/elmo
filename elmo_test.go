package elmo

import "testing"

func expectOneLine(t *testing.T) func(*node32, string) {
	return func(ast *node32, buf string) {
		if !TestEqRules(ChildrenRules(ast), []pegRule{ruleLine}) {
			t.Error("does not contain one line")
		}
	}
}

func expectOneLineWith(t *testing.T, testChildren func([]*node32)) func(*node32, string) {
	return func(ast *node32, buf string) {
		if !TestEqRules(ChildrenRules(ast), []pegRule{ruleLine}) {
			t.Error("does not contain one line")
		}

		testChildren(Children(ast.up))
	}
}

func expectTwoLines(t *testing.T) func(*node32, string) {
	return func(ast *node32, buf string) {
		if !TestEqRules(ChildrenRules(ast), []pegRule{ruleLine, ruleLine}) {
			t.Error("does not contain two lines")
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

func TestParseCommandWithOneParameter(t *testing.T) {
	ParseAndTest(t, "chipotle sauce", expectOneLineWith(t, func(children []*node32) {
		// TODO!!!

	}))
}
