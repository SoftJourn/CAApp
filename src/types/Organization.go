package types

import "time"

type Organization struct {
	OrgId string `json:"orgId"`
	ChaincodeApiUrl string `json:"chaincodeApiUrl"`
	UsersChaincodeName string `json:"usersChaincodeName"`
	CaCertificatePath string `json:"caCertificatePath"`
	CaKeyPath string `json:"caKeyPath"`
	CaKvsPath string `json:"caKvsPath"`
	Peers []string `json:"peers"`
	ChannelName string `json:"channelName"`
	HttpClientTimeoutSec time.Duration `json:"httpClientTimeoutSec"`
}