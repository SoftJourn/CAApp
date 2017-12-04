package controllers

import (
	"net/http"
	"encoding/json"
	"github.com/SoftJourn/CAApp/src/ldap"
	"fmt"
	"github.com/SoftJourn/CAApp/src/ca"
)

type CertificateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CertificateResponse struct {
	Certificate string `json:"certificate"`
	PublicKey string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func (app *Application) CertificateHandler(responseWriter http.ResponseWriter, request *http.Request) {

	var certificateRequest CertificateRequest

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&certificateRequest)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusExpectationFailed)
		return
	}

	if len(certificateRequest.Username) == 0 {
		http.Error(responseWriter, "Username is required", http.StatusExpectationFailed)
		return
	}

	if len(certificateRequest.Username) == 0 {
		http.Error(responseWriter, "Password is required", http.StatusExpectationFailed)
		return
	}

	ldapUser, _, err := ldap.GetUser(certificateRequest.Username, certificateRequest.Password)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusExpectationFailed)
		return
	}

	fmt.Printf("LDAP User: %v", ldapUser)

	certificateInfo, err := ca.Generate(ldapUser.Email, app.Configuration.Organization.CaCertificatePath, app.Configuration.Organization.CaKeyPath)
	if err != nil {
		fmt.Errorf("failed to generate certificate: %s", err)
	}

	fmt.Printf("certificateInfo: %v", certificateInfo)

	certificateResponse := CertificateResponse{
		Certificate: certificateInfo.Certificate,
		PublicKey: certificateInfo.PublicKey,
		PrivateKey: certificateInfo.PrivateKey,
	}

	jsonCertificateResponse, err := json.Marshal(certificateResponse)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusExpectationFailed)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(jsonCertificateResponse)
	return
}