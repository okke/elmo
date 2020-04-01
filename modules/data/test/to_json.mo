data: (load data)

suite: {

    testStringToJSON: (func {
        data.toJSON "chipotle" |eq `"chipotle"` |assert
    })

    testDictionaryToJSON: (func {
        d: {
            shu: {
                min: 30000
                max: 50000
            }
            notAsHotAs: ["habanero", "santaka"]
            Pepper: "Cayenne"
        }
        json: (data.toJSON $d)
        eq $json `{"Pepper":"Cayenne","notAsHotAs":["habanero","santaka"],"shu":{"max":50000,"min":30000}}` |assert
    })
}

test suite