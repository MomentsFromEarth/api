package main

import (
	"os"

	app "github.com/MomentsFromEarth/api/internal"
)

func main() {
	os.Setenv("API_KEY", "<api_key>")
	os.Setenv("AWS_ACCESS_KEY_ID", "<access_key_id>")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "<secret_access_key")
	appEngine := app.Init()
	appEngine.Run()
}
