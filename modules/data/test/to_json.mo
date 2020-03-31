data: (load data)

suite: {

    testStringToJSON: (func {
        data.toJSON "chipotle" |eq `"chipotle"` |assert
    })

    testDictionaryToJSON: (func {
        d: {
            Pepper: "Cayenne"
            shu: {
                min: 30000
                max: 50000
            }
            notAsHotAs: ["habanero", "santaka"]
        }
        json: (data.toJSON $d)
        eq $json `{"Pepper":"Cayenne","notAsHotAs":["habanero","santaka"],"shu":{"max":50000,"min":30000}}` |assert
    })
}

test suite