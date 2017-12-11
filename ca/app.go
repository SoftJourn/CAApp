package ca

import (
	"net/http"
	"fmt"
)

func Serve(app Application) {
	http.HandleFunc("/api/certificate", app.CertificateHandler)

	fmt.Println("Listening (http://0.0.0.0:3000/) ...")
	http.ListenAndServe("0.0.0.0:3000", nil)
}