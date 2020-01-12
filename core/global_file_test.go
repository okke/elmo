package elmo

import "testing"

func TestFile(t *testing.T) {

	ParseTestAndRunBlock(t,
		`f: (file "file_testdata/peppers.txt"); f.name`, ExpectValue(t, NewStringLiteral("peppers.txt")))

	ParseTestAndRunBlock(t,
		`((file "file_testdata/peppers.txt") string)`, ExpectValue(t, NewStringLiteral("chipotle,jalapeno,habanero")))

	ParseTestAndRunBlock(t,
		`((file "file_testdata/peppers.txt") binary) | type`, ExpectValue(t, NewIdentifier("binary")))

	ParseTestAndRunBlock(t,
		`((file "file_testdata/peppers.txt") exists)`, ExpectValue(t, NewBooleanLiteral(true)))

	ParseTestAndRunBlock(t,
		`((file "file_testdata") exists)`, ExpectValue(t, NewBooleanLiteral(true)))

	ParseTestAndRunBlock(t,
		`((file "file_testdata") isDir)`, ExpectValue(t, NewBooleanLiteral(true)))

	ParseTestAndRunBlock(t,
		`((file "file_testdata/nofile") exists)`, ExpectValue(t, NewBooleanLiteral(false)))

	ParseTestAndRunBlock(t,
		`((file ".") string)`, ExpectErrorValueAt(t, 1))

}
