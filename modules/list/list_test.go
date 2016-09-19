package el

import (
	"testing"

	"github.com/okke/elmo/core"
)

func listContext() elmo.RunContext {
	context := elmo.NewGlobalContext()
	context.RegisterModule(Module)
	return context
}

func TestAppend(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l: [1 2 3]
    l: (el.append (l) 4)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 1 2 3 4")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l1: [1 2 3]
    l2: (el.append (l1) 4)
		l1`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 1 2 3")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l: [1 2 3]
    l: (el.append (l) 4 5 6)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 1 2 3 4 5 6")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l: [1 2 3]
    l: (el.append l 4 )
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 1 2 3 4")))
}

func TestPrepend(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l: [1 2 3]
    l: (el.prepend l 4)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 4 1 2 3")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l: [1 2 3]
    l: (el.prepend l 4 5 6)
		l`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 6 5 4 1 2 3")))
}

func TestEach(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l: [1 2 3]
    el.each l v {
		  once result (list)
			result: (el.append result (v))
	  }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 1 2 3")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
	    l: [a b c]
	    el.each l v i {
			  once result (list)
				result: (el.append result (i) (v))
		  }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 0 a 1 b 2 c")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
			l: [1 2 3]
			touch: 99
			el.each l v {
				touch: (v)
			}
			touch`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(99)))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
			l: [1 2 3]
			stepsize: 99
			el.each l v {
				incr index (stepsize)
			}`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(297)))

}

func TestMap(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
    l: [a b c]
    el.map l v {
		  true
	  }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list (true) (true) (true)")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
		l: [a b c]
		el.map l v {
			incr index
		}`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 1 2 3")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
		 l: [a b c]
		 el.map l v {
			 v
		 }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list a b c")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
 		 l: [a b c]
		 el.map l v {
		   nil
		 }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list (nil) (nil) (nil)")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
	 	  l: [a b c]
	 	  el.map l v i {
	 			list (i) (v)
	 		}`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list (list 0 a) (list 1 b) (list 2 c)")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
			l: [a b c]
			el.map l v i {
				if (eq (i) 1) {
					list (i) (v)
				} else {
					return 99
				}
			}`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 99 (list 1 b) 99")))

}

func TestFilter(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
     el.filter [a b c] v {
		   true
	   }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list a b c")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
	   el.filter [1 2 3] v {
	 		  eq (v) 2
	 	 }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 2")))

	elmo.ParseTestAndRunBlockWithinContext(t, listContext(),
		`el: (load "el")
		 el.filter [1 2 3] v {
				ne (v) 2
		 }`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "list 1 3")))
}
