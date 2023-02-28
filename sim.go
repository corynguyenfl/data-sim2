package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	breakermodule "gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/breakermodule"
	commonmodule "gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/commonmodule"
	essmodule "gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/essmodule"
	generationmodule "gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/generationmodule"
	solarmodule "gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/solarmodule"
	loadmodule "gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/loadmodule"
	metermodule "gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/metermodule"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Nats struct {
		Url string `yaml:"url"`
	}
	Pcc struct {
		MRID     string  `yaml:"mrid"`
		W        float64 `yaml:"w"`
		IsClosed bool    `yaml:"is_closed"`
	}
	Ess struct {
		MRID        string  `yaml:"mrid"`
		ReadingMRID string  `yaml:"reading_mrid"`
		SOC         float64 `yaml:"soc"`
		Mode        int     `yaml:"mode"`
		IsOn        bool    `yaml:"is_on"`
		W           float64 `yaml:"w"`
	}
	Solar struct {
		MRID         string  `yaml:"mrid"`
		reading_mrid string  `yaml:"reading_mrid"`
		IsOn         bool    `yaml:"is_on"`
		W            float64 `yaml:"w"`
	}
	Generator struct {
		MRID        string  `yaml:"mrid"`
		ReadingMRID string  `yaml:"reading_mrid"`
		IsOn        bool    `yaml:"is_on"`
		W           float64 `yaml:"w"`
	}
	ShopMeter struct {
		MRID string  `yaml:"mrid"`
		W    float64 `yaml:"w"`
	} `yaml:"shop-meter"`
	LoadBank struct {
		MRID        string  `yaml:"mrid"`
		ReadingMRID string  `yaml:"reading_mrid"`
		IsOn        bool    `yaml:"is_on"`
		W           float64 `yaml:"w"`
	} `yaml:"load-bank"`
}

func main() {

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	fmt.Println("Press CTRL-C to exit...")

	config, _ := readConf("config/app.yaml")

	nc, _ := nats.Connect(config.Nats.Url)

	for {

		config, _ := readConf("config/app.yaml")

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
		 publishLoadBankStatus(nc, config.LoadBank.MRID, config.LoadBank.IsOn);
		 publishMeterReading(nc, config.LoadBank.MRID, config.ShopMeter.W)
 
		 // publish solar status		
		 publishSolarStatus(nc, config.Solar.MRID, config.Solar.IsOn);
		 publishSolarReading(nc, config.Solar.MRID, config.Solar.W);
 
		 // publish generator status
		 publishGeneratorStatus(nc, config.Generator.MRID, config.Generator.IsOn);
		 publishGeneratorReading(nc, config.Generator.MRID, config.Generator.W);

		time.Sleep(1 * time.Second)
	}
}

