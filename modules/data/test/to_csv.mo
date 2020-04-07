data: (load data)

suite: {

    
    testStringToCSV: (func {

        expected: `name,rating
Ghost Pepper,4.0
Dorset Naga,4.0
Red Savina Habanero,4.0
Scotch Bonnet,3.5
Santaka,3.0
Cayenne,3.0
Manzano,2.5
Yellow Wax Pepper,2.0
`

        # Load csv data
        # 
        peppers: (((file "./test/peppers.csv") string) |data.fromCSV)
        
        # And produce csv again. But only use name and rating properties
        #
        csv: (data.toCSV ["name" "rating"] $peppers)

        assert (eq $csv $expected)
    })
}

test suite