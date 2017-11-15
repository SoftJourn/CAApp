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

var organization Organization

var netClient = &http.Client{
	Timeout: time.Second * 20,
}

type UserService struct {
	ChaincodeApiUrl   string
	UserChaincodeName string
	OrgId             string
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

func NewUserService(org Organization) *UserService {
	organization = org
	netClient.Timeout = time.Second * org.HttpClientTimeoutSec

	return &UserService{
		ChaincodeApiUrl:   org.ChaincodeApiUrl,
		UserChaincodeName: org.UsersChaincodeName,
		OrgId:             org.OrgId,
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

	response, err := us.invokeChaincode(organization.UsersChaincodeName, enrollData.Token, organization.Peers, "getUserDataById", []string{userId})
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

	caCertInfo, err := ca.Generate(userData.Email, organization.CaCertificatePath, organization.CaKeyPath)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}
	err = ca.Deploy(userData.Email, caCertInfo,  organization.CaKvsPath)
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

	_, err = us.invokeChaincode(organization.UsersChaincodeName, enrollData.Token, organization.Peers, "addUser", []string{string(userDataBytes)})
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

	url := us.ChaincodeApiUrl + "channels/"+ organization.ChannelName + "/chaincodes/" + chaincodeName
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

func (us *UserService) GetOrganization() Organization {
	return organization
}