http: (load http)

# https://raw.githubusercontent.com/okke/elmo/master/modules/elmohttp/test/getdata.txt

suite: {
    client: (http.client "https://raw.githubusercontent.com")

    testHttpClientHasItsOwnType: (func {
        type $client |eq httpClient |assert
    })

    testHttpClientUsesUrlAsString: (func {
        to_s $client |eq "https://raw.githubusercontent.com" |assert
    })

    testHttpClientCanGetContentAsString: (func {
        http.get $client |type |eq string |assert
    })

    testHttpClientWillReturnErrorOn404: (func {
        http.get $client "/okke/elmo/master/404.go" |type |eq error |assert
        http.get $client "/okke/elmo/master/README.md" |type |ne error |assert
    })

    testHttpClientCanGetContentByPath: (func {
        http.get $client "/okke/elmo/master/modules/elmohttp/test/getdata.txt" |eq "chipotle" |assert
    })
}

(suite.testHttpClientHasItsOwnType)
(suite.testHttpClientUsesUrlAsString)
(suite.testHttpClientCanGetContentAsString)
(suite.testHttpClientWillReturnErrorOn404)
(suite.testHttpClientCanGetContentByPath)
