module vimseqserver

go 1.16

replace mrz.io/itermctl => ./../../dist/go_itermctl/.

require (
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d
	github.com/sirupsen/logrus v1.6.0
	mrz.io/itermctl v0.0.3
)

replace dev.local/localserver => ./../../local_server/.
