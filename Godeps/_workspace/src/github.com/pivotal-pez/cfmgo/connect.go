package cfmgo

import (
	"fmt"

	"github.com/cloudfoundry-community/go-cfenv"
	"gopkg.in/mgo.v2"
)

//GetServiceBinding parses cfenv.App object and returns a URI for the specified service.
func GetServiceBinding(serviceName string, serviceURIName string, appEnv *cfenv.App) (serviceURI string) {

	if service, err := appEnv.Services.WithName(serviceName); err == nil {
		if serviceURI = service.Credentials[serviceURIName].(string); serviceURI == "" {
			panic(fmt.Sprintf("we pulled an empty connection string %s from %v - %v", serviceURI, service, service.Credentials))
		}

	} else {
		panic(fmt.Sprint("Experienced an error trying to grab service binding information:", err.Error()))
	}
	return
}

//SetupDB connects to the specified database and returns a Collection object for the specified collection.
func SetupDB(dialer CollectionDialer, URI string, collectionName string) (collection Collection) {
	var (
		err      error
		dialInfo *mgo.DialInfo
	)

	if dialInfo, err = mgo.ParseURL(URI); err != nil || dialInfo.Database == "" {
		panic(fmt.Sprintf("cannot parse given URI %s due to error: %s", URI, err.Error()))
	}

	if collection, err = dialer(URI, dialInfo.Database, collectionName); err != nil {
		panic(fmt.Sprintf("cannot dial connection due to error: %s URI:%s col:%s db:%s", err.Error(), URI, collectionName, dialInfo.Database))
	}
	return
}
