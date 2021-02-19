package main

import (
	"DStorage/service/auth/route"
)

func main() {
	r := route.Router()
	r.Run(":8080")
}
