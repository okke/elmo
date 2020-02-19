http: (load http)
dict: (load dict)

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

    # TODO: do not use github but own internal server
    #
    testHttpClientCanCatchCookies: (func {

        # first check cookies before request
        #
        github: (http.client "https://github.com")
        cookies: (http.cookies $github)
        dict.knows $cookies _gh_sess |not |assert

        # and then do a request and expect a github session cookie
        #
        http.get $github "/"
        cookies: (http.cookies $github)
        dict.knows $cookies _gh_sess |assert
    })
}

test suite