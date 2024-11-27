module exampleplugin

go 1.23

replace github.com/okke/elmo => ../../..

require github.com/okke/elmo v0.0.0-20200205202441-fe94bce044de

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.13.0 // indirect
)
