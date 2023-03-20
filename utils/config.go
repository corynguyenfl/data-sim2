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
	Nats struct {
		Url string `yaml:"url"`
	}
	MicrogridConfiguration MicrogridConfiguration `yaml:"microgrid-controller"`
	CvrConfiguration       CvrConfiguration       `yaml:"cvr"`
}

type MicrogridConfiguration struct {
	Enabled bool `yaml:"enabled"`
	Pcc     struct {
		MRID     string `yaml:"mrid"`
		W        float64
		IsClosed bool `yaml:"is_closed"`
	}
	Ess struct {
		MRID        string  `yaml:"mrid"`
		ReadingMRID string  `yaml:"reading_mrid"`
		SOC         float64 `yaml:"soc"`
		Mode        int     `yaml:"mode"`
		IsOn        bool    `yaml:"is_on"`
		W           float64
	}
	Solar struct {
		MRID         string `yaml:"mrid"`
		reading_mrid string `yaml:"reading_mrid"`
		IsOn         bool   `yaml:"is_on"`
		W            float64
	}
	Generator struct {
		MRID        string `yaml:"mrid"`
		ReadingMRID string `yaml:"reading_mrid"`
		IsOn        bool   `yaml:"is_on"`
		W           float64
	}
	ShopMeter struct {
		MRID string `yaml:"mrid"`
		W    float64
	} `yaml:"shop-meter"`
	LoadBank struct {
		MRID        string `yaml:"mrid"`
		ReadingMRID string `yaml:"reading_mrid"`
		IsOn        bool   `yaml:"is_on"`
		W           float64
	} `yaml:"load-bank"`
}

type CvrConfiguration struct {
	Enabled   bool `yaml:"enabled"`
	Recloser1 struct {
		MRID     string `yaml:"mrid"`
		W        float64
		IsClosed bool `yaml:"is_closed"`
	} `yaml:"recloser1"`
	Recloser2 struct {
		MRID     string `yaml:"mrid"`
		W        float64
		IsClosed bool `yaml:"is_closed"`
	} `yaml:"recloser2"`
	VR1 struct {
		MRID                   string  `yaml:"mrid"`
		Pos                    int32   `yaml:"pos"`
		VolLmHi                bool    `yaml:"volLmHi"`
		VolLmLo                bool    `yaml:"volLmLo"`
		VoltageSetPointEnabled bool    `yaml:"voltageSetPointEnabled"`
		SourcePrimaryVolage    float64 `yaml:"source_primary_voltage"`
		SourceSecondaryVolage  float64 `yaml:"source_secondary_voltage"`
		LoadPrimaryVolage      float64 `yaml:"load_primary_voltage"`
		LoadSecondaryVolage    float64 `yaml:"load_secondary_voltage"`
	} `yaml:"vr1"`
	VR2 struct {
		MRID                   string  `yaml:"mrid"`
		Pos                    int32   `yaml:"pos"`
		VolLmHi                bool    `yaml:"volLmHi"`
		VolLmLo                bool    `yaml:"volLmLo"`
		VoltageSetPointEnabled bool    `yaml:"voltageSetPointEnabled"`
		SourcePrimaryVolage    float64 `yaml:"source_primary_voltage"`
		SourceSecondaryVolage  float64 `yaml:"source_secondary_voltage"`
		LoadPrimaryVolage      float64 `yaml:"load_primary_voltage"`
		LoadSecondaryVolage    float64 `yaml:"load_secondary_voltage"`
	} `yaml:"vr2"`
	VR3 struct {
		MRID                   string  `yaml:"mrid"`
		Pos                    int32   `yaml:"pos"`
		VolLmHi                bool    `yaml:"volLmHi"`
		VolLmLo                bool    `yaml:"volLmLo"`
		VoltageSetPointEnabled bool    `yaml:"voltageSetPointEnabled"`
		SourcePrimaryVolage    float64 `yaml:"source_primary_voltage"`
		SourceSecondaryVolage  float64 `yaml:"source_secondary_voltage"`
		LoadPrimaryVolage      float64 `yaml:"load_primary_voltage"`
		LoadSecondaryVolage    float64 `yaml:"load_secondary_voltage"`
	} `yaml:"vr3"`
	CapBank struct {
		MRID     string `yaml:"mrid"`
		Manual   bool   `yaml:"manual"`
		IsClosed bool   `yaml:"is_closed"`
		VolLmt   bool   `yaml:"volLmt"`
		VarLmt   bool   `yaml:"varLmt"`
		TempLmt  bool   `yaml:"tempLmt"`
		Ia       float64
		Ib       float64
		Ic       float64
		Va       float64
		Vb       float64
		Vc       float64
		V2a      float64
		V2b      float64
		V2c      float64
		Wa       float64
		Wb       float64
		Wc       float64
	} `yaml:"capbank"`
	Load1 struct {
		MRID     string `yaml:"mrid"`
		Ia       float64
		Ib       float64
		Ic       float64
		Va       float64
		Vb       float64
		Vc       float64
		Apparent float64
		Reactive float64
		W        float64
	} `yaml:"load1"`
	Load2 struct {
		MRID     string `yaml:"mrid"`
		Ia       float64
		Ib       float64
		Ic       float64
		Va       float64
		Vb       float64
		Vc       float64
		Apparent float64
		Reactive float64
		W        float64
	} `yaml:"load2"`
	Load3 struct {
		MRID     string `yaml:"mrid"`
		Ia       float64
		Ib       float64
		Ic       float64
		Va       float64
		Vb       float64
		Vc       float64
		Apparent float64
		Reactive float64
		W        float64
	} `yaml:"load3"`
	Load4 struct {
		MRID     string `yaml:"mrid"`
		Ia       float64
		Ib       float64
		Ic       float64
		Va       float64
		Vb       float64
		Vc       float64
		Apparent float64
		Reactive float64
		W        float64
	} `yaml:"load4"`
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
