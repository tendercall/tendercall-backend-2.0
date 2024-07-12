package twilio

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"tendercall.com/main/data"
)

// Constant Declaration

const appTimeout = time.Second * 10

// sendSMS method declaration
// sendSMS is a method that sends SMS to the user

func (app *Config) sendSMS() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Initializes a context (context.Background()) with a timeout (appTimeout)

		_, cancel := context.WithTimeout(context.Background(), appTimeout)
		var payload data.OTPData
		defer cancel()

		/*
			app.validateBody(c, &payload) validates and parses the request body (c) into the payload structure.
			This function ensures that the required fields are present and correctly formatted.
		*/

		app.validateBody(c, &payload)

		/*
		  Prepares the necessary data for sending an OTP SMS.
		  The data includes the user's phone number
		*/

		newData := data.OTPData{
			PhoneNumber: payload.PhoneNumber,
		}

		// Sends the SMS to the user

		_, err := app.twilioSendOTP(newData.PhoneNumber)
		if err != nil {
			app.errorJSON(c, err)
			return
		}

		app.writeJSON(c, http.StatusAccepted, "OTP sent successfully")
	}
}

// verifySMS function declaration
// verifySMS is a method that verifies the OTP sent to the user

func (app *Config) verifySMS() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Initializes a context (context.Background()) with a timeout (appTimeout)

		_, cancel := context.WithTimeout(context.Background(), appTimeout)
		var payload data.VerifyData
		defer cancel()

		/*
			app.validateBody(c, &payload) validates and parses the request body (c) into the payload structure.
			This function ensures that the required fields are present and correctly formatted.
		*/

		app.validateBody(c, &payload)

		/*
		  Prepares data for OTP verification
		  The data includes the user's phone number and the OTP
		*/

		newData := data.VerifyData{
			User: payload.User,
			Code: payload.Code,
		}

		// OTP verification

		err := app.twilioVerifyOTP(newData.User.PhoneNumber, newData.Code)
		fmt.Println("err: ", err)
		if err != nil {
			app.errorJSON(c, err)
			return
		}

		app.writeJSON(c, http.StatusAccepted, "OTP verified successfully")
	}
}
