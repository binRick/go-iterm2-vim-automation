module osxkeystrokes

go 1.16

replace dev.local/types => ./../types/.

replace dev.local/utils => ./../utils/.

require (
	dev.local/utils v0.0.0-00010101000000-000000000000
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/hpcloud/tail v1.0.0
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/mattn/go-colorable v0.1.8 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)
