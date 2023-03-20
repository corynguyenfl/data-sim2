package cvr

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/openenergysolutions/data-sim/utils"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/capbankmodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/reclosermodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/regulatormodule"
)

type CVR struct {
	m sync.Mutex
}

func (a *CVR) Start() {
	appconfig, _ := utils.ReadAppConfig("config/app.yaml")

	nc, _ := nats.Connect(appconfig.Nats.Url)

	nc.Subscribe("openfmb.reclosermodule.RecloserDiscreteControlProfile.>", func(m *nats.Msg) {
		go a.processRecloserControl(m)
	})

	nc.Subscribe("openfmb.capbankmodule.CapBankDiscreteControlProfile.>", func(m *nats.Msg) {
		go a.processCapBankControl(m)
	})

	nc.Subscribe("openfmb.regulatormodule.RegulatorDiscreteControlProfile.>", func(m *nats.Msg) {
		go a.processRegulatorControl(m)
	})

	for {

		a.m.Lock()
		appconfig, _ := utils.ReadAppConfig("config/app.yaml")
		a.m.Unlock()

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

func (a *CVR) processRecloserControl(m *nats.Msg) (err error) {

	var profile reclosermodule.RecloserDiscreteControlProfile

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR:: Failed to process recloser control: ", profile)
			err = r.(error)
		}
	}()

	var data = m.Data
	err = proto.Unmarshal(data, &profile)

	if err != nil {
		fmt.Println("ERROR:: failed to decode message")
		return err
	} else {
		mrid := profile.Recloser.ConductingEquipment.MRID
		pos := profile.RecloserDiscreteControl.RecloserDiscreteControlXCBR.DiscreteControlXCBR.Pos.Phs3.CtlVal

		a.m.Lock()
		appconfig, _ := utils.ReadAppConfig("config/app.yaml")

		config := appconfig.CvrConfiguration

		if config.Recloser1.MRID == mrid {
			appconfig.CvrConfiguration.Recloser1.IsClosed = pos
			fmt.Println("Updated app config for recloser 1: IsClosed = ", pos)
		} else if config.Recloser2.MRID == mrid {
			appconfig.CvrConfiguration.Recloser2.IsClosed = pos
			fmt.Println("Updated app config for recloser 2: IsClosed = ", pos)
		}
		appconfig.Save()
		a.m.Unlock()

	}
	return nil
}

func (a *CVR) processCapBankControl(m *nats.Msg) (err error) {

	var profile capbankmodule.CapBankDiscreteControlProfile

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR:: Failed to process capbank control: ", profile)
			err = r.(error)
		}
	}()

	var data = m.Data
	err = proto.Unmarshal(data, &profile)

	if err != nil {
		fmt.Println("ERROR:: failed to decode message")
		return err
	} else {
		mrid := profile.CapBankSystem.ConductingEquipment.MRID

		a.m.Lock()
		appconfig, _ := utils.ReadAppConfig("config/app.yaml")

		if appconfig.CvrConfiguration.CapBank.MRID == mrid {
			a.handleCapbankPos(&profile, appconfig)
		}

		a.m.Unlock()

	}
	return nil
}

func (a *CVR) handleCapbankPos(profile *capbankmodule.CapBankDiscreteControlProfile, appconfig *utils.AppConfig) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("No POS data in profile.")
			err = r.(error)
		}
	}()

	pos := profile.CapBankControl.CapBankDiscreteControlYPSH.Control.Pos.Phs3.CtlVal

	appconfig.CvrConfiguration.CapBank.IsClosed = pos
	fmt.Println("Updated app config for capbank: IsClosed = ", pos)
	appconfig.Save()

	return nil
}

func (a *CVR) processRegulatorControl(m *nats.Msg) (err error) {

	var profile regulatormodule.RegulatorDiscreteControlProfile

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR:: Failed to process regulator control: ", profile)
			err = r.(error)
		}
	}()

	var data = m.Data
	err = proto.Unmarshal(data, &profile)

	if err != nil {
		fmt.Println("ERROR:: failed to decode message")
		return err
	} else {
		mrid := profile.RegulatorSystem.ConductingEquipment.MRID

		a.m.Lock()
		appconfig, _ := utils.ReadAppConfig("config/app.yaml")

		if appconfig.CvrConfiguration.VR1.MRID == mrid {
			a.handleRaiseTap(&profile, appconfig, &appconfig.CvrConfiguration.VR1)
		} else if appconfig.CvrConfiguration.VR2.MRID == mrid {
			a.handleRaiseTap(&profile, appconfig, &appconfig.CvrConfiguration.VR2)
		} else if appconfig.CvrConfiguration.VR3.MRID == mrid {
			a.handleRaiseTap(&profile, appconfig, &appconfig.CvrConfiguration.VR3)
		}

		a.m.Unlock()

	}
	return nil
}

func (a *CVR) handleRaiseTap(profile *regulatormodule.RegulatorDiscreteControlProfile, appconfig *utils.AppConfig, vr *utils.VoltageRegulator) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("No TapOpL or TapOpR data in profile.")
			err = r.(error)
		}
	}()

	if profile.RegulatorDiscreteControl != nil {
		if profile.RegulatorDiscreteControl.RegulatorControlATCC != nil {
			if profile.RegulatorDiscreteControl.RegulatorControlATCC.TapOpL != nil {
				if profile.RegulatorDiscreteControl.RegulatorControlATCC.TapOpL.Phs3 != nil {
					cmd := profile.RegulatorDiscreteControl.RegulatorControlATCC.TapOpL.Phs3.CtlVal
					if cmd {
						vr.Pos = vr.Pos - 1
					} else {
						vr.Pos = vr.Pos + 1
					}
				}
			} else if profile.RegulatorDiscreteControl.RegulatorControlATCC.TapOpR != nil {
				if profile.RegulatorDiscreteControl.RegulatorControlATCC.TapOpR.Phs3 != nil {
					cmd := profile.RegulatorDiscreteControl.RegulatorControlATCC.TapOpR.Phs3.CtlVal
					if cmd {
						vr.Pos = vr.Pos + 1
					} else {
						vr.Pos = vr.Pos - 1
					}
				}
			}
		}
	}
	fmt.Println("Updated app config for VR: Pos = ", vr.Pos)
	appconfig.Save()

	return nil
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

func ReadMessage(m *nats.Msg) (profile interface{}, err error) {
	var data = m.Data

	var parts = strings.Split(m.Subject, ".")
	var subj = parts[2]

	switch subj {
	case "RecloserDiscreteControlProfile":
		{
			var obj reclosermodule.RecloserDiscreteControlProfile
			err = proto.Unmarshal(data, &obj)
			profile = &obj
		}
		break
	}

	if err != nil {
		fmt.Println("ERROR:: error parsing proto obj: ", err)
	}

	return profile, err
}
