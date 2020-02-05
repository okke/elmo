# example to show how to load a static library in elmo

## compile the library

```
cd gosrc
go build -buildmode=plugin
```

This will produce the exampleplugin.so

## run the example code

```
elmo main.mo
```
