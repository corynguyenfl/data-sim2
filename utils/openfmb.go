package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/breakermodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/capbankmodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/commonmodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/essmodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/generationmodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/loadmodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/metermodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/reclosermodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/regulatormodule"
	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/solarmodule"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Voltage struct {
	Primary   float64
	Secondary float64
}

var LogMessageEnabled bool

func Publish(nc *nats.Conn, mrid string, profile protoreflect.ProtoMessage) {
	data, _ := proto.Marshal(profile)
	subject := fmt.Sprintf("%s.%s", getSubjectPrefix(profile), mrid)
	result := nc.Publish(subject, data)

	if result != nil {
		fmt.Println("There is an error publishing message")
	} else {
		if LogMessageEnabled {
			fmt.Println(profile)
		}
	}
}

func getSubjectPrefix(profile interface{}) string {
	return strings.Replace(reflect.TypeOf(profile).String(), "*", "openfmb.", 1)
}

func getDbPosKind(flag bool) commonmodule.DbPosKind {
	if flag == true {
		return commonmodule.DbPosKind_DbPosKind_closed
	}
	return commonmodule.DbPosKind_DbPosKind_open
}

func getStateKind(flag bool) commonmodule.StateKind {
	if flag == true {
		return commonmodule.StateKind_StateKind_on
	}
	return commonmodule.StateKind_StateKind_off
}

func getControlModeKind(manual bool) commonmodule.ControlModeKind {
	if manual == true {
		return commonmodule.ControlModeKind_ControlModeKind_manual
	}
	return commonmodule.ControlModeKind_ControlModeKind_auto
}

func CreateRecloserStatus(mrid string, pos bool) *reclosermodule.RecloserStatusProfile {
	return &reclosermodule.RecloserStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		Recloser: &reclosermodule.Recloser{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		RecloserStatus: &reclosermodule.RecloserStatus{
			StatusAndEventXCBR: &commonmodule.StatusAndEventXCBR{
				Pos: &commonmodule.PhaseDPS{
					Phs3: &commonmodule.StatusDPS{
						StVal: getDbPosKind(pos),
					},
				},
			},
		},
	}
}

func CreateRecloserReading(mrid string, va float64, vb float64, vc float64, w float64) *reclosermodule.RecloserReadingProfile {
	profile := &reclosermodule.RecloserReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		Recloser: &reclosermodule.Recloser{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
	}

	profile.RecloserReading = append(profile.RecloserReading, &reclosermodule.RecloserReading{
		ReadingMMXU: &commonmodule.ReadingMMXU{
			PhV: &commonmodule.WYE{
				PhsA: &commonmodule.CMV{
					CVal: &commonmodule.Vector{
						Mag: va,
					},
				},
				PhsB: &commonmodule.CMV{
					CVal: &commonmodule.Vector{
						Mag: vb,
					},
				},
				PhsC: &commonmodule.CMV{
					CVal: &commonmodule.Vector{
						Mag: vc,
					},
				},
			},
			W: &commonmodule.WYE{
				Net: &commonmodule.CMV{
					CVal: &commonmodule.Vector{
						Mag: w,
					},
				},
			},
		},
	})

	return profile
}

func CreateBreakerStatus(mrid string, pos bool) *breakermodule.BreakerStatusProfile {
	return &breakermodule.BreakerStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		Breaker: &breakermodule.Breaker{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		BreakerStatus: &breakermodule.BreakerStatus{
			StatusAndEventXCBR: &commonmodule.StatusAndEventXCBR{
				Pos: &commonmodule.PhaseDPS{
					Phs3: &commonmodule.StatusDPS{
						StVal: getDbPosKind(pos),
					},
				},
			},
		},
	}
}

func CreateBreakerReading(mrid string, wattage float64) *breakermodule.BreakerReadingProfile {
	profile := &breakermodule.BreakerReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		Breaker: &breakermodule.Breaker{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
	}

	profile.BreakerReading = append(profile.BreakerReading, &breakermodule.BreakerReading{
		ReadingMMXU: &commonmodule.ReadingMMXU{
			W: &commonmodule.WYE{
				Net: &commonmodule.CMV{
					CVal: &commonmodule.Vector{
						Mag: wattage,
					},
				},
			},
		},
	})

	return profile
}

func CreateEssStatus(mrid string, soc float64, isOn bool, mode int) *essmodule.ESSStatusProfile {
	return &essmodule.ESSStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		Ess: &commonmodule.ESS{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		EssStatus: &essmodule.ESSStatus{
			EssStatusZBAT: &essmodule.EssStatusZBAT{
				Soc: &commonmodule.MV{
					Mag: soc,
				},
			},
			EssStatusZGEN: &essmodule.ESSStatusZGEN{
				ESSEventAndStatusZGEN: &essmodule.ESSEventAndStatusZGEN{
					PointStatus: &essmodule.ESSPointStatus{
						State: &commonmodule.Optional_StateKind{
							Value: getStateKind(isOn),
						},
						Mode: &commonmodule.ENG_GridConnectModeKind{
							SetVal: commonmodule.GridConnectModeKind(mode),
						},
					},
				},
			},
		},
	}
}

func CreateLoadStatus(mrid string, isOn bool) *loadmodule.LoadStatusProfile {
	return &loadmodule.LoadStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		EnergyConsumer: &commonmodule.EnergyConsumer{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		LoadStatus: &loadmodule.LoadStatus{
			LoadStatusZGLD: &loadmodule.LoadStatusZGLD{
				LoadEventAndStatusZGLD: &loadmodule.LoadEventAndStatusZGLD{
					PointStatus: &loadmodule.LoadPointStatus{
						State: &commonmodule.Optional_StateKind{
							Value: getStateKind(isOn),
						},
					},
				},
			},
		},
	}
}

func CreateGeneratorStatus(mrid string, isOn bool) *generationmodule.GenerationStatusProfile {
	return &generationmodule.GenerationStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		GeneratingUnit: &generationmodule.GeneratingUnit{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		GenerationStatus: &generationmodule.GenerationStatus{
			GenerationStatusZGEN: &generationmodule.GenerationStatusZGEN{
				GenerationEventAndStatusZGEN: &generationmodule.GenerationEventAndStatusZGEN{
					PointStatus: &generationmodule.GenerationPointStatus{
						State: &commonmodule.Optional_StateKind{
							Value: getStateKind(isOn),
						},
					},
				},
			},
		},
	}
}

func CreateSolarStatus(mrid string, isOn bool) *solarmodule.SolarStatusProfile {
	return &solarmodule.SolarStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		SolarInverter: &solarmodule.SolarInverter{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		SolarStatus: &solarmodule.SolarStatus{
			SolarStatusZGEN: &solarmodule.SolarStatusZGEN{
				SolarEventAndStatusZGEN: &solarmodule.SolarEventAndStatusZGEN{
					PointStatus: &solarmodule.SolarPointStatus{
						State: &commonmodule.Optional_StateKind{
							Value: getStateKind(isOn),
						},
					},
				},
			},
		},
	}
}

func CreateGeneratorReading(mrid string, wattage float64) *generationmodule.GenerationReadingProfile {
	return &generationmodule.GenerationReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		GeneratingUnit: &generationmodule.GeneratingUnit{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		GenerationReading: &generationmodule.GenerationReading{
			ReadingMMXU: &commonmodule.ReadingMMXU{
				W: &commonmodule.WYE{
					Net: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: wattage,
						},
					},
				},
			},
		},
	}
}

func CreateSolarReading(mrid string, wattage float64) *solarmodule.SolarReadingProfile {
	return &solarmodule.SolarReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		SolarInverter: &solarmodule.SolarInverter{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		SolarReading: &solarmodule.SolarReading{
			ReadingMMXU: &commonmodule.ReadingMMXU{
				W: &commonmodule.WYE{
					Net: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: wattage,
						},
					},
				},
			},
		},
	}
}

func CreateMeterReading(mrid string, wattage float64) *metermodule.MeterReadingProfile {
	return &metermodule.MeterReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		Meter: &commonmodule.Meter{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		MeterReading: &metermodule.MeterReading{
			ReadingMMXU: &commonmodule.ReadingMMXU{
				W: &commonmodule.WYE{
					Net: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: wattage,
						},
					},
				},
			},
		},
	}
}

func CreateRegulatorStatus(
	mrid string,
	tapPos int32,
	volLmtHi bool,
	volLmtLo bool,
	voltageSetPointEnabled bool) *regulatormodule.RegulatorStatusProfile {
	return &regulatormodule.RegulatorStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		RegulatorSystem: &regulatormodule.RegulatorSystem{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		RegulatorStatus: &regulatormodule.RegulatorStatus{
			RegulatorEventAndStatusANCR: &regulatormodule.RegulatorEventAndStatusANCR{
				PointStatus: &regulatormodule.RegulatorEventAndStatusATCC{
					TapPos: &commonmodule.PhaseINS{
						Phs3: &commonmodule.StatusINS{
							StVal: tapPos,
						},
					},
					VolLmtHi: &commonmodule.PhaseSPS{
						Phs3: &commonmodule.StatusSPS{
							StVal: volLmtHi,
						},
					},
					VolLmtLo: &commonmodule.PhaseSPS{
						Phs3: &commonmodule.StatusSPS{
							StVal: volLmtLo,
						},
					},
					VoltageSetPointEnabled: &commonmodule.StatusSPS{
						StVal: voltageSetPointEnabled,
					},
				},
			},
		},
	}
}

func CreateRegulatorReading(
	mrid string,
	sourceVoltage *Voltage,
	loadVoltage *Voltage) *regulatormodule.RegulatorReadingProfile {
	profile := &regulatormodule.RegulatorReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		RegulatorSystem: &regulatormodule.RegulatorSystem{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
	}

	profile.RegulatorReading = append(profile.RegulatorReading, &regulatormodule.RegulatorReading{
		ReadingMMXU: &commonmodule.ReadingMMXU{
			PhV: &commonmodule.WYE{
				Net: &commonmodule.CMV{
					CVal: &commonmodule.Vector{
						Mag: sourceVoltage.Primary,
					},
				},
			},
		},
		SecondaryReadingMMXU: &commonmodule.ReadingMMXU{
			PhV: &commonmodule.WYE{
				Net: &commonmodule.CMV{
					CVal: &commonmodule.Vector{
						Mag: sourceVoltage.Secondary,
					},
				},
			},
		},
	})

	profile.RegulatorReading = append(profile.RegulatorReading, &regulatormodule.RegulatorReading{
		ReadingMMXU: &commonmodule.ReadingMMXU{
			PhV: &commonmodule.WYE{
				Net: &commonmodule.CMV{
					CVal: &commonmodule.Vector{
						Mag: loadVoltage.Primary,
					},
				},
			},
		},
		SecondaryReadingMMXU: &commonmodule.ReadingMMXU{
			PhV: &commonmodule.WYE{
				Net: &commonmodule.CMV{
					CVal: &commonmodule.Vector{
						Mag: loadVoltage.Secondary,
					},
				},
			},
		},
	})

	return profile
}

func CreateCapBankStatus(
	mrid string,
	controlMode int32,
	pos bool,
	volLmt bool,
	varLmt bool,
	tempLmt bool) *capbankmodule.CapBankStatusProfile {
	return &capbankmodule.CapBankStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		CapBankSystem: &capbankmodule.CapBankSystem{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		CapBankStatus: &capbankmodule.CapBankStatus{
			CapBankEventAndStatusYPSH: &capbankmodule.CapBankEventAndStatusYPSH{
				CtlMode: &commonmodule.Optional_ControlModeKind{
					Value: commonmodule.ControlModeKind(controlMode),
				},
				Pos: &commonmodule.PhaseDPS{
					Phs3: &commonmodule.StatusDPS{
						StVal: getDbPosKind(pos),
					},
				},
				VolLmt: &commonmodule.PhaseSPS{
					Phs3: &commonmodule.StatusSPS{
						StVal: volLmt,
					},
				},
				VArLmt: &commonmodule.PhaseSPS{
					Phs3: &commonmodule.StatusSPS{
						StVal: varLmt,
					},
				},
				TempLmt: &commonmodule.PhaseSPS{
					Phs3: &commonmodule.StatusSPS{
						StVal: tempLmt,
					},
				},
			},
		},
	}
}

func CreateCapBankReading(
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
	wattC float64) *capbankmodule.CapBankReadingProfile {
	return &capbankmodule.CapBankReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		CapBankSystem: &capbankmodule.CapBankSystem{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		CapBankReading: &capbankmodule.CapBankReading{
			ReadingMMXU: &commonmodule.ReadingMMXU{
				A: &commonmodule.WYE{
					PhsA: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: currentA,
						},
					},
					PhsB: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: currentB,
						},
					},
					PhsC: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: currentC,
						},
					},
				},
				PhV: &commonmodule.WYE{
					PhsA: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: voltageA,
						},
					},
					PhsB: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: voltageB,
						},
					},
					PhsC: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: voltageC,
						},
					},
				},
				W: &commonmodule.WYE{
					PhsA: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: wattA,
						},
					},
					PhsB: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: wattB,
						},
					},
					PhsC: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: wattC,
						},
					},
				},
			},
			SecondaryReadingMMXU: &commonmodule.ReadingMMXU{
				PhV: &commonmodule.WYE{
					PhsA: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: voltage2A,
						},
					},
					PhsB: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: voltage2B,
						},
					},
					PhsC: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: voltage2C,
						},
					},
				},
			},
		},
	}
}

func CreateLoadReading(
	mrid string,
	currentA float64,
	currentB float64,
	currentC float64,
	voltageA float64,
	voltageB float64,
	voltageC float64,
	va float64,
	vAr float64,
	watt float64) *loadmodule.LoadReadingProfile {
	return &loadmodule.LoadReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: Now(),
			},
		},
		EnergyConsumer: &commonmodule.EnergyConsumer{
			ConductingEquipment: &commonmodule.ConductingEquipment{
				MRID: mrid,
			},
		},
		LoadReading: &loadmodule.LoadReading{
			ReadingMMXU: &commonmodule.ReadingMMXU{
				A: &commonmodule.WYE{
					PhsA: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: currentA,
						},
					},
					PhsB: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: currentB,
						},
					},
					PhsC: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: currentC,
						},
					},
				},
				PhV: &commonmodule.WYE{
					PhsA: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: voltageA,
						},
					},
					PhsB: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: voltageB,
						},
					},
					PhsC: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: voltageC,
						},
					},
				},
				VA: &commonmodule.WYE{
					Net: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: va,
						},
					},
				},
				VAr: &commonmodule.WYE{
					Net: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: vAr,
						},
					},
				},
				W: &commonmodule.WYE{
					Net: &commonmodule.CMV{
						CVal: &commonmodule.Vector{
							Mag: watt,
						},
					},
				},
			},
		},
	}
}
