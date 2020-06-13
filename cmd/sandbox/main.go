package main

import (
	"fmt"

	api "github.com/MomentsFromEarth/api/internal"
)

func main() {
	fmt.Println("sandbox.main")
	apiEngine := api.Init()
	apiEngine.Run()
}
