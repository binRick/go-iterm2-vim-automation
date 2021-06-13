module local_server

go 1.16

replace dev.local/types => ./../types/.

replace mrz.io/itermctl => ./../dist/go_itermctl

require (
	dev.local/types v0.0.0-00010101000000-000000000000
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/k0kubun/pp v3.0.1+incompatible
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d
	github.com/pterm/pterm v0.12.23
	github.com/shirou/gopsutil v3.21.5+incompatible
	github.com/sirupsen/logrus v1.6.0
	github.com/tklauser/go-sysconf v0.3.6 // indirect
	mrz.io/itermctl v0.0.3
)

replace dev.local/utils => ./../utils/.
