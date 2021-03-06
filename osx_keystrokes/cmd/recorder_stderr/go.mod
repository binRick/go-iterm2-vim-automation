module recorder

go 1.16

replace dev.local/osxkeystrokes => ../../

replace dev.local/utils => ../../../utils/.

require (
	dev.local/osxkeystrokes v0.0.0-00010101000000-000000000000
	dev.local/utils v0.0.0-00010101000000-000000000000
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20210208195552-ff826a37aa15 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/paulbellamy/ratecounter v0.2.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.4.0
)
