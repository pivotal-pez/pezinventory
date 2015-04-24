package main

import (
	"os"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/pivotalservices/pezauth/keycheck"
	pez "github.com/pivotalservices/pezinventory/service"
)

func main() {
	appEnv, _ := cfenv.Current()
	validatorServiceName := os.Getenv("UPS_PEZVALIDATOR_NAME")
	targetKeyName := os.Getenv("UPS_PEZVALIDATOR_TARGET")
	service, _ := appEnv.Services.WithName(validatorServiceName)
	validationTargetUrl := service.Credentials[targetKeyName]
	s := pez.NewServer(keycheck.NewAPIKeyCheckMiddleware(validationTargetUrl).Handler())
	s.Run()
}
