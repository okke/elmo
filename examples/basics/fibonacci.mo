

fib: (func n {
  prev: -1
  result: 1
  sum: 0
  i: 0

  while (lte (i) (n)) {

    sum: (plus (result) (prev))
    prev: (result)
    result: (sum)

    incr i
  }

  return (result)
})


recfib: (func n {

  if (lte (n) 1) {
    return (n)
  }

  # very slow when n > 25 but it works
  #
  return (plus (recfib (minus (n) 1)) (recfib (minus (n) 2)))

})



puts "fibonacci non recursive"

puts (fib 1) "," (fib 2) "," (fib 3)
puts (fib 4) "," (fib 5) "," (fib 6)
puts (fib 7) "," (fib 8) "," (fib 9)

puts "fibonacci recursive"

puts (recfib 1) "," (recfib 2) "," (recfib 3)
puts (recfib 4) "," (recfib 5) "," (recfib 6)
puts (recfib 7) "," (recfib 8) "," (recfib 9)
