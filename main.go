package main
import (
	"github.com/SoftJourn/CAApp/web"
	"github.com/SoftJourn/CAApp/web/controllers"
	"github.com/SoftJourn/CAApp/src/services/faceService"
	"github.com/SoftJourn/CAApp/src/services/userService"
)


type User struct {
	Name string
}

func main() {

	baseUrl := "https://westeurope.api.cognitive.microsoft.com/face/v1.0/"
	appKey := "ddafd2521bb64f64b79d915c70ecb35e"
	faceService := faceService.NewFaceService(baseUrl, appKey)


	//personGroupId := "1"
	//imageBytes := []byte(`{"url":"https://www.smileexpo.ru/public/upload/speakers/tn3_robert_wiecko_15063270100748_image.jpg"}`)
	//
	//faceId, _ := faceService.DetectFace(imageBytes)
	//personId, _ := faceService.CreatePerson(personGroupId, "Jim")
	//faceService.AddFaceToPerson(personGroupId, personId , imageBytes)
	//faceService.VerifyFace(faceId, personId, personGroupId)


	userService := userService.NewUserService("http://node.example.com:4000/", "usr", "org1")

	// Make the web application listening
	app := controllers.Application{
		FaceService: *faceService,
		UserService: *userService,
	}
	web.Serve(app)
}