# Working with files

## Reading files

Elmo comes with a build-in function to read files. The ``file`` function. It actually does a bit more than just reading, it provides information about a file returned as an elmo dictionary.

```elmo
readme: (file README.md)
puts $readme
```

The ``file`` function returns a dictionary with the following properties:

- ``absPath`` (gives the absolute path)
- ``exists`` (checks if file exists)
- ``isDir`` (checks if file is a directory)
- ``mode`` (unix style string representation of the file mode. e.g. -rw-r--r--)
- ``name`` (well, the name of the file)
- ``path`` (the relative path)
- ``binary`` (a function returning the content as binary value)
- ``string`` (a function returning the content as string value)

To retrieve the content of a file, simply use the returned ``string`` function:

```elmo
content: ((file README.md) string)
puts $content
```

## Reading CSV files

Elmo comes with a ``data`` module which contains some handy functions to work with structured data. It, among other formats, support the parsing of comma separated value files.

```
data: (load data)
peppers: (((file "peppers.csv") string) |data.csv)
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
peppers: (((file "peppers.json") string) |data.json)
```

Here the ``data.json`` function does the trick. It assumes the input data contains a JSON object and converts it to an elmo dictionary.
