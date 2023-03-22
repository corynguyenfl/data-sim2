package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/commonmodule"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	LogMessageEnabled bool `yaml:"log-message-enabled"`
	Nats              struct {
		Url string `yaml:"url"`
	} `yaml:"nats"`
	MicrogridConfiguration MicrogridConfiguration `yaml:"microgrid-controller"`
	CvrConfiguration       CvrConfiguration       `yaml:"cvr"`
	file                   string
}

type MicrogridConfiguration struct {
	Enabled bool `yaml:"enabled"`
	Pcc     struct {
		MRID     string  `yaml:"mrid"`
		W        float64 `yaml:"W"`
		IsClosed bool    `yaml:"is_closed"`
	}
	Ess struct {
		MRID        string  `yaml:"mrid"`
		ReadingMRID string  `yaml:"reading_mrid"`
		SOC         float64 `yaml:"soc"`
		Mode        int     `yaml:"mode"`
		IsOn        bool    `yaml:"is_on"`
		W           float64 `yaml:"W"`
	}
	Solar struct {
		MRID         string `yaml:"mrid"`
		reading_mrid string `yaml:"reading_mrid"`
		IsOn         bool   `yaml:"is_on"`
		W            float64
	}
	Generator struct {
		MRID        string  `yaml:"mrid"`
		ReadingMRID string  `yaml:"reading_mrid"`
		IsOn        bool    `yaml:"is_on"`
		W           float64 `yaml:"W"`
	}
	ShopMeter struct {
		MRID string  `yaml:"mrid"`
		W    float64 `yaml:"W"`
	} `yaml:"shop-meter"`
	LoadBank struct {
		MRID        string  `yaml:"mrid"`
		ReadingMRID string  `yaml:"reading_mrid"`
		IsOn        bool    `yaml:"is_on"`
		W           float64 `yaml:"W"`
	} `yaml:"load-bank"`
}

type CvrConfiguration struct {
	Enabled   bool             `yaml:"enabled"`
	Recloser1 Recloser         `yaml:"recloser1"`
	Recloser2 Recloser         `yaml:"recloser2"`
	VR1       VoltageRegulator `yaml:"vr1"`
	VR2       VoltageRegulator `yaml:"vr2"`
	VR3       VoltageRegulator `yaml:"vr3"`
	CapBank   CapBank          `yaml:"capbank"`
	Load1     Load             `yaml:"load1"`
	Load2     Load             `yaml:"load2"`
	Load3     Load             `yaml:"load3"`
	Load4     Load             `yaml:"load4"`
}

type VoltageRegulator struct {
	Name                   string  `yaml:"name"`
	MRID                   string  `yaml:"mrid"`
	Pos                    int32   `yaml:"pos"`
	VolLmHi                bool    `yaml:"volLmHi"`
	VolLmLo                bool    `yaml:"volLmLo"`
	VoltageSetPointEnabled bool    `yaml:"voltageSetPointEnabled"`
	SourcePrimaryVolage    float64 `yaml:"source_primary_voltage"`
	SourceSecondaryVolage  float64 `yaml:"source_secondary_voltage"`
	LoadPrimaryVolage      float64 `yaml:"load_primary_voltage"`
	LoadSecondaryVolage    float64 `yaml:"load_secondary_voltage"`
}

type Load struct {
	Name     string  `yaml:"name"`
	MRID     string  `yaml:"mrid"`
	Ia       float64 `yaml:"Ia"`
	Ib       float64 `yaml:"Ib"`
	Ic       float64 `yaml:"Ic"`
	Va       float64 `yaml:"Va"`
	Vb       float64 `yaml:"Vb"`
	Vc       float64 `yaml:"Vc"`
	Apparent float64 `yaml:"Apparent"`
	Reactive float64 `yaml:"Reactive"`
	W        float64 `yaml:"W"`
}

type Recloser struct {
	Name     string  `yaml:"name"`
	MRID     string  `yaml:"mrid"`
	Va       float64 `yaml:"Va"`
	Vb       float64 `yaml:"Vb"`
	Vc       float64 `yaml:"Vc"`
	W        float64 `yaml:"W"`
	IsClosed bool    `yaml:"is_closed"`
}

type CapBank struct {
	Name        string  `yaml:"name"`
	MRID        string  `yaml:"mrid"`
	ControlMode int32   `yaml:"control-mode"`
	IsClosed    bool    `yaml:"is_closed"`
	VolLmt      bool    `yaml:"volLmt"`
	VarLmt      bool    `yaml:"varLmt"`
	TempLmt     bool    `yaml:"tempLmt"`
	Ia          float64 `yaml:"Ia"`
	Ib          float64 `yaml:"Ib"`
	Ic          float64 `yaml:"Ic"`
	Va          float64 `yaml:"Va"`
	Vb          float64 `yaml:"Vb"`
	Vc          float64 `yaml:"Vc"`
	V2a         float64 `yaml:"V2a"`
	V2b         float64 `yaml:"V2b"`
	V2c         float64 `yaml:"V2c"`
	Wa          float64 `yaml:"Wa"`
	Wb          float64 `yaml:"Wb"`
	Wc          float64 `yaml:"Wc"`
}

func (c *AppConfig) Save() {
	data, err := yaml.Marshal(c)

	if err != nil {
		fmt.Println("ERROR:: failed to serialize app config")
	}

	err = ioutil.WriteFile(c.file, data, 0)

	if err != nil {
		fmt.Println("ERROR:: failed to write app config")
	}
}

func ReadAppConfig(defaultFileName string) (*AppConfig, error) {

	filename := os.Getenv("APP_CONF")

	if filename == "" {
		filename = defaultFileName
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &AppConfig{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}
	c.file = filename
	return c, err
}

func Now() *commonmodule.Timestamp {
	nano := time.Now().UTC().UnixNano()
	seconds := float64(nano) / 1e9
	return &commonmodule.Timestamp{
		Seconds:     uint64(seconds),
		Nanoseconds: 0,
	}
}
