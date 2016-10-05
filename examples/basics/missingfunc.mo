
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

# declare a dictionary with functions
#
sauce: {

  # declare missing func for this specific dictionary
  #
  ?: (func name args {
    puts "chipotles will do"
  })

  # and let the ingredients be missing
  #
  ingredients: (func {
    this.guess
  })
}

# simply do some undefined calls
#
sauce.ingredients
sauce.anything
