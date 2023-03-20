package microgrid

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/openenergysolutions/data-sim/utils"
)

type Microgrid struct {
	m sync.Mutex
}

func (m *Microgrid) Start() {

	appconfig, _ := utils.ReadAppConfig("config/app.yaml")

	nc, _ := nats.Connect(appconfig.Nats.Url)

	for {

		appconfig, _ := utils.ReadAppConfig("config/app.yaml")

		config := appconfig.MicrogridConfiguration

		// publish PCC status and PCC reading (PCC reading is also used for ESS reading)
		publishPccReading(nc, config.Pcc.MRID, config.Pcc.W)
		publishPccStatus(nc, config.Pcc.MRID, config.Pcc.IsClosed)

		// publish ESS status
		soc := config.Ess.SOC
		isOn := config.Ess.IsOn
		mode := config.Ess.Mode
		publishEssStatus(nc, config.Ess.MRID, soc, isOn, mode)

		// publish load reading (shop meter)
		publishMeterReading(nc, config.ShopMeter.MRID, config.ShopMeter.W)

		// publish loadbank reading (meter for the loadbank)
		publishLoadBankStatus(nc, config.LoadBank.MRID, config.LoadBank.IsOn)
		publishMeterReading(nc, config.LoadBank.MRID, config.ShopMeter.W)

		// publish solar status
		publishSolarStatus(nc, config.Solar.MRID, config.Solar.IsOn)
		publishSolarReading(nc, config.Solar.MRID, config.Solar.W)

		// publish generator status
		publishGeneratorStatus(nc, config.Generator.MRID, config.Generator.IsOn)
		publishGeneratorReading(nc, config.Generator.MRID, config.Generator.W)

		time.Sleep(1 * time.Second)
	}
}

func publishPccStatus(nc *nats.Conn, mrid string, pos bool) {
	profile := utils.CreateBreakerStatus(mrid, pos)
	utils.Publish(nc, mrid, profile)
}

func publishPccReading(nc *nats.Conn, mrid string, wattage float64) {
	profile := utils.CreateBreakerReading(mrid, wattage)
	utils.Publish(nc, mrid, profile)
}

func publishEssStatus(nc *nats.Conn, mrid string, soc float64, isOn bool, mode int) {
	profile := utils.CreateEssStatus(mrid, soc, isOn, mode)
	utils.Publish(nc, mrid, profile)
}

func publishLoadBankStatus(nc *nats.Conn, mrid string, isOn bool) {
	profile := utils.CreateLoadStatus(mrid, isOn)
	utils.Publish(nc, mrid, profile)
}

func publishSolarStatus(nc *nats.Conn, mrid string, isOn bool) {
	profile := utils.CreateSolarStatus(mrid, isOn)
	utils.Publish(nc, mrid, profile)
}

func publishGeneratorStatus(nc *nats.Conn, mrid string, isOn bool) {
	profile := utils.CreateGeneratorStatus(mrid, isOn)
	utils.Publish(nc, mrid, profile)
}

func publishSolarReading(nc *nats.Conn, mrid string, wattage float64) {
	profile := utils.CreateSolarReading(mrid, wattage)
	utils.Publish(nc, mrid, profile)
}

func publishGeneratorReading(nc *nats.Conn, mrid string, wattage float64) {
	profile := utils.CreateGeneratorReading(mrid, wattage)
	utils.Publish(nc, mrid, profile)
}

func publishMeterReading(nc *nats.Conn, mrid string, wattage float64) {
	profile := utils.CreateMeterReading(mrid, wattage)
	utils.Publish(nc, mrid, profile)
}
