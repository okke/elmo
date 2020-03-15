inspect: (load inspect)
list: (load list)
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
        list.len $calls |eq 2 |assert
    })

    testPipedCallIsSingleCall: (func {
        calls: (inspect.calls {
            plus 5 3 | plus 8 |eq 16 |assert
        })
        
        list.len $calls |eq 1 |assert
    })

    testCallsHaveMetaData: (func {
        calls: (inspect.calls {
            plus 5 3 | plus 8 |eq 16 |assert
        })
        eq $((inspect.meta (calls 0)) code) "plus 5 3 | plus 8 |eq 16 |assert" |assert
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