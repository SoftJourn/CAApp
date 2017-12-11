package ca

type Configuration struct {

	AuthRSAPublicKey string `json:"authRSAPublicKey"`

	Organization Organization `json:"organization"`
}

