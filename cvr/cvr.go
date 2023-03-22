package cvr

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/corynguyenfl/data-sim2/utils"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/capbankmodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/commonmodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/reclosermodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/regulatormodule"
)

type CVR struct {
	m          sync.Mutex
	configFile string
}

func (a *CVR) Start(configFile string) {
	a.configFile = configFile
	appconfig, _ := utils.ReadAppConfig(configFile)

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
		appconfig, _ := utils.ReadAppConfig(configFile)
		a.m.Unlock()

		config := appconfig.CvrConfiguration

		// Reclosers
		go publishRecloserStatus(nc, config.Recloser1.MRID, config.Recloser1.IsClosed)
		go publishRecloserStatus(nc, config.Recloser2.MRID, config.Recloser2.IsClosed)
		go publishRecloserReading(nc, config.Recloser1.MRID, config.Recloser1.Va, config.Recloser1.Vb, config.Recloser1.Vc, config.Recloser1.W)
		go publishRecloserReading(nc, config.Recloser2.MRID, config.Recloser2.Va, config.Recloser2.Vb, config.Recloser2.Vc, config.Recloser2.W)

		// Regulators
		go publishRegulatorStatus(nc, config.VR1.MRID, config.VR1.Pos, config.VR1.VolLmHi, config.VR1.VolLmLo, config.VR1.VoltageSetPointEnabled)
		go publishRegulatorStatus(nc, config.VR2.MRID, config.VR2.Pos, config.VR2.VolLmHi, config.VR2.VolLmLo, config.VR2.VoltageSetPointEnabled)
		go publishRegulatorStatus(nc, config.VR3.MRID, config.VR3.Pos, config.VR3.VolLmHi, config.VR3.VolLmLo, config.VR3.VoltageSetPointEnabled)

		go publishRegulatorReading(nc, config.VR1.MRID,
			&utils.Voltage{Primary: config.VR1.SourcePrimaryVolage, Secondary: config.VR1.SourceSecondaryVolage},
			&utils.Voltage{Primary: config.VR1.LoadPrimaryVolage, Secondary: config.VR1.LoadSecondaryVolage})

		go publishRegulatorReading(nc, config.VR2.MRID,
			&utils.Voltage{Primary: config.VR2.SourcePrimaryVolage, Secondary: config.VR2.SourceSecondaryVolage},
			&utils.Voltage{Primary: config.VR2.LoadPrimaryVolage, Secondary: config.VR2.LoadSecondaryVolage})

		go publishRegulatorReading(nc, config.VR3.MRID,
			&utils.Voltage{Primary: config.VR3.SourcePrimaryVolage, Secondary: config.VR3.SourceSecondaryVolage},
			&utils.Voltage{Primary: config.VR3.LoadPrimaryVolage, Secondary: config.VR3.LoadSecondaryVolage})

		// CapBank
		go publishCapbankStatus(nc, config.CapBank.MRID, config.CapBank.ControlMode, config.CapBank.IsClosed, config.CapBank.VolLmt, config.CapBank.VarLmt, config.CapBank.TempLmt)
		go publishCapbankReading(
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
		go publishLoadReading(
			nc,
			config.Load1.MRID,
			config.Load1.Ia,
			config.Load1.Ib,
			config.Load1.Ic,
			config.Load1.Va,
			config.Load1.Vb,
			config.Load1.Vc,
			config.Load1.Apparent,
			config.Load1.Reactive,
			config.Load1.W)

		go publishLoadReading(
			nc,
			config.Load2.MRID,
			config.Load2.Ia,
			config.Load2.Ib,
			config.Load2.Ic,
			config.Load2.Va,
			config.Load2.Vb,
			config.Load2.Vc,
			config.Load2.Apparent,
			config.Load2.Reactive,
			config.Load2.W)

		go publishLoadReading(
			nc,
			config.Load3.MRID,
			config.Load3.Ia,
			config.Load3.Ib,
			config.Load3.Ic,
			config.Load3.Va,
			config.Load3.Vb,
			config.Load3.Vc,
			config.Load3.Apparent,
			config.Load3.Reactive,
			config.Load3.W)

		go publishLoadReading(
			nc,
			config.Load4.MRID,
			config.Load4.Ia,
			config.Load4.Ib,
			config.Load4.Ic,
			config.Load4.Va,
			config.Load4.Vb,
			config.Load4.Vc,
			config.Load4.Apparent,
			config.Load4.Reactive,
			config.Load4.W)

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
		appconfig, _ := utils.ReadAppConfig(a.configFile)

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
		appconfig, _ := utils.ReadAppConfig(a.configFile)

		if appconfig.CvrConfiguration.CapBank.MRID == mrid {
			a.handleCapbankPos(&profile, appconfig)
			a.handleCapbankRemote(&profile, appconfig)
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

	if profile.CapBankControl != nil {
		if profile.CapBankControl.CapBankDiscreteControlYPSH != nil {
			if profile.CapBankControl.CapBankDiscreteControlYPSH.Control != nil {
				if profile.CapBankControl.CapBankDiscreteControlYPSH.Control.Pos != nil {
					if profile.CapBankControl.CapBankDiscreteControlYPSH.Control.Pos.Phs3 != nil {
						appconfig.CvrConfiguration.CapBank.IsClosed = profile.CapBankControl.CapBankDiscreteControlYPSH.Control.Pos.Phs3.CtlVal
						fmt.Println("Updated app config for capbank: IsClosed = ", appconfig.CvrConfiguration.CapBank.IsClosed)
						appconfig.Save()
					}
				}
			}
		}
	}

	return nil
}

func (a *CVR) handleCapbankRemote(profile *capbankmodule.CapBankDiscreteControlProfile, appconfig *utils.AppConfig) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("No CtlMode data in profile.")
			err = r.(error)
		}
	}()

	if profile.CapBankControl != nil {
		if profile.CapBankControl.CapBankDiscreteControlYPSH != nil {
			if profile.CapBankControl.CapBankDiscreteControlYPSH.Control != nil {
				if profile.CapBankControl.CapBankDiscreteControlYPSH.Control.CtlModeOvrRd != nil {
					if profile.CapBankControl.CapBankDiscreteControlYPSH.Control.CtlModeOvrRd.CtlVal {
						appconfig.CvrConfiguration.CapBank.ControlMode = 4
					} else {
						appconfig.CvrConfiguration.CapBank.ControlMode = 1
					}
					fmt.Println("Updated app config for capbank: ControlMode = ", appconfig.CvrConfiguration.CapBank.ControlMode)
					appconfig.Save()
				}
			}
		}
	}

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
		appconfig, _ := utils.ReadAppConfig(a.configFile)

		if appconfig.CvrConfiguration.VR1.MRID == mrid {
			a.handleDirMode(&profile, appconfig, &appconfig.CvrConfiguration.VR1)
			a.handleTapOp(&profile, appconfig, &appconfig.CvrConfiguration.VR1)
		} else if appconfig.CvrConfiguration.VR2.MRID == mrid {
			a.handleDirMode(&profile, appconfig, &appconfig.CvrConfiguration.VR2)
			a.handleTapOp(&profile, appconfig, &appconfig.CvrConfiguration.VR2)
		} else if appconfig.CvrConfiguration.VR3.MRID == mrid {
			a.handleDirMode(&profile, appconfig, &appconfig.CvrConfiguration.VR3)
			a.handleTapOp(&profile, appconfig, &appconfig.CvrConfiguration.VR3)
		}

		a.m.Unlock()

	}
	return nil
}

