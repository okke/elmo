# Working with files

## Reading and writing files

Elmo comes with a build-in function to read files. The ``file`` function. It actually does a bit more than just reading, it provides information about a file returned as an elmo dictionary.

```elmo
readme: (file README.md)
puts $readme
```

The ``file`` function returns a dictionary with the following properties:

- ``absPath`` (gives the absolute path)
- ``exists`` (checks if file exists)
- ``isDir`` (checks if file is a directory)
- ``size`` (file size)
- ``mode`` (unix style string representation of the file mode. e.g. -rw-r--r--)
- ``name`` (well, the name of the file)
- ``path`` (the relative path)
- ``binary`` (a function returning the content as binary value)
- ``string`` (a function returning the content as string value)
- ``write`` (a function to write content to the file)
- ``append`` (a function to append content to the file)


### read

To retrieve the content of a file, simply use the returned ``string`` function:

```elmo
content: ((file README.md) string)
puts $content
```

### write

And to write something to a file, simply use the returned ``write`` function:

```elmo
f: (file data.txt)
f.write "nice!"
```

A simple trick to read content from a file, process it and write it to another file is by combining read and write in a pipe construction:

```elmo
str: (load string)
((file in.data) string) | (func s { str.replace all $s "a" "o" }) | ((file out.data) write)
```

### append

Writing to files will overwrite existing content. It's also possible to append content to a file using a file's ``append`` function.

```elmo
f: (file data.txt)
f.append "first line"
f.append "\n"
f.append "second line"
```

## Temporary files

Elmo has build in support to work with temporary fileswhich are removed after being using. This is done through the ``tempFile`` function. ``tempFile`` Creates a temporary file in the default directory for temporary files. It then executes the provided block of elmo code, passing the associated file dictionary as variable, and removed the file.

```elmo
result: (tempFile tmp {
    
    # do something with $tmp

    return $tmp.absPath
})
```

## Reading CSV files

Elmo comes with a ``data`` module which contains some handy functions to work with structured data. It, among other formats, support the parsing of comma separated value files.

```
data: (load data)
peppers: (((file "peppers.csv") string) |data.fromCSV)
```

Suppose the peppers.csv file contains data like

```csv
shumin, shumax, rating, name
855000, 1463000, 4.0, Ghost Pepper
876000, 970000, 4.0, Dorset Naga
```

Then the ``data.csv`` function (assuming the data module is loaded as ``data``), returns a list of dictionaries with the properties found on the first line of the CSV file. To assert this try:

```elmo
(peppers 0) name |eq "Ghost Pepper" |assert
(peppers 1) shumax |eq 970000 |assert
```

## Reading JSON files

The same kind of mechanism can be used to read a file containing JSON data.

```
data: (load data)
peppers: (((file "peppers.json") string) |data.fromJSON)
```

Here the ``data.json`` function does the trick. It assumes the input data contains a JSON object and converts it to an elmo dictionary.

## Producing CSV data

Elmo can convert a list of dictionaries into a string containing CSV data.

```elmo
data: (load data)

# Load csv data (which returns a list of dictionaries)
# 
peppers: (((file "./test/peppers.csv") string) |data.fromCSV)

# And produce csv again. But only use name and rating properties
#
puts (data.toCSV ["name" "rating"] $peppers)
```

## Producing JSON data

Elmo dictionaries (and other value types) can be easily converted to strings containing JSON data.

```elmo
data: {
    Pepper: "Cayenne"
    shu: {
        min: 30000
        max: 50000
    }
    notAsHotAs: ["habanero", "santaka"]
}
json: (data.toJSON $data)
puts $json
```

This will print one line of JSON (which looks almost the same as the original Elmo dictionary):

```elmo
{"Pepper":"Cayenne","notAsHotAs":["habanero","santaka"],"shu":{"max":50000,"min":30000}}
```

Note, the order of all properties is sorted by key name when converted to JSON.

