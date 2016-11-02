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

func TestNewWithoutParent(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 d.new`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
	 	 d.new "soep"`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.new [1 2 3 4])
 		 h 1`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(2)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.new [1 2 3 4])
  	 h 3`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(4)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.new [1 2 3 4])
  	 h 5`, elmo.ExpectNothing(t))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.new ["1" 2 3 4])
	   h 1`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(2)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.new [a 2 b 4])
	 	 h b`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(4)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 l: [a 2 b 4]
		 h: (d.new $l)
 	 	 h b`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(4)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 h: (d.new {
			set b 4
		 })
 	 	 h b`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(4)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 val: "chipotle"
		 h: (d.new {
 			set b $val
 		 })
  	 h b`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 val: "chipotle"
 		 h: (d.new {
  		set val "galapeno"
			set b $val
  	 })
   	 h b`, elmo.ExpectValue(t, elmo.NewStringLiteral("galapeno")))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 peppers: {
		  hot: (func {
		  	return "chipotle"
		  })
		}
		sauce: (d.new $peppers)
		sauce.hot`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))
}

func TestNewWithParent(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 sauce: (d.new "soup" {
		 })`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
	 	 peppers: {
	 		hot: (func {
	 			return "chipotle"
	 		})
	 	 }
	 	 sauce: (d.new $peppers "soup")`, elmo.ExpectErrorValueAt(t, 7))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 peppers: {
		  hot: (func {
		  	return "chipotle"
		  })
		 }
		 sauce: (d.new $peppers {
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

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
 		 peppers: {
 		  hot: (func {
 		  	return "chipotle"
 		  })
 		 }
 		 sauce: (d.new $peppers [same (func {return (this.hot)})])
 		 sauce.same`, elmo.ExpectValue(t, elmo.NewStringLiteral("chipotle")))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
 		 peppers: {
 		  hot: (func {
 		  	return (this.favourite)
 		  })
 		 }
		 more: {
		 	favourite: "jalapeno"
	   }
 		 sauce: (d.new (peppers) $more)
 		 sauce.hot`, elmo.ExpectValue(t, elmo.NewStringLiteral("jalapeno")))

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
		d.keys $peppers`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[chipotle galapeno]")))

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
    d.keys $peppers`, elmo.ExpectValue(t, elmo.ParseAndRun(elmo.NewGlobalContext(), "[chipotle galapeno]")))

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
		d.knows $peppers chipotle`, elmo.ExpectValue(t, elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		peppers: { chipotle: { heat: 2 } }
		d.knows $peppers jalapeno`, elmo.ExpectValue(t, elmo.False))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 peppers: {
		  hot: "jalapeno"
		 }
		 snacks: (d.new (peppers) {
		   lame: "chipotle"
	   })
		 d.knows $snacks hot`, elmo.ExpectValue(t, elmo.True))

}

func TestGet(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 d.get`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 d.get "chipotle" chipotle`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
    peppers: { chipotle: "hot" }
		d.get $peppers chipotle`, elmo.ExpectValues(t, elmo.NewStringLiteral("hot"), elmo.True))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		peppers: { chipotle: "lame" }
		d.get $peppers jalapeno`, elmo.ExpectValues(t, elmo.Nothing, elmo.False))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 peppers: {
		  hot: "jalapeno"
		 }
		 snacks: (d.new (peppers) {
		   lame: "chipotle"
	   })
		 d.get $snacks hot`, elmo.ExpectValues(t, elmo.NewStringLiteral("jalapeno"), elmo.True))

}

func TestMerge(t *testing.T) {

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 d.merge`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
 		 d.merge {}`, elmo.ExpectErrorValueAt(t, 2))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
  	 d.merge {} {}`, elmo.ExpectValue(t, elmo.NewDictionaryValue(nil, map[string]elmo.Value{})))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
		 m1: {}
		 m2: {}
	 	 d.merge $m1 $m2`, elmo.ExpectValue(t, elmo.NewDictionaryValue(nil, map[string]elmo.Value{})))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
 		 m1: {a:1}
 	 	 m2: (d.merge $m1 {b:2})
		 m2.b`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(2)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
  	 m: (d.merge (d.new [a 1 b 2]) (d.new [b 3]))
 		 m.b`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(3)))

	elmo.ParseTestAndRunBlockWithinContext(t, dictContext(),
		`d: (load "dict")
   	 m: (d.merge (d.new [a 1 b 2]) (d.new [b 3]) (d.new [a 2 b 4]))
  	 multiply $m.a $m.b`, elmo.ExpectValue(t, elmo.NewIntegerLiteral(8)))

}
