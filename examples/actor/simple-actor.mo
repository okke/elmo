actor: (load "actor")

# start an actor that will count messages
# endlessly
#
counter: (actor.new {

  puts (actor.current) " is counting"

  do {

    # receive a message but ignore its content
    #
    actor.receive

    # and simply count
    #
    incr count
    puts "count " (count)

  } while (true)

})

# start an actor that will print messages
# until it receives a 0 value
#
bg: (actor.new {

  puts (actor.current) " is receiving messages"

  do {
    v: (actor.receive)

    puts "received " (v)

    # tell the counter to count (send true as message)
    #
    actor.send counter

  } until (eq (v) 0)

})

# start an actor that sends some messages
#
actor.new {

  puts (actor.current) " is sending messages"

  actor.send bg 1
  actor.send bg 2
  actor.send bg 3
  actor.send bg 0
}

# give the actors some time to interact
#
sleep 1000
