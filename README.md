[![wercker status](https://app.wercker.com/status/918a2f54ea2bee6f8ec0c1e04c19ca79/m "wercker status")](https://app.wercker.com/project/bykey/918a2f54ea2bee6f8ec0c1e04c19ca79)

# Pez Inventory Service

The inventory capability of Pez.

### Configure for local development

Clone this repo and navigate to the project root.

```
$ cd $GOPATH/src/github.com/pivotal-pez
$ git clone https://github.com/pivotal-pez/pezinventory.git && cd pezinventory
```

Run the `init_developer_environment` script to copy the files from the samples directory into a local_wercker_configs directory in your project root.

```
$ ./init_developer_environment
```

Configure and start a local docker environment.  If you already have docker configured locally, you can safely ignore these next steps.  The following steps show how to install and configure [docker-machine](http://docs.docker.com/machine/). All subsequent instructions will assume you are using docker-machine for your local docker environment.

```
$ brew update
$ brew install docker-machine
$ docker-machine create --driver virtualbox dev
$ eval "$(docker-machine env dev)"
$ docker-machine ip dev
192.168.99.100
```

Update `local_wercker_configs/myenv` with the IP address of your docker host.

```
...

# Local DOCKER IP Address
X_DOCKER_IP=192.168.99.100
```

The environment should now be configured sufficiently to execute locally.

**To run the wercker pipeline locally:**

```
$ ./testrunner
```

This will execute the wercker pipeline and then exit.  You should see the results of unit tests, integration tests, linting, and coverage analysis.

###### Example:
``` 
--> Running step: wercker-init
--> Running step: setup-go-workspace
package-dir option not set, will use default: /gopath/src/github.com/pivotal-pez/pezinventory
$WERCKER_SOURCE_DIR now points to: /gopath/src/github.com/pivotal-pez/pezinventory
Go workspace setup finished
--> Running step: golint
service/inventory.go:45:3: can probably use "var items []RedactedInventoryItem" instead
service/inventory.go:130:50: func parameter leaseId should be leaseID
--> Running step: go get
go version go1.4 linux/amd64
--> Running step: go unit tests
?   	github.com/pivotal-pez/pezinventory	[no test files]
=== RUN TestInventory
Running Suite: Pez Inventory Suite
==================================
Random Seed: 1442329508
Will run 5 of 5 specs

•••••
Ran 5 of 5 Specs in 0.009 seconds
SUCCESS! -- 5 Passed | 0 Failed | 0 Pending | 0 Skipped --- PASS: TestInventory (0.01s)
PASS
ok  	github.com/pivotal-pez/pezinventory/service	0.045s
?   	github.com/pivotal-pez/pezinventory/service/fakes	[no test files]
?   	github.com/pivotal-pez/pezinventory/service/integrations	[no test files]
--> Running step: go integration test
=== RUN TestOsutils
Running Suite: Pez Dispenser Suite
==================================
Random Seed: 1442329511
Will run 6 of 6 specs

[negroni] Started GET /v1/inventory
[negroni] Completed 200 OK in 652.924µs
•[negroni] Started POST /v1/inventory
[negroni] Completed 200 OK in 24.057341ms
•[negroni] Started GET /v1/inventory
[negroni] Completed 200 OK in 1.057174ms
•[negroni] Started POST /v1/leases
[negroni] Completed 200 OK in 2.931449ms
•[negroni] Started GET /v1/leases/55f833a759d11100b0000002
[negroni] Completed 200 OK in 710.381µs
•[negroni] Started POST /v1/leases
[negroni] Completed 404 Not Found in 1.784894ms
•
Ran 6 of 6 Specs in 0.061 seconds
SUCCESS! -- 6 Passed | 0 Failed | 0 Pending | 0 Skipped --- PASS: TestOsutils (0.06s)
PASS
ok  	github.com/pivotal-pez/pezinventory/service/_tests	0.124s
--> Running step: go test cover
? github.com/pivotal-pez/pezinventory [no test files]
ok github.com/pivotal-pez/pezinventory/service 0.012s coverage: 22.3% of statements
22.3
Coverage - 22%
Watermark- %
? github.com/pivotal-pez/pezinventory/service/fakes [no test files]
? github.com/pivotal-pez/pezinventory/service/integrations [no test files]
--> Running step: build the server binary
--> Running step: add repo to artifact
--> Steps passed
--> Pipeline finished
```

*Note: The initial run will take some time as the docker containers will be downloaded from dockerhub.  Subsequent calls will execute more quickly.*

**To deploy the application locally**

```
$ ./runlocaldeploy local_wercker_configs/myenv
```

This will run a different wercker pipeline, governed by the files in `local_wercker_configs`.  As it is configured by default, this, too, will run the unit & integration tests, linting, and coverage analysis.  However, this can be easily changed by editing `wercker_local_deploy.yml`.

The `runlocaldeploy` command differs from `testrunner` in that it will start the server and leave it running so that you can interact with it.  It publishes the server on the docker host IP address on port 80.

```
...

--> Step passed: build the server binary
--> Running step: watch
set +e
PORT=80 ./main

[negroni] listening on :80
```

With the server now running you can interact with it using your docker host IP.

```
$ docker-machine ip dev
192.168.99.100
$ curl http://192.168.99.100/v1/inventory
{
  "status": "success",
  "data": []
}
```
