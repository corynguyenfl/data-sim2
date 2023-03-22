package utils

import (
	"fmt"

	"gitlab.com/openfmb/psm/ops/protobuf/go-openfmb-ops-protobuf/v2/openfmb/breakermodule"
)

type MissingMridError struct {
}

func (e *MissingMridError) Error() string {
	return "MissingMridError"
}

func Mrid(message interface{}) (mrid string, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR:: Failed to extract MRID from message: ", message)
			err = r.(error)
		}
	}()

	profile, ok := message.(breakermodule.BreakerReadingProfile)
	if ok {
		mrid = profile.Breaker.ConductingEquipment.MRID
		return mrid, nil
	}

	return "", &MissingMridError{}
}
