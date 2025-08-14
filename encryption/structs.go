package encryption

type Fingerprint struct {
	EncryptedEmail    string
	EncryptedPassword string
	EncryptedInfo     string
	EncryptedActivity string
	EncryptedXGrid    string
	XGridB            string
}

type EncParams struct {
	Email     string
	Password  string
	UserAgent string
	EncKeys   Keys
}

type Keys struct {
	EmailKey       Key
	CredentialsKey Key
	InfoKey        Key
	ActivityKey    Key
	XGridKey       Key
}

type Key struct {
	KeyId  string
	PubKey string
}

type XGrid struct {
	BP          string `json:"bP"`
	CH          string `json:"cH"`
	WH          string `json:"wH"`
	Platform    string `json:"p"`
	NavigatorOS string `json:"os"`
	ColorDepth  int    `json:"cD"`
	Concurrency int    `json:"nC"`
	TouchScreen bool   `json:"tS"`
}

type XGridB struct {
	GCV string `json:"gCV"`
	GCN string `json:"gCN"`
	AB  string `json:"aB"`
	SR  string `json:"sR"`
	SL  string `json:"sL"`
	SF  string `json:"sF"`
	SFC int    `json:"sFC"`
	ST  string `json:"sT"`
}
