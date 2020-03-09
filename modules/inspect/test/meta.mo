inspect: (load inspect)

suite: {

    testMetaOfBlock: (func {
        meta: (inspect.meta { })

        eq (meta.beginsAt) 92 | assert
        eq (meta.length) 3 | assert
        eq (meta.fileName) "./test/meta.mo" | assert
        eq (meta.code) "{ }" | assert
    })
}

test suite