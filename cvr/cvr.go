package cvr

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/openenergysolutions/data-sim/utils"
)

type CVR struct {
}

func (a CVR) Start() {
	appconfig, _ := utils.ReadAppConfig("config/app.yaml")

	nc, _ := nats.Connect(appconfig.Nats.Url)

	for {

		appconfig, _ := utils.ReadAppConfig("config/app.yaml")

		config := appconfig.CvrConfiguration

		// Reclosers
		publishRecloserStatus(nc, config.Recloser1.MRID, config.Recloser1.IsClosed)
		publishRecloserStatus(nc, config.Recloser2.MRID, config.Recloser1.IsClosed)
		publishRecloserReading(nc, config.Recloser1.MRID, config.Recloser1.W)
		publishRecloserReading(nc, config.Recloser2.MRID, config.Recloser2.W)

		publishRegulatorStatus(nc, config.VR1.MRID, config.VR1.Pos, config.VR1.VolLmHi, config.VR1.VolLmLo, config.VR1.VoltageSetPointEnabled)
		publishRegulatorStatus(nc, config.VR2.MRID, config.VR2.Pos, config.VR2.VolLmHi, config.VR2.VolLmLo, config.VR2.VoltageSetPointEnabled)
		publishRegulatorStatus(nc, config.VR3.MRID, config.VR3.Pos, config.VR3.VolLmHi, config.VR3.VolLmLo, config.VR3.VoltageSetPointEnabled)

		// Regulators
		publishRegulatorReading(nc, config.VR1.MRID,
			&utils.Voltage{Primary: config.VR1.SourcePrimaryVolage, Secondary: config.VR1.SourceSecondaryVolage},
			&utils.Voltage{Primary: config.VR1.LoadPrimaryVolage, Secondary: config.VR1.LoadSecondaryVolage})

		publishRegulatorReading(nc, config.VR2.MRID,
			&utils.Voltage{Primary: config.VR2.SourcePrimaryVolage, Secondary: config.VR2.SourceSecondaryVolage},
			&utils.Voltage{Primary: config.VR2.LoadPrimaryVolage, Secondary: config.VR2.LoadSecondaryVolage})

		publishRegulatorReading(nc, config.VR3.MRID,
			&utils.Voltage{Primary: config.VR3.SourcePrimaryVolage, Secondary: config.VR3.SourceSecondaryVolage},
			&utils.Voltage{Primary: config.VR3.LoadPrimaryVolage, Secondary: config.VR3.LoadSecondaryVolage})

		// CapBank
		publishCapbankStatus(nc, config.CapBank.MRID, config.CapBank.Manual, config.CapBank.IsClosed, config.CapBank.VolLmt, config.CapBank.VarLmt, config.CapBank.TempLmt)
		publishCapbankReading(
			nc,
			config.CapBank.MRID,
			config.CapBank.Ia,
			config.CapBank.Ib,
			config.CapBank.Ic,
			config.CapBank.Va,
			config.CapBank.Vb,
			config.CapBank.Vc,
			config.CapBank.V2a,
			config.CapBank.V2b,
			config.CapBank.V2c,
			config.CapBank.Wa,
			config.CapBank.Wb,
			config.CapBank.Wc)

		// Loads

		time.Sleep(1 * time.Second)
	}
}

func publishRecloserStatus(nc *nats.Conn, mrid string, pos bool) {
	profile := utils.CreateRecloserStatus(mrid, pos)
	utils.Publish(nc, mrid, profile)
}

func publishRecloserReading(nc *nats.Conn, mrid string, wattage float64) {
	profile := utils.CreateRecloserReading(mrid, wattage)
	utils.Publish(nc, mrid, profile)
}

func publishRegulatorStatus(nc *nats.Conn, mrid string, tapPos int32, volLmHi bool, volLmLo bool, voltageSetPointEnabled bool) {
	profile := utils.CreateRegulatorStatus(mrid, tapPos, volLmHi, volLmLo, voltageSetPointEnabled)
	utils.Publish(nc, mrid, profile)
}

func publishRegulatorReading(nc *nats.Conn, mrid string, source *utils.Voltage, load *utils.Voltage) {
	profile := utils.CreateRegulatorReading(mrid, source, load)
	utils.Publish(nc, mrid, profile)
}

func publishCapbankStatus(
	nc *nats.Conn,
	mrid string,
	manual bool,
	pos bool,
	volLmt bool,
	varLmt bool,
	tempLmt bool) {
	profile := utils.CreateCapBankStatus(mrid, manual,
		pos,
		volLmt,
		varLmt,
		tempLmt)
	utils.Publish(nc, mrid, profile)
}

func publishCapbankReading(
	nc *nats.Conn,
	mrid string,
	currentA float64,
	currentB float64,
	currentC float64,
	voltageA float64,
	voltageB float64,
	voltageC float64,
	voltage2A float64,
	voltage2B float64,
	voltage2C float64,
	wattA float64,
	wattB float64,
	wattC float64) {
	profile := utils.CreateCapBankReading(
		mrid,
		currentA,
		currentB,
		currentC,
		voltageA,
		voltageB,
		voltageC,
		voltage2A,
		voltage2B,
		voltage2C,
		wattA,
		wattB,
		wattC)
	utils.Publish(nc, mrid, profile)
}

func publishLoadReading(
	nc *nats.Conn,
	mrid string,
	currentA float64,
	currentB float64,
	currentC float64,
	voltageA float64,
	voltageB float64,
	voltageC float64,
	va float64,
	vAr float64,
	watt float64) {
	profile := utils.CreateLoadReading(
		mrid,
		currentA,
		currentB,
		currentC,
		voltageA,
		voltageB,
		voltageC,
		va,
		vAr,
		watt)
	utils.Publish(nc, mrid, profile)
}
