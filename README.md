# ELMO, embeddable and extensible scripting language for Go


![elmologo](/doc/images/logo.jpg)

Note, Elmo is still in draft status and under development (on master branch). But feel free to experiment with it. Or even help developing it. I started Elmo, back in the summer of 2016, as a pet project. And after one year of development, my interests shifted away to other things. But at the end of 2019, somehow I felt the urge to work on it again. Or better explained, to spend some leisure time coding again. So, it's still a pet project which maybe someday grow into a mature scripting language that can be used to get things done with a smile.

## Documentation

[introduction](/doc/introduction.md)

[far from complete manual](/doc/manual.md)

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
bin:  (load bin)
data: (load data)
http: (load http)
help
```

do not want to run the REPL but just want to execute an elmo source file:

```bash
elmo example.mo
```

want to run a file and then start the REPL:

```bash
elmo -repl example.mo
```
