http: (load http)
dict: (load dict)
str: (load string)

server: (http.testServer (func request response {
    response.sendCookie "pepper" "jalapeno" 3600
    response.write "chipotle"
}))

server404: (http.testServer (func request response {
    response.sendStatus 404
}))

suite: {
  
    testClient: (http.client (http.testURL $server))
    testClient404: (http.client (http.testURL $server404))


    testHttpClientHasItsOwnType: (func {
        type $testClient |eq httpClient |assert
    })

    testHttpClientUsesUrlAsString: (func {
        to_s $testClient |str.startsWith "http://" |assert
    })

    testHttpClientCanGetContentAsString: (func {
        http.get $testClient |type |eq string |assert
    })

    testHttpClientWillReturnErrorOn404: (func {
        http.get $testClient404 "" |type |eq error |assert
        http.get $testClient "" |type |ne error |assert
    })

    testHttpClientCanGetContentByPath: (func {
        http.get $testClient "" |eq "chipotle" |assert
    })

    testHttpClientCanCatchCookies: (func {
        http.get $testClient ""
        dict.knows (http.cookies $testClient) pepper |assert
    })
}

test suite
close $server
close $server404