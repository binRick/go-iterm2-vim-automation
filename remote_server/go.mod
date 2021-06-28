module remote_server

go 1.16

replace dev.local/types => ./../types/.

replace mrz.io/itermctl => ./../dist/go_itermctl

require (
	dev.local/types v0.0.0-00010101000000-000000000000 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/jackpal/gateway v1.0.7
	github.com/shirou/gopsutil/v3 v3.21.5
	github.com/sirupsen/logrus v1.8.1
)
