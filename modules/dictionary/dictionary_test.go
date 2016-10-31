package dict

import (
	"testing"

	"github.com/okke/elmo/core"
)

func dictContext() elmo.RunContext {
	context := elmo.NewGlobalContext()
	context.RegisterModule(Module)
	return context
}

func TestDictionary(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.dict 1 2 3 4)
		 h 1`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(2)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.dict [1 2 3 4])
 		 h 1`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(2)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
	   h: (d.dict 1 2 3 4)
 		 h 3`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(4)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.dict [1 2 3 4])
  	 h 3`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(4)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.dict 1 2 3 4)
  	 h 5`, elmo.ExpectNothing(t))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.dict "1" 2 3 4)
	   h 1`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(2)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.dict a 2 b 4)
	 	 h b`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(4)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 l: [a 2 b 4]
		 h: (d.dict (l))
 	 	 h b`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(4)))
}

func TestDictionaryWithBlock(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.dict {
			set b 4
		 })
 	 	 h b`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(4)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 val: "chipotle"
		 h: (d.dict {
 			set b (val)
 		 })
  	 h b`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 val: "chipotle"
 		 h: (d.dict {
  		set val "galapeno"
			set b (val)
  	 })
   	 h b`, elmo.ExpectValue(t, elmo.NewStringLiteral("galapeno")))
}

func TestNewConstructsDictionary(t *testing.T) {
	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 peppers: {
		  hot: (func {
		  	return "chipotle"
		  })
		}
		sauce: (d.new (peppers))
		sauce.hot`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 peppers: {
		  hot: (func {
		  	return "chipotle"
		  })
		 }
		 sauce: (d.new (peppers) {
		   same: (func {
		 	  return (this.hot)
		   })
	   })
		 sauce.same`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 peppers: {
		  hot: (func {
		  	return "chipotle"
		  })
		 }
		 sauce: (d.new (peppers) {
		  hot: (func {
			  return "galapeno"
		  })
	   })
		 sauce.hot`, elmo.ExpectValue(t, elmo.NewStringLiteral("galapeno")))
}

func TestKeys(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		d.keys`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		d.keys "chipotle"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
    peppers: {
      chipotle: {
        heat: 2
      }
      galapeno: {
        heat: 3
      }
    }
		d.keys peppers`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[chipotle galapeno]")))

	// keys function should return keys in sorted order
	//
	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
    peppers: {
      galapeno: {
        heat: 3
      }
      chipotle: {
        heat: 2
      }
    }
    d.keys peppers`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[chipotle galapeno]")))

}

func TestKnows(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 d.knows`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 d.knows "chipotle" chipotle`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
    peppers: { chipotle: { heat: 2 } }
		d.knows peppers chipotle`, elmo.ExpectValue(t, elmo.NewBooleanLiteral(true)))
}
