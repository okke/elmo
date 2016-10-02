
# declare func missing function
#
?: (func name args {
  puts "no idea how to resolve " $name " " $args
  return $name
})

# and simply do some undefined calls
#
chipotle
chipotle "sauce"
chipotle "sauce" | clueless
