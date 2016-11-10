
str: (load "string")

soup: (func {
  return "lovely beans with jalapeno soup"
})

funkysoup: (func pepper {
  message: (str.concat "lovely beans with " $pepper " soup")
  return (func {
    return $message
  })
})

puts "regular soup " $soup
puts "funky soup " ((funkysoup "chipotle"))
