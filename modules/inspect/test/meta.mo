inspect: (load inspect)

toBeInspected: (func pepper sku { puts $pepper $sku })

suite: {

    testMetaOfBlock: (func {
        meta: (inspect.meta { })

        eq (meta.beginsAt) 148 | assert
        eq (meta.length) 3 | assert
        eq (meta.fileName) "./test/meta.mo" | assert
        eq (meta.code) "{ }" | assert
    })

    testMetaOfFunc: (func {
        meta: (inspect.meta &toBeInspected)

        # normal meta data of a function should be the meta
        # data of its coresponding code block
        #
        eq (meta.beginsAt) 57 | assert
        eq (meta.length) 21 | assert
        eq (meta.fileName) "./test/meta.mo" | assert
        eq (meta.code) "{ puts $pepper $sku }" | assert
        
        # functions should have arguments
        #
        eq (meta.arguments) [pepper sku] | assert
    })
}

test suite