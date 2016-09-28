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
}

func TestEach(t *testing.T) {

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
}
