
int_func: (func i {
  assert (eq (type $i) int) "first argument must be an int"

  puts "can perform integer calculation"

  return (plus $i $i)
})

puts (int_func 5)
puts (int_func "chipotle")
