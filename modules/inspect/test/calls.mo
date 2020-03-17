inspect: (load inspect)
str: (load string)

suite: {

    testEmptyBlockHasNoCalls: (func {
        inspect.calls {} | eq [] | assert
    })

    testBlockWithCallsHasCalls: (func {
        calls: (inspect.calls {
            puts chipotle
            puts jalapeno
        })
        len $calls |eq 2 |assert
    })

    testPipedCallIsSingleCall: (func {
        calls: (inspect.calls {
            plus 5 3 | plus 8 |eq 16 |assert
        })
        
        len $calls |eq 1 |assert
    })

    testCallsHaveMetaData: (func {
        calls: (inspect.calls {
            plus 5 3 | plus 8 |eq 16 |assert
        })
        eq $((inspect.meta (calls 0)) code) "plus 5 3 | plus 8 |eq 16 |assert" |assert
    })

    testCallsHaveNameInMetaData: (func {
        calls: (inspect.calls {
            plus 5 3 
        })
        first: (inspect.meta (calls 0))
        eq (first name) "plus" |assert
    })

    testPipedCallIsInMeta: (func {
        calls: (inspect.calls {
            plus 5 3 | plus 8
        })
        first: (inspect.meta (calls 0))
        second: (inspect.meta (first pipe))
        eq $(second code) "plus 8" |assert
    })    
}

test suite