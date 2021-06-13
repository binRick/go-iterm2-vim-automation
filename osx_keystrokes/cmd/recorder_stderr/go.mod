module recorder

go 1.16

require (
	dev.local/osxkeystrokes v0.0.0-00010101000000-000000000000
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20210208195552-ff826a37aa15 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

replace dev.local/osxkeystrokes => ./../../../osx_keystrokes/.

replace dev.local/utils => ./../../../utils/.

replace dev.local/types => ./../../../types/.
