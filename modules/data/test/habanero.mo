data: (load data)

habanero: (((file "./test/habanero.json") string) |data.fromJSON)

type (habanero.name) |eq string |assert
type (habanero.shu) |eq list |assert
type (habanero.rating) |eq float |assert
type (habanero.intrating) |eq int |assert

habanero.shu 0 |eq 350000 |assert
habanero.shu 1 |eq 855000 |assert

(habanero.comments 0) from |eq joe |assert 
(habanero.comments 0) text |eq "pretty damn hot but still eatable" |assert 
(habanero.comments 1) from |eq han |assert 
(habanero.comments 1) text |eq "chews like a cucumber lollypop" |assert 
