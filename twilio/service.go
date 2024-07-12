package twilio

import (
	"errors"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

func envACCOUNTSID() string {
	return "AC69a9784ab0e8b478796082e5dc31c1c3"
}

func envAUTHTOKEN() string {
	return "ef8320b60ac75cccf2e180b19aa5f4e7"
}

func envSERVICESID() string {
	return "VA29225696496d8141658520ccc6e0f81e"
}

// Initializing a client object from the Twilio SDK (twilio.RestClient)
//	using the provided credentials fetched from environment variables
// twilio.NewRestClientWithParams creates a new Twilio REST client with the specified parameters

var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: envACCOUNTSID(),
	Password: envAUTHTOKEN(),
})

// twilioSendOTP method declaration
// It sends an OTP (One-Time Password) SMS using Twilio's Verify API.

func (app *Config) twilioSendOTP(phoneNumber string) (string, error) {
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	// Initiate the OTP verification process using the Twilio service identified
	// by the SERVICE_SID environment variable

	resp, err := client.VerifyV2.CreateVerification(envSERVICESID(), params)
	if err != nil {
		return "", err
	}
	println(resp.Sid)
	return *resp.Sid, nil
}

// twilioVerifyOTP method declaration
// It verifies the OTP (One-Time Password) using Twilio's Verify API.

func (app *Config) twilioVerifyOTP(phoneNumber string, code string) error {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	// check the OTP against the Twilio service identified
	// by the SERVICE_SID environment variable

	resp, err := client.VerifyV2.CreateVerificationCheck(envSERVICESID(), params)
	if err != nil {
		return err
	}

	// BREAKING CHANGE IN THE VERIFY API
	// https://www.twilio.com/docs/verify/quickstarts/verify-totp-change-in-api-response-when-authpayload-is-incorrect
	if *resp.Status != "approved" {
		return errors.New("not a valid code")
	}

	return nil
}
