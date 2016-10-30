package str

import (
	"testing"

	elmo "github.com/okke/elmo/core"
)

func strContext() elmo.RunContext {
	context := elmo.NewGlobalContext()
	context.RegisterModule(Module)
	return context
}

func TestAt(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.at`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.at "chipotle" 2`, elmo.ExpectValue(t, elmo.NewStringLiteral("i")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     pepper: "chipotle"
     str.at (pepper 0 99) 2`, elmo.ExpectErrorValueAt(t, 3))

}

func TestLen(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.len`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.len "chipotle"`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(8)))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     pepper: "chipotle"
		 str.len (pepper 0 2)`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(3)))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     pepper: "chipotle"
     str.len (pepper 0 99)`, elmo.ExpectErrorValueAt(t, 3))

}

func TestConcat(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.concat`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.concat "chipotle" " " "sauce"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle sauce")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     pepper: "chipotle"
     str.concat "sauce" (pepper 0 99)`, elmo.ExpectErrorValueAt(t, 3))

}

func TestTrim(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim (error "chipotle")`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim "chipotle" (error "chipotle")`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim " chipotle "`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim left " chipotle "`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle ")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
	   str.trim left "!chipotle!" "!"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle!")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim left (error "!chipotle!") "!"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim right "!chipotle!" "!"`, elmo.ExpectValue(t, elmo.NewStringLiteral("!chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim right " chipotle "`, elmo.ExpectValue(t, elmo.NewStringLiteral(" chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim prefix "jar_with_chipotle" "jar_with_"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim suffix "chipotle_in_a_jar" "_in_a_jar"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim soup " chipotle "`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.trim "<chipotle>" "<>"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))
}

func TestReplace(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace "chipotle in a jar" "in" "out"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle out a jar")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace (error "chipotle in a jar") "in" "out"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace "chipotle in a jar" (error "in") "out"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace "chipotle in a jar" "in" (error "out")`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
      str.replace "chipotle in a jar jar" "jar" "big"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle in a big jar")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace first "chipotle in a jar jar" "jar" "big"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle in a big jar")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace first (error "chipotle in a jar jar") "jar" "big"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace all "chipotle in a jar jar" "jar" "clueless"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle in a clueless clueless")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace last "chipotle in a jar jar binks" "jar" "with"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle in a jar with binks")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
      str.replace last "chipotle in a jar jar" " jar" "!"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle in a jar!")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace last "chipotle in a jar" "botle" "clueless"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle in a jar")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
      str.replace last "chipotle in a jar" "chipotle" "jalapeno"`, elmo.ExpectValue(t, elmo.NewStringLiteral("jalapeno in a jar")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.replace soup "chipotle in a jar" "in" "out"`, elmo.ExpectErrorValueAt(t, 2))

}

func TestFind(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find "chipotle in a jar" "in"`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(9)))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
      str.find (error "chipotle in a jar") "in" `, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find first (error "chipotle in a jar") "in" `, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find "chipotle in a jar" (error "in") `, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find soup "chipotle in a jar" "in" `, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find first "chipotle in a jar" "in"`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(9)))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find last "chipotle in a bin" "in"`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(15)))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find all "galalalapeno" "la"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[2 4 6]")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find all "gagagagalapeno" "gaga"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[0 4]")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find all "galapeno" "galapeno"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[0]")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.find all "galapeno" "peno"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[4]")))
}

func TestCount(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.count`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.count (error "chipotle") "chipotle"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.count "chipotle" (error "chipotle")`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.count "chipotle" "chipotle"`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(1)))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.count "jalapeno" "a"`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(2)))
}
