http: (load http)
str: (load string)

suite: {

    testTestServerAcceptsCallbackFunction: (func {
        server: (http.testServer (func request response {

        }))
        close $server
    })

    testTestServerHasTestURL: (func {
        server: (http.testServer (func request response {}))
        url: (http.testURL $server)
        str.startsWith $url "http://" |assert  
        close $server
    })

    testTestServerHasNoURLAfterClose: (func {
        server: (http.testServer (func erquest response {}))
        close $server
        url: (http.testURL $server)
        eq $url "" |assert
    })

}

test suite