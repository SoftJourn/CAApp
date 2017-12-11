package main
import (
	"os"
	"fmt"
	"encoding/json"
	"CAApp/ca"
)

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := ca.Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	app := ca.Application{
		Configuration: configuration,
	}
	ca.Serve(app)
}