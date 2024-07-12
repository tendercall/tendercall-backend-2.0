package twilio

import "github.com/gin-gonic/gin"

type Config struct {
	Router *gin.Engine
}

// Routes method declaration

func (app *Config) Routes() {
	app.Router.POST("/otp", app.sendSMS())
	app.Router.POST("/verifyOTP", app.verifySMS())
}
