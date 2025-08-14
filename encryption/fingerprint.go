package encryption

import "fmt"

func GenerateFingerprint(params EncParams) (Fingerprint, error) {
	encEmail, err := EncryptJoin(params.Email, params.EncKeys.EmailKey.PubKey, params.EncKeys.EmailKey.KeyId)
	if err != nil {
		return Fingerprint{}, err
	}

	activity, err := GenerateActivity(params.Email)
	if err != nil {
		return Fingerprint{}, err
	}

	encActivity, err := EncryptJoin(activity, params.EncKeys.InfoKey.PubKey, params.EncKeys.InfoKey.KeyId)
	if err != nil {
		return Fingerprint{}, err
	}

	encPass, err := EncryptSmall(params.Password, params.EncKeys.CredentialsKey.PubKey, params.EncKeys.CredentialsKey.KeyId)
	if err != nil {
		return Fingerprint{}, err
	}

	encInfo, err := EncryptJoin(fmt.Sprintf(`{"userAgent": "%s"}`, params.UserAgent), params.EncKeys.InfoKey.PubKey, params.EncKeys.InfoKey.KeyId)
	if err != nil {
		return Fingerprint{}, err
	}

	xGrid, err := GenerateXGrid()
	if err != nil {
		return Fingerprint{}, err
	}

	encXGrid, err := EncryptJoin(xGrid, params.EncKeys.XGridKey.PubKey, params.EncKeys.XGridKey.KeyId)
	if err != nil {
		return Fingerprint{}, err
	}

	xGridB, err := GenerateXGridB()
	if err != nil {
		return Fingerprint{}, err
	}

	return Fingerprint{
		EncryptedEmail:    encEmail,
		EncryptedActivity: encActivity,
		EncryptedPassword: encPass,
		EncryptedInfo:     encInfo,
		EncryptedXGrid:    encXGrid,
		XGridB:            xGridB,
	}, nil
}
