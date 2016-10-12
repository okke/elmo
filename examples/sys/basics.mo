
# directly access sys and list functions
#
mixin (load sys)
mixin (load list)

# when running an os command, we only get a command value
#
puts "ls command: " (ls)

# also when we pipe the output to another command
#
puts "ls piped to wc command: " (ls |wc)

# so we need to execute the actual command (exec is an elmo function from sys)
# which will give us a list (of lines)
#
puts "executed directly: " (ls |exec)

# same when piping
#
puts "piped: " (ls | wc "-l"| exec)

# we can also store commands and execute them later
#
cmd: (ls | wc)
puts "executed later: " (exec $cmd)

# but instead of using exec, we can use list functions
#
ls | each x {
  puts "element of list: " $x
}
