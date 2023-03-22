package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/corynguyenfl/data-sim2/cvr"
	"github.com/corynguyenfl/data-sim2/microgrid"
	"github.com/corynguyenfl/data-sim2/utils"
)

func main() {

	appConfig, _ := utils.ReadAppConfig("config/app.yaml")

	utils.LogMessageEnabled = appConfig.LogMessageEnabled

	if appConfig.MicrogridConfiguration.Enabled {
		microgrid := &microgrid.Microgrid{}

		go microgrid.Start()
	}

	if appConfig.CvrConfiguration.Enabled {
		cvr := &cvr.CVR{}

		go cvr.Start()
	}

	fmt.Println("Press CTRL-C to exit...")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	func() {
		<-c
		os.Exit(1)
	}()
}
