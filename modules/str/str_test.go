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

func TestReplaceAll(t *testing.T) {
	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
 		 str.replaceAll "chipotle in a jar" "in" "out"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle out a jar")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
 		 str.replaceAll "chipotle in a jar of chipotles" "chipotle" "jalapeno"`, elmo.ExpectValue(t, elmo.NewStringLiteral("jalapeno in a jar of jalapenos")))
}

func TestReplaceFirst(t *testing.T) {
	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
 		 str.replaceFirst "chipotle in a jar in a box" "in" "out"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle out a jar in a box")))
}

func TestReplaceLast(t *testing.T) {
	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
 		 str.replaceLast "chipotle in a jar in a box" "in" "out"`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle in a jar out a box")))
}

func TestFindFirst(t *testing.T) {
	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
 		 str.findFirst "chipotle in a jar" "in"`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(9)))
}

func TestFindLast(t *testing.T) {
	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
         str.findLast "chipotle in a bin" "in"`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(15)))
}

func TestFindAll(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     	 str.findAll "galalalapeno" "la"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[2 4 6]")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     	 str.findAll "gagagagalapeno" "gaga"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[0 4]")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     	 str.findAll "galapeno" "galapeno"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[0]")))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
     	 str.findAll "galapeno" "peno"`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[4]")))
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

func TestPadLeft(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padLeft "soup" "?"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padLeft "soup" 8 |len |eq 8 |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		str.padLeft "soup" 0 |len |eq 0 |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		str.padLeft "soup" 8 |eq "    soup" |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padLeft "soup" 8 "+" |eq "++++soup" |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padLeft "soup" 8 "+-" |eq "+-+-soup" |assert`, elmo.ExpectValue(t, elmo.True))
}

func TestPadRight(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padRight "soup" "?"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padRight "soup" 8 |len |eq 8 |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		str.padRight "soup" 0 |len |eq 0 |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		str.padRight "soup" 8 |eq "soup    " |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padRight "soup" 8 "+" |eq "soup++++" |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padRight "soup" 8 "+-" |eq "soup+-+-" |assert`, elmo.ExpectValue(t, elmo.True))
}

func TestPadBoth(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padBoth "soup" "?"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padBoth "soup" 8 |len |eq 8 |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		str.padBoth "soup" 0 |len |eq 0 |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		str.padBoth "soup" 8 |eq "  soup  " |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padBoth "soup" 8 "+" |eq "++soup++" |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		 str.padBoth "soup" 9 "+" |eq "++soup+++" |assert`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, strContext(),
		`str: (load "string")
		  str.padBoth "sou" 8 "+" |eq "+++sou++" |assert`, elmo.ExpectValue(t, elmo.True))

}
