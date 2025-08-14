package encryption

import "encoding/json"

// * Generates an XGrid json payload and returns it as a string to be encrypted
func GenerateXGrid() (string, error) {
	// todo: randomization
	xgrid := XGrid{
		BP:          "86dd3e6ea46118c6c572252debbca7ef516f2ffd",
		CH:          "f257132fefc2c6000a936c2c93d4d8c2e7c88f13",
		WH:          "609ff3db47b19ce8cd1d1573f564a4bd93343714",
		Platform:    "Win32",   // navigator.platform
		NavigatorOS: "Windows", // navigator OS
		ColorDepth:  24,
		Concurrency: 16,
		TouchScreen: false,
	}

	jBytes, err := json.Marshal(xgrid)
	if err != nil {
		return "", err
	}

	return string(jBytes), nil
}

// * Generates an XGridB json payload and returns it as a string
func GenerateXGridB() (string, error) {
	// todo: randomization
	// xGridB := XGridB{
	// 	GCV: "Google Inc. (NVIDIA)",
	// 	GCN: "ANGLE (NVIDIA, NVIDIA GeForce RTX 3080 (0x00002206) Direct3D11 vs_5_0 ps_5_0, D3D11)",
	// 	AB:  "101.45646868248605",
	// 	SR:  "2560, 1440",
	// 	SL:  "en-US",
	// 	SF:  "27ad0b308d5a51f5345a04ceea3df69b494492ea",
	// 	SFC: 57,
	// 	ST:  "GMT-0500 (Eastern Standard Time)",
	// }

	// {"gCV":"Google Inc. (NVIDIA)","gCN":"ANGLE (NVIDIA, NVIDIA GeForce RTX 5070 (0x00002F04) Direct3D11 vs_5_0 ps_5_0, D3D11)","aB":"101.45646868248605","sR":"2560, 1440","sL":"en-US","sF":"27ad0b308d5a51f5345a04ceea3df69b494492ea","sFC":57,"sT":"GMT-0400 (Eastern Daylight Time)"}
	xGridB := XGridB{
		GCV: "Google Inc. (NVIDIA)",
		GCN: "ANGLE (NVIDIA, NVIDIA GeForce RTX 5070 (0x00002F04) Direct3D11 vs_5_0 ps_5_0, D3D11)",
		AB:  "101.45646868248605",
		SR:  "2560, 1440",
		SL:  "en-US",
		SF:  "27ad0b308d5a51f5345a04ceea3df69b494492ea",
		SFC: 57,
		ST:  "GMT-0400 (Eastern Daylight Time)",
	}

	jBytes, err := json.Marshal(xGridB)
	if err != nil {
		return "", err
	}

	return string(jBytes), nil
}