func (a *CVR) handleTapOp(profile *regulatormodule.RegulatorDiscreteControlProfile, appconfig *utils.AppConfig, vr *utils.VoltageRegulator) (err error) {
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
					fmt.Printf("Updated app config for  %s: Pos = %d\n", vr.Name, vr.Pos)
					appconfig.Save()
				}
			} else if profile.RegulatorDiscreteControl.RegulatorControlATCC.TapOpR != nil {
				if profile.RegulatorDiscreteControl.RegulatorControlATCC.TapOpR.Phs3 != nil {
					cmd := profile.RegulatorDiscreteControl.RegulatorControlATCC.TapOpR.Phs3.CtlVal
					if cmd {
						vr.Pos = vr.Pos + 1
					} else {
						vr.Pos = vr.Pos - 1
					}
					fmt.Printf("Updated app config for %s: Pos = %d\n", vr.Name, vr.Pos)
					appconfig.Save()
				}
			}
		}
	}

	return nil
}

func (a *CVR) handleDirMode(profile *regulatormodule.RegulatorDiscreteControlProfile, appconfig *utils.AppConfig, vr *utils.VoltageRegulator) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("No DirMode data in profile.")
			err = r.(error)
		}
	}()

	if profile.RegulatorDiscreteControl != nil {
		if profile.RegulatorDiscreteControl.RegulatorControlATCC != nil {
			if profile.RegulatorDiscreteControl.RegulatorControlATCC.DirMode != nil {
				mode := profile.RegulatorDiscreteControl.RegulatorControlATCC.DirMode.Value
				if mode == commonmodule.DirectionModeKind_DirectionModeKind_locked_forward {
					vr.VoltageSetPointEnabled = true
				} else {
					vr.VoltageSetPointEnabled = false
				}
				fmt.Printf("Updated app config for %s: VoltageSetPointEnabled = %t\n", vr.Name, vr.VoltageSetPointEnabled)
				appconfig.Save()
			}
		}
	}

	return nil
}

func publishRecloserStatus(nc *nats.Conn, mrid string, pos bool) {
	profile := utils.CreateRecloserStatus(mrid, pos)
	utils.Publish(nc, mrid, profile)
}

func publishRecloserReading(nc *nats.Conn, mrid string, va float64, vb float64, vc float64, w float64) {
	profile := utils.CreateRecloserReading(mrid, va, vb, vc, w)
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
	controlMode int32,
	pos bool,
	volLmt bool,
	varLmt bool,
	tempLmt bool) {
	profile := utils.CreateCapBankStatus(mrid, controlMode,
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
