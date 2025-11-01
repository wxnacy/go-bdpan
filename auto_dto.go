package bdpan

type GetDeviceTokenRes struct {
	ExpiresIn     int32  `json:"expires_in,omitempty"`
	RefreshToken  string `json:"refresh_token,omitempty"`
	AccessToken   string `json:"access_token,omitempty"`
	SessionSecret string `json:"session_secret,omitempty"`
	SessionKey    string `json:"session_key,omitempty"`
	Scope         string `json:"scope,omitempty"`
}

type GetDeviceCodeRes struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	QrcodeURL       string `json:"qrcode_url"`
	ExpiresIn       int64  `json:"expires_in"`
	Interval        int64  `json:"interval"`
}
