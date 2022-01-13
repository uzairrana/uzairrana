package MatchMaking

import (
	"context"
	"database/sql"

	//	"encoding/json"
	//	"github.com/heroiclabs/nakama/api"
	"github.com/heroiclabs/nakama-common/rtapi"
	"github.com/heroiclabs/nakama-common/runtime"
)

func MatchMakerBefore(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, envelope *rtapi.Envelope) (*rtapi.Envelope, error) {

	logger.Info("======GetNumericProperties====='%v'", envelope.GetMatchmakerAdd().GetNumericProperties())

	numericProps := envelope.GetMatchmakerAdd().GetNumericProperties()

	envelope.GetMatchmakerAdd().NumericProperties = make(map[string]float64)
	for k, v := range numericProps {
		logger.Info("%v======loop====='%v'", k, v)

		envelope.GetMatchmakerAdd().NumericProperties[k] = v
	}

	return envelope, nil
}
