
data: (load data)

peppers: (((file "./test/peppers.csv") string) |data.csv)

first: (peppers 0)

type (first.name) |eq string |assert
type (first.SHUMIN) |eq int |assert
type (first.rating) |eq float |assert