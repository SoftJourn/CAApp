package ca

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"errors"
	"strings"
)

type CertificateRequest struct {
	Ldap string `json:"ldap"`
}

type CertificateResponse struct {
	Certificate string `json:"certificate"`
	PublicKey   string `json:"publicKey"`
	PrivateKey  string `json:"privateKey"`
}

type Application struct {
	Configuration Configuration
}

func (app *Application) CertificateHandler(responseWriter http.ResponseWriter, request *http.Request) {

	var certificateRequest CertificateRequest

	verErr := verifyToken(request.Header.Get("Authorization"), app.Configuration.AuthRSAPublicKey)

	if verErr != nil {
		http.Error(responseWriter, verErr.Error(), http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&certificateRequest)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusExpectationFailed)
		return
	}

	if len(certificateRequest.Ldap) == 0 {
		http.Error(responseWriter, "Ldap name is required", http.StatusExpectationFailed)
		return
	}

	certificateInfo, err := Generate(certificateRequest.Ldap, app.Configuration.Organization.CaCertificatePath, app.Configuration.Organization.CaKeyPath)
	if err != nil {
		fmt.Errorf("failed to generate certificate: %s", err)
	}

	fmt.Printf("certificateInfo: %v", certificateInfo)

	certificateResponse := CertificateResponse{
		Certificate: certificateInfo.Certificate,
		PublicKey:   certificateInfo.PublicKey,
		PrivateKey:  certificateInfo.PrivateKey,
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

func verifyToken(tokenString string, rsaKey string) error {
	if len(tokenString) > 0 {

		if !strings.Contains(tokenString, "Bearer ") {
			return errors.New("JWT token format is wrong")
		}

		runes := []rune(tokenString)
		tokenString := string(runes[7:len(tokenString)])

		_, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			data, err := ioutil.ReadFile(rsaKey)
			if err != nil {
				return nil, err
			}
			rsa, err := jwt.ParseRSAPublicKeyFromPEM(data)
			if err != nil {
				return nil, err
			}
			return rsa, nil
		})
		return err
	} else {
		return errors.New("Authorization header was not found")
	}
}