http: (load http)
dict: (load dict)
list: (load list)
str: (load string)

server: (http.testServer (func request response {
    response.sendCookie "pepper" "jalapeno" 3600
    response.write "chipotle"
}))

server404: (http.testServer (func request response {
    response.sendStatus 404
}))

serverEcho: (http.testServer (func request response {
    dict.keys request |list.each k {
        response.write $k "=" (dict.get request $k |first) ";"
    }
}))


suite: {
  
    testClient: (http.client (http.testURL $server))
    testClient404: (http.client (http.testURL $server404))
    testClientEcho: (http.client (http.testURL $serverEcho))


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

    testHttpClientCanSendQueryParameters: (func {
        p: {pepper:"jalapeno"}
        http.get $testClientEcho "" $p |eq "pepper=jalapeno;" |assert
    })

    testHttpClientCanCatchCookies: (func {
        http.get $testClient ""
        dict.knows (http.cookies $testClient) pepper |assert
    })

    testHttpClientCanSendPostRequest: (func {
        http.post $testClientEcho "jalapeno" "" |eq "body=jalapeno;" |assert       
    })
}

test suite
close $server
close $server404
close $serverEcho