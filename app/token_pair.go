package app

type TokenPair struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}
