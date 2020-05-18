
loaded: (load loadedtemplates)

suite: {

    testScopeOfLoadedTemplates: (func {
        (loaded.t1 "42") |eq "4242" |assert
    })

}

test suite