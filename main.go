package main

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"

    "fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(router)
	fmt.Printl("server run on 8000");

	router.Run("localhost:8000")
}
