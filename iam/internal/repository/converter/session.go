package converter

import (
	"encoding/json"

	servicemodel "iam/internal/model"
)

func SessionToJSON(session *servicemodel.Session) ([]byte, error) {
	return json.Marshal(session)
}

func JSONToSession(data []byte) (*servicemodel.Session, error) {
	var session servicemodel.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}