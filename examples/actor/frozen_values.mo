actor: (load actor)
list: (load list)

# start an actor that will modify a list
#
modify: (actor.new {

  puts (actor.current) " will modify incoming lists"

  do {

    # receive a value
    #
    value: (actor.receive)

    puts (actor.current) " has received incoming value " $value

    # will result in an error
    #
    result: (list.append! $value "unsafe")

    # but resulting error is not fatal so it can be printed
    #
    puts "result: " $result

  } while (true)

})

# create an array
#
values: [1 2 3]

# as long as this array is not send to an actor, it can be modified
#
list.append! $values 4

# but not after it has been send
#
puts "send value " $values
actor.send $modify $values

err: (list.append! $values 5)
puts "result after sending to actor: " $err

# give the actors some time to interact
#
sleep 1000

puts "after sending the list looks like " $values
