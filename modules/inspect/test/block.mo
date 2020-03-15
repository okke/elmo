inspect: (load inspect)

suite: {

    testBlockOfFunction: (func {
        f: (func x {plus $x 3})
        meta: (inspect.block &f | inspect.meta)
        eq $(meta.code) "{plus $x 3}" | assert
    })    
}

test suite