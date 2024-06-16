package main

import (
	"golang-dev-logic-challenge/routes"
)

// @Title Golang Dev Logic Challenge
// @Version 1.0

func main() {
	router := routes.SetupRouter()
	router.Run(":8080")
}
