module basic

go 1.16

replace dev.local/osxkeystrokes => ./../../../osx_keystrokes/.

replace dev.local/utils => ./../../../utils/.

replace dev.local/types => ./../../../types/.

require (
	dev.local/osxkeystrokes v0.0.0-00010101000000-000000000000
	dev.local/utils v0.0.0-00010101000000-000000000000
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/mattn/go-colorable v0.1.8 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
)