func readConf(defaultFileName string) (*Configuration, error) {

	filename := os.Getenv("APP_CONF")

	if filename == "" {
		filename = defaultFileName
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Configuration{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}

func now() *commonmodule.Timestamp {
	nano := time.Now().UTC().UnixNano()
	seconds := float64(nano) / 1e9
	return &commonmodule.Timestamp{
		Seconds:     uint64(seconds),
		Nanoseconds: 0,
	}
}

func publishPccStatus(nc *nats.Conn, mrid string, pos bool) {
	profile := createBreakerStatus(mrid, pos)
	data, _ := proto.Marshal(profile)

	result := nc.Publish(fmt.Sprintf("openfmb.breakermodule.BreakerStatusProfile.%s", mrid), data)

	if result != nil {
		fmt.Println("There is an error publishing message")
	} else {
		fmt.Println(profile)
	}
}

func publishPccReading(nc *nats.Conn, mrid string, wattage float64) {
	profile := createBreakerReading(mrid, wattage)
	data, _ := proto.Marshal(profile)

	result := nc.Publish(fmt.Sprintf("openfmb.breakermodule.BreakerReadingProfile.%s", mrid), data)

	if result != nil {
		fmt.Println("There is an error publishing message")
	} else {
		fmt.Println(profile)
	}
}

func publishEssStatus(nc *nats.Conn, mrid string, soc float64, isOn bool, mode int) {
	profile := createEssStatus(mrid, soc, isOn, mode)
	data, _ := proto.Marshal(profile)

	result := nc.Publish(fmt.Sprintf("openfmb.essmodule.ESSStatusProfile.%s", mrid), data)

	if result != nil {
		fmt.Println("There is an error publishing message")
	} else {
		fmt.Println(profile)
	}
}

func publishLoadBankStatus(nc *nats.Conn, mrid string, isOn bool) {
	profile := createLoadStatus(mrid, isOn)
	data, _ := proto.Marshal(profile)

	result := nc.Publish(fmt.Sprintf("openfmb.loadmodule.LoadStatusProfile.%s", mrid), data)

	if result != nil {
		fmt.Println("There is an error publishing message")
	} else {
		fmt.Println(profile)
	}
}

func publishSolarStatus(nc *nats.Conn, mrid string, isOn bool) {
	profile := createSolarStatus(mrid, isOn)
	data, _ := proto.Marshal(profile)

	result := nc.Publish(fmt.Sprintf("openfmb.solarmodule.SolarStatusProfile.%s", mrid), data)

	if result != nil {
		fmt.Println("There is an error publishing message")
	} else {
		fmt.Println(profile)
	}
}

func publishGeneratorStatus(nc *nats.Conn, mrid string, isOn bool) {
	profile := createGeneratorStatus(mrid, isOn)
	data, _ := proto.Marshal(profile)

	result := nc.Publish(fmt.Sprintf("openfmb.generationmodule.GenerationStatusProfile.%s", mrid), data)

	if result != nil {
		fmt.Println("There is an error publishing message")
	} else {
		fmt.Println(profile)
	}
}

func publishSolarReading(nc *nats.Conn, mrid string, wattage float64) {
	profile := createSolarReading(mrid, wattage)
	data, _ := proto.Marshal(profile)

	result := nc.Publish(fmt.Sprintf("openfmb.solarmodule.SolarReadingProfile.%s", mrid), data)

	if result != nil {
		fmt.Println("There is an error publishing message")
	} else {
		fmt.Println(profile)
	}
}

func publishGeneratorReading(nc *nats.Conn, mrid string, wattage float64) {
	profile := createGeneratorReading(mrid, wattage)
	data, _ := proto.Marshal(profile)

	result := nc.Publish(fmt.Sprintf("openfmb.generationmodule.GenerationReadingProfile.%s", mrid), data)

	if result != nil {
		fmt.Println("There is an error publishing message")
	} else {
		fmt.Println(profile)
	}
}

func publishMeterReading(nc *nats.Conn, mrid string, wattage float64) {
	profile := createMeterReading(mrid, wattage)
	data, _ := proto.Marshal(profile)

	result := nc.Publish(fmt.Sprintf("openfmb.metermodule.MeterReadingProfile.%s", mrid), data)

	if result != nil {
		fmt.Println("There is an error publishing message")
	} else {
		fmt.Println(profile)
	}
}

func createMeterReading(mrid string, wattage float64) *metermodule.MeterReadingProfile {
	return &metermodule.MeterReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: now(),
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

func createBreakerStatus(mrid string, pos bool) *breakermodule.BreakerStatusProfile {
	return &breakermodule.BreakerStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: now(),
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

func createBreakerReading(mrid string, wattage float64) *breakermodule.BreakerReadingProfile {
	profile := &breakermodule.BreakerReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: now(),
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

func createEssStatus(mrid string, soc float64, isOn bool, mode int) *essmodule.ESSStatusProfile {
	return &essmodule.ESSStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: now(),
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

func createLoadStatus(mrid string, isOn bool) *loadmodule.LoadStatusProfile {
	return &loadmodule.LoadStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: now(),
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

func createGeneratorStatus(mrid string, isOn bool) *generationmodule.GenerationStatusProfile {
	return &generationmodule.GenerationStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: now(),
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

func createSolarStatus(mrid string, isOn bool) *solarmodule.SolarStatusProfile {
	return &solarmodule.SolarStatusProfile{
		StatusMessageInfo: &commonmodule.StatusMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: now(),
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

func createGeneratorReading(mrid string, wattage float64) *generationmodule.GenerationReadingProfile {
	return &generationmodule.GenerationReadingProfile{
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: now(),
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

func createSolarReading(mrid string, wattage float64) *solarmodule.SolarReadingProfile {
	return &solarmodule.SolarReadingProfile {
		ReadingMessageInfo: &commonmodule.ReadingMessageInfo{
			MessageInfo: &commonmodule.MessageInfo{
				IdentifiedObject: &commonmodule.IdentifiedObject{
					MRID: &wrapperspb.StringValue{
						Value: uuid.New().String(),
					},
				},
				MessageTimeStamp: now(),
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
