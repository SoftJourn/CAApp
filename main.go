package main
import (
	"github.com/SoftJourn/CAApp/web"
	"github.com/SoftJourn/CAApp/web/controllers"
	"github.com/SoftJourn/CAApp/src/services/faceService"
	"github.com/SoftJourn/CAApp/src/services/userService"
	"os"
	"fmt"
	"encoding/json"
	"github.com/SoftJourn/CAApp/src/types"
)

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := types.Configuration{}
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
		Configuration: configuration,
	}
	web.Serve(app)
}