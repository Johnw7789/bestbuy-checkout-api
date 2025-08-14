package login

func (l *LoginClient) Login() (bool, string, error) {
	l.UpdateStatus("Getting login data")

	// * Get the data needed for the authentication request
	ld, err := l.getLoginSensorData()
	if err != nil {
		return false, "", err
	}

	l.UpdateStatus("Getting encryption data")

	// * Get the fingerprint which contains encrypted email, useragent, xgrid, xgridb
	encFingerprint, err := l.getEncData()
	if err != nil {
		return false, "", err
	}

	// * Attempt to login
	success, _, err := l.getBestbuyLogin(ld, encFingerprint)

	if err.Error() == ErrStepUpRequired {
		return l.handleMFA(ld, encFingerprint)
	}

	return success, "", err
}
