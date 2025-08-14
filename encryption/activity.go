package encryption

import (
	"encoding/json"
	"time"
)

type UserActivity struct {
	Mousemoved         bool   `json:"mouseMoved"`
	Keyboardused       bool   `json:"keyboardUsed"`
	Fieldreceivedinput bool   `json:"fieldReceivedInput"`
	Fieldreceivedfocus bool   `json:"fieldReceivedFocus"`
	Timestamp          string `json:"timestamp"`
	Email              string `json:"email"`
}

func GenerateActivity(email string) (string, error) {
	activity := UserActivity{
		Mousemoved:         true,
		Keyboardused:       true,
		Fieldreceivedinput: true,
		Fieldreceivedfocus: true,
		Timestamp:          time.Now().Format("2006-01-02T15:04:05.000Z"),
		Email:              email,
	}

	activityJson, err := json.Marshal(activity)
	if err != nil {
		return "", err
	}

	return string(activityJson), nil
}
