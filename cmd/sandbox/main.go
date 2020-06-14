package main

import (
	app "github.com/MomentsFromEarth/api/internal"
)

func main() {
	appEngine := app.Init()
	appEngine.Run()
}
