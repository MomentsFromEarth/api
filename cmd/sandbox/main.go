package main

import (
	"os"

	app "github.com/MomentsFromEarth/api/internal"
)

func main() {
	os.Setenv("API_KEY", "palebluedot")
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIA3JY5AC66I34TJBQS")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "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")
	appEngine := app.Init()
	appEngine.Run()
}
