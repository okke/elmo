http: (load http)
dict: (load dict)
str: (load string)

server: (http.testServer (func request response {
    response.write "chipotle"        
}))

suite: {
  
    client: (http.client "https://raw.githubusercontent.com")
    testClient: (http.client (http.testURL $server))


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
        http.get $client "/okke/elmo/master/404.go" |type |eq error |assert
        http.get $client "/okke/elmo/master/README.md" |type |ne error |assert
    })

    testHttpClientCanGetContentByPath: (func {
        http.get $testClient "" |eq "chipotle" |assert
    })

    # TODO: do not use github but own internal server
    #
    testHttpClientCanCatchCookies: (func {

        # first check cookies before request
        #
        github: (http.client "https://github.com")
        dict.knows (http.cookies $github) _gh_sess |not |assert

        # and then do a request and expect a github session cookie
        #
        http.get $github "/"
        dict.knows (http.cookies $github) _gh_sess |assert
    })
}

test suite
close $server