package main

import (
	"os"

	app "github.com/MomentsFromEarth/api/internal"
)

func main() {
	os.Setenv("API_KEY", "palebluedot")
	appEngine := app.Init()
	appEngine.Run()
}
