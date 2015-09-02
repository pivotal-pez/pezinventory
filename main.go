package main

import (
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	pez "github.com/pivotal-pez/pezinventory/service"
)

func main() {
	appEnv, _ := cfenv.Current()
	// validatorServiceName := os.Getenv("UPS_PEZVALIDATOR_NAME")
	// targetKeyName := os.Getenv("UPS_PEZVALIDATOR_TARGET")
	// service, _ := appEnv.Services.WithName(validatorServiceName)
	// validationTargetUrl := service.Credentials[targetKeyName]
	// s := pez.NewServer(keycheck.NewNegroniAPIKeyCheckMiddleware(validationTargetUrl.(string)))

	s := pez.NewServer(appEnv)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	s.Run(":" + port)
}
