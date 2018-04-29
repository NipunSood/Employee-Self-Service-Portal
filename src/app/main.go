package main

import (
	"github.com/NipunSood/Employee-Self-Service-Portal/src/controller"
)

func main() {

	r := controller.RegisterRouters()
	r.Run(":3000")
}
