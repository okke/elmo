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
		`((file "file_testdata/peppers.txt") exists)`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`((file "file_testdata") exists)`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`((file "file_testdata") isDir)`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`((file "file_testdata/nofile") exists)`, ExpectValue(t, False))

	ParseTestAndRunBlock(t,
		`f:  (file "file_testdata/nofile")
		 defined f.exists`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`f:  (file "file_testdata/nofile")
		 defined f.string`, ExpectValue(t, True))

	ParseTestAndRunBlock(t,
		`f:  (file "file_testdata/nofile")
		 defined f.absPath`, ExpectValue(t, False))

	ParseTestAndRunBlock(t,
		`((file ".") string)`, ExpectErrorValueAt(t, 1))

}

func TestTempFile(t *testing.T) {
	ParseTestAndRunBlock(t,
		`f: (file (tempFile tmp { return $tmp.absPath }))
		not $f.exists |assert`, ExpectValue(t, True))
}

func TestFileWrite(t *testing.T) {
	ParseTestAndRunBlock(t,
		`tempFile tmp { 
			f: (tmp.write "chipotle")
			return $f.size
		}`, ExpectValue(t, NewIntegerLiteral(8)))

	ParseTestAndRunBlock(t,
		`tempFile tmp { 
			f: (tmp.write "chipotle")
			f: (tmp.write "jalapeno")
			return $f.size
		}`, ExpectValue(t, NewIntegerLiteral(8)))
}

func TestFileAppend(t *testing.T) {
	ParseTestAndRunBlock(t,
		`tempFile tmp { 
			f: (tmp.append "chipotle")
			return $f.size
		}`, ExpectValue(t, NewIntegerLiteral(8)))

	ParseTestAndRunBlock(t,
		`tempFile tmp { 
			f: (tmp.append "chipotle")
			f: (tmp.append "jalapeno")
			return $f.size
		}`, ExpectValue(t, NewIntegerLiteral(16)))
}
