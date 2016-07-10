package elmo

import "testing"

func withScript(s string, testFunc func(ScriptMetaData)) {
	testFunc(NewScriptMetaData("test", s))
}

func expectLineNoAndColumn(t *testing.T, pos int, lineno int, column int) func(ScriptMetaData) {
	return func(meta ScriptMetaData) {
		l, c := meta.PositionOf(pos)
		if l != lineno {
			t.Errorf("expected lineno (%d), got (%d)", lineno, l)
		}

		if c != column {
			t.Errorf("expected column (%d), got (%d)", column, c)
		}
	}
}

func TestGetPositionOf(t *testing.T) {
	// postion 0 should always be on line 1 and column 1
	//
	withScript("soup\nsoup", expectLineNoAndColumn(t, 0, 1, 1))

	// last char on first line should be on line 1, column 4
	//
	withScript("soup\nsoup", expectLineNoAndColumn(t, 3, 1, 4))

	// newline of first line should be on the next line (2) at column 0
	withScript("soup\nsoup", expectLineNoAndColumn(t, 4, 2, 0))

	// first character of last line should be on line 3, column 1
	//
	withScript("soup\nsoup\nsoup", expectLineNoAndColumn(t, 10, 3, 1))
}
