package userService

import (
	"github.com/SoftJourn/CAApp/src/types"
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"time"
	"io/ioutil"
	"github.com/SoftJourn/CAApp/src/ca"
)

//TODO move out


var caCertificatePath string = "/go/src/github.com/SoftJourn/CAApp/fixtures/channel/crypto-config/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem"
var caKeyPath string = "/go/src/github.com/SoftJourn//CAApp/fixtures/channel/crypto-config/peerOrganizations/org1.example.com/ca/8791d1363e89515f9afa042b0693a2c704bb8dd95d28f97d3549a2b9e3c4352d_sk"

//var caCertificatePath string = "/home/vitaliy/projects/gocode/src/CAApp/fixtures/channel/crypto-config/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem"
//var caKeyPath string = "/home/vitaliy/projects/gocode/src/github.com/SoftJourn/CAApp/fixtures/channel/crypto-config/peerOrganizations/org1.example.com/ca/8791d1363e89515f9afa042b0693a2c704bb8dd95d28f97d3549a2b9e3c4352d_sk"
var caKvsPath string = "/tmp/fabric-client-kvs_peerOrg1"
var peers []string = []string {"peer0.org1.example.com:7051"}
var channelName = "mychannel"
var chaincodeName string = "usr"

var netClient = &http.Client{
	Timeout: time.Second * 20,
}

type UserService struct {
	ChaincodeApiUrl string
	ChaincodeName string
	OrgId string
}

type EnrollRequest struct {
	Username string `json:"username"`
	OrgName string `json:"orgName"`
}

type EnrollResponse struct {
	Success bool `json:"success"`
	Secret string `json:"secret"`
	Message string `json:"message"`
	Token string `json:"token"`
}

type InvokeRequest struct {
	Peers []string `json:"peers"`
	Fcn string `json:"fcn"`
	Args []string `json:"args"`
}

type InvokeResponse struct {
	StatusCode int    `json:"statusCode"`
	BodyBytes  []byte `json:"bodyBytes"`
}

func NewUserService(chaincodeApiUrl string, chaincodeName string, orgId string) *UserService {
	return &UserService{
		ChaincodeApiUrl: chaincodeApiUrl,
		ChaincodeName: chaincodeName,
		OrgId: orgId,
	}
}

func (us *UserService) GetUserById(userId string) (types.UserData, error) {

	var invokeData types.InvokeData
	var userData types.UserData

	enrollData, err := us.enrollUser(userId, us.OrgId)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return userData, err
	}
	fmt.Printf("enrollData: %s\n", enrollData)

	response, err := us.invokeChaincode(chaincodeName, enrollData.Token, peers, "getUserDataById", []string{userId})
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return userData, err
	}
	fmt.Printf("GetUserById response.BodyBytes: %s\n", response.BodyBytes)

	err = json.Unmarshal(response.BodyBytes, &invokeData)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return userData, err
	}

	err = json.Unmarshal([]byte(invokeData.Payload), &userData)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return userData, err
	}

	fmt.Printf("GetUserById userData: %v\n", userData)

	return userData, err
}

func (us *UserService) AddUser(userData types.UserData) error {

	caCertInfo, err := ca.Generate(userData.Email, caCertificatePath, caKeyPath)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}
	err = ca.Deploy(userData.Email, caCertInfo, caKvsPath)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}

	enrollData, err := us.enrollUser(userData.Email, us.OrgId)
	fmt.Printf("enrollData: %s\n", enrollData)

	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}

	userDataBytes, err := json.Marshal(userData)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}

	_, err = us.invokeChaincode(chaincodeName, enrollData.Token, peers, "addUser", []string{string(userDataBytes)})
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}
	return nil
}

func (us *UserService) enrollUser(name string, orgId string) (EnrollResponse, error) {

	body := EnrollRequest {
		Username: name,
		OrgName:  orgId,
	}

	var enrollResponse EnrollResponse

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return  enrollResponse, err
	}

	request, _ := http.NewRequest(http.MethodPost, us.ChaincodeApiUrl + "users", bytes.NewReader(bodyBytes))
	request.Header.Add("Content-Type", "application/json")

	response, err := netClient.Do(request)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return  enrollResponse, err
	}
	defer response.Body.Close()

	responseBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return  enrollResponse, err
	}

	err = json.Unmarshal(responseBodyBytes, &enrollResponse)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return  enrollResponse, err
	}
	fmt.Printf("enrollResponse: %v\n", enrollResponse)

	return  enrollResponse, err
}

func (us *UserService) invokeChaincode(chaincodeName string, token string, peers []string, fcnName string, args []string) (InvokeResponse, error) {

	body := InvokeRequest{
		Peers: peers,
		Fcn: fcnName,
		Args: args,
	}
	var invokeResponse InvokeResponse

	bodyBytes, err := json.Marshal(body)
	fmt.Printf("bodyBytes: %s\n", string(bodyBytes))
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return  invokeResponse, err
	}

	url := us.ChaincodeApiUrl + "channels/"+ channelName + "/chaincodes/" + chaincodeName
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyBytes))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer " + token)

	response, err := netClient.Do(request)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return  invokeResponse, err
	}
	defer response.Body.Close()

	invokeResponse.StatusCode = response.StatusCode
		invokeResponse.BodyBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return  invokeResponse, err
	}
	fmt.Printf("invoke responseBodyBytes: %s\n", invokeResponse.BodyBytes)
	return  invokeResponse, err
}