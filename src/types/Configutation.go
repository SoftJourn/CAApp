package types

type Configuration struct {

	FaceServiceBaseUrl string `json:"faceServiceBaseUrl"`
	FaceServiceAppKey string `json:"faceServiceAppKey"`

	Organization Organization `json:"organization"`
}

