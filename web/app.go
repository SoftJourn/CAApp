package web

import (
	"net/http"
	"fmt"
	"github.com/SoftJourn/CAApp/web/controllers"
	"github.com/SoftJourn/CAApp/src/storage"
	"github.com/gorilla/sessions"
)

func Serve(app controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/login", app.LoginHandler)
	http.HandleFunc("/logout", app.LogoutHandler)
	http.HandleFunc("/generate", app.GenerateHandler)
	http.HandleFunc("/deploy", app.GenerateHandler)
	http.HandleFunc("/api/register", app.RegisterHandler)
	http.HandleFunc("/api/login", app.FaceLoginHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	})

	var storage = storage.GetInstance()
	storage.Store = sessions.NewCookieStore([]byte("something-very-secret"))

	fmt.Println("Listening (http://0.0.0.0:3000/) ...")
	http.ListenAndServe("0.0.0.0:3000", nil)
}