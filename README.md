# ELMO, embeddable scripting language for Go


![elmologo](/doc/images/logo.jpg)

Note, Elmo is 'brand new', far from complete and still under heavy development. But feel free to experiment with it. Or even help developing it.

## Documentation

[introduction](/doc/introduction.md)

[manual](/doc/manual.md)

## Quick start

Install elmo using Go:

```bash
go get github.com/okke/elmo/tools/elmo
go install github.com/okke/elmo/tools/elmo
elmo
```

and you should see:

```
e>mo:
```

to see all basic commands type:

```elmo
help
```

to get help on specific command:
```elmo
help puts
```

want to see more:
```elmo
list: (load list)
dict: (load dict)
str:  (load string)
sys:  (load sys)
act:  (load actor)
help
```

do not want to run the REPL but execute an elmo source file:

```bash
elmo example.mo
```
