package main
import (
	"github.com/SoftJourn/CAApp/web"
	"github.com/SoftJourn/CAApp/web/controllers"
	"github.com/SoftJourn/CAApp/src/services/faceService"
	"github.com/SoftJourn/CAApp/src/services/userService"
	"os"
	"fmt"
	"encoding/json"
)

type Configuration struct {

	FaceServiceBaseUrl string `json:"faceServiceBaseUrl"`
	FaceServiceAppKey string `json:"faceServiceAppKey"`

	Organization userService.Organization `json:"organization"`
}

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(configuration)

	faceService := faceService.NewFaceService(configuration.FaceServiceBaseUrl, configuration.FaceServiceAppKey)

	userService := userService.NewUserService(configuration.Organization)

	app := controllers.Application{
		FaceService: *faceService,
		UserService: *userService,
	}
	web.Serve(app)
}