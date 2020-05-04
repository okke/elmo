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
		 str.trim " chipotle "`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.trim "++chipotle--" "+-"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

}

func TestTrimLeft(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     	 str.trimLeft`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
  		 str.trimLeft " chipotle "`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle ")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
  		 str.trimLeft "+chipotle+" "+"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle+")))
}

func TestTrimRight(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     	 str.trimRight`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
  		 str.trimRight " chipotle "`, elmo.ExpectValue(t, elmo.NewStringLiteral(" chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
  		 str.trimRight "+chipotle+" "+"`, elmo.ExpectValue(t, elmo.NewStringLiteral("+chipotle")))
}

func TestTrimSuffix(t *testing.T) {
	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
         str.trimSuffix "chipotle_in_a_jar" "_in_a_jar"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))
}

func TestTrimPrefix(t *testing.T) {
	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
         str.trimPrefix "jar_with_chipotle" "jar_with_"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))
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

func TestSplit(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.split`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.split (error "chipotle") "chipotle"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.split "chipotle" (error "chipotle")`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.split "chipotle;jalapeno;chilli" ";"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), `["chipotle" "jalapeno" "chilli"]`)))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.split "chipotle"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), `["c" "h" "i" "p" "o" "t" "l" "e"]`)))
}

func TestStartsWith(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.startsWith`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.startsWith (error "chipotle") "chipotle"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.startsWith "chipotle" (error "chipotle")`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.startsWith "chipotle" "chi"`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.startsWith "chipotle" "cho"`, elmo.ExpectValue(t, elmo.False))

}

func TestEndsWith(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     str.endsWith`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.endsWith (error "chipotle") "chipotle"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.endsWith "chipotle" (error "chipotle")`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.endsWith "chipotle" "tle"`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
    str.endsWith "chipotle" "tla"`, elmo.ExpectValue(t, elmo.False))

}

func TestUpper(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
         str.upper "upper" |eq "UPPER" |assert`, elmo.ExpectValue(t, elmo.True))
}

func TestLower(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
         str.lower "LOWER" |eq "lower" |assert`, elmo.ExpectValue(t, elmo.True))
}
