package list

import (
	"testing"

	"github.com/okke/elmo/core"
)

func listContext() elmo.RunContext {
	context := elmo.NewGlobalContext()
	context.RegisterModule(Module)
	return context
}

func TestNew(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    l: (list.new 1 2 3)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 2 3]")))
}

func TestLen(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    list.len [1 2 3]`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(3)))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    list.len "chipotle"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    list.len`, elmo.ExpectErrorValueAt(t, 2))
}

func TestAt(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		list.at`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		list.at 1 2`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    list.at [1 2 3] 0`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(1)))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		list.at [1 2 3] 0 0`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		list.at [1 2 3] 0 2`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 2 3]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		list.at [1 2 3] -1`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(3)))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		list.at [1 2 3] -1 -2`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[3 2]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		list.at [1 2 3] -2 -1`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[2 3]")))

}

func TestAppend(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    l: [1 2 3]
    l: (list.append (l) 4)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 2 3 4]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    l1: [1 2 3]
    l2: (list.append (l1) 4)
		l1`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 2 3]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    l: [1 2 3]
    l: (list.append (l) 4 5 6)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 2 3 4 5 6]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    l: [1 2 3]
    l: (list.append l 4 )
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 2 3 4]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
	   l: [1 2 3]
	   list.append! $l 4
		 l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 2 3 4]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
	 	 l: [1 2 3]
		 freeze! $l
	 	 list.append! $l 4`, elmo.ExpectErrorValueAt(t, 4))

}

func TestPrepend(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    l: [1 2 3]
    l: (list.prepend l 4)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[4 1 2 3]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    l: [1 2 3]
    l: (list.prepend l 4 5 6)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[6 5 4 1 2 3]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		l: [1 2 3]
		list.prepend! l 4 5 6
	  l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[6 5 4 1 2 3]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
	 	 l: [1 2 3]
		 freeze! $l
	 	 list.prepend! $l 4`, elmo.ExpectErrorValueAt(t, 4))
}

func TestEach(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		 [1 2 3] |list.each v soep`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    l: [1 2 3]
    list.each l v {
		  once result []
			result: (list.append result (v))
	  }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 2 3]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
	    l: [a b c]
	    list.each l v i {
			  once result []
				result: (list.append result (i) (v))
		  }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[0 a 1 b 2 c]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
			l: [1 2 3]
			touch: 99
			list.each l v {
				touch: (v)
			}
			touch`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(99)))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
			l: [1 2 3]
			stepsize: 99
			list.each l v {
				incr index (stepsize)
			}`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(297)))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		 f: (func x {return $x})
	 	 [1 2 3] |list.each &f`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(3)))

}

func TestMap(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
    l: [a b c]
    list.map l v {
		  true
	  }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[(true) (true) (true)]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		l: [a b c]
		list.map l v {
			incr index
		}`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 2 3]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		 l: [a b c]
		 list.map l v {
			 v
		 }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[a b c]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
 		 l: [a b c]
		 list.map l v {
		   nil
		 }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[(nil) (nil) (nil)]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
	 	  l: [a b c]
	 	  list.map l v i {
	 			[(i) (v)]
	 		}`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[[0 a] [1 b] [2 c]]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load list)
		 f: (func i {return (multiply $i $i)})
		 [1 2 3] |list.map &f`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 4 9]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
			l: [a b c]
			list.map l v i {
				if (eq (i) 1) {
					[(i) (v)]
				} else {
					return 99
				}
			}`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[99 [1 b] 99]")))

}

func TestFilter(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
     list.filter [a b c] v {
		   true
	   }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[a b c]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
	   list.filter [1 2 3] v {
	 		  eq (v) 2
	 	 }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[2]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		 list.filter [1 2 3] v {
				ne (v) 2
		 }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 3]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
	 	 list.filter [2 1 0] v i {
	 	   ne $v $i
	 	 }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[2 0]")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`list: (load "list")
		 f: (func v {return (ne (v) 2)})
	 	 list.filter [1 2 3] &f`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[1 3]")))
}
