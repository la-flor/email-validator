package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/la-flor/email-validator/email"
)

type Email struct {
	Email	string	`json: email`
}

func validateEmail (c *gin.Context) {
	var emailInput Email
	
	if err := c.BindJSON(&emailInput); err != nil {
		log.Println("unable to bind request body email value", emailInput)
		return
	}

	invalidated, message := email.CheckIfInvalid(emailInput.Email)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "invalidated": invalidated, "message": message})
}


func main() {
	route := gin.Default()
	route.POST("/validate", validateEmail)
	
	route.Run(":8080")
}
