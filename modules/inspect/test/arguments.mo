inspect: (load inspect)
list: (load list)
str: (load string)

suite: {

    testCallsHaveArgumentsWhichAreInspectable: (func {
        calls: (inspect.calls {
            plus 5 3.0 
        })
        args: (inspect.arguments (calls 0))
        first: (inspect.meta (args 0))
        second: (inspect.meta (args 1))
        eq (first value) 5
        eq (first type) int
        eq (second value) 3.0
        eq (second type) float
    })

}

test suite