# cfmgo
A MongoDB integration package for Cloud Foundry

## Overview

`cfmgo` is a package to assist you in connecting Go applications running on Cloud Foundry to MongoDB.  

## Usage

`go get github.com/pivotal-pez/cfmgo`

```go
appEnv, _ := cfenv,Current() //relies on github.com/cloudfoundry-community/go-cfenv
serviceName := os.Getenv("DB_NAME")
serviceURIName := os.Getenv("DB_URI")
serviceURI := cfmgo.GetServiceBinding(serviceName, serviceURIName, appEnv)
collection := cfmgo.Connect(cfmgo.NewCollectionDialer, serviceURI, "my-collection")
```

### cfmgo/params

`params` will extract query parameters from the query string of a request into a `RequestParams` object.  RequestParams satisfies the `Params` interface defined in the `cfmgo` base package and used by the `cfmgo.Collection.Find()` method.

`RequestParams` object provides 4 methods to yield components of a MongoDB query and satisfy the `cfmgo.Params` interface:
1. **Selector()** returns a bson.M object that is used to filter the records returned by a query
1. **Scope()** returns a bson.M object that is use to filter the fields returned by a query
1. **Limit()** returns an integer that constrains the number of records returned in a result set
1. **Offset()** returns an integer that is used to skip over a number of matching records; used for paging

`params.Extract()` will interrogate the `url.Values` object of an HTTP request and scan for the following:
* `scope` -- will build a properly formatted bson.M object off of a provided set of comma-delimited fields to be used as the Select() argument in a MongoDB query.  If not provided, an empty bson.M object will be provided, which results in all fields being returned in the result set.
* `limit` -- will convert the value to an integer; if not provided, it will default to 10.  
* `offset` -- will convert the value to an integer; if not provided, it will default to 0.
* all other `name=value` pairs are assumed to represent the query selector and will be converted into a bson.M object and passed as the argument to Find()

#### Example
```go
func ListInventoryItemsHandler(collection cfmgo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		collection.Wake()

		params := params.Extract(req.URL.Query())

		items := make([]RedactedInventoryItem, 0)

		if count, err := collection.Find(params, &items); err == nil {
			Formatter().JSON(w, http.StatusOK, wrapper.Collection(&items, count))
		} else {
			Formatter().JSON(w, http.StatusNotFound, wrapper.Error(err.Error()))
		}
	}
}
```


### cfmgo/wrapper

`wrapper` is a simple helper to wrap API response data and errors in a consistent structure.  

```go
//ResponseWrapper provides a standard structure for API responses.
type ResponseWrapper struct {
	//Status indicates the result of a request as "success" or "error"
	Status string `json:"status"`
	//Data holds the payload of the response
	Data interface{} `json:"data,omitempty"`
	//Message contains the nature of an error
	Message string `json:"message,omitempty"`
	//Count contains the number of records in the result set
	Count int `json:"count,omitempty"`
}
```
#### Examples

`wrapper.Error(err)` yields:

```json
{
"status": "error",
"message": "error message text"
}
```

`wrapper.One(&someRecord)` yields:

```json
{
"status": "success",
"data": {
	"id": 1,
	"name": "fluffy"
	}
}
```

`wrapper.Collection(&someResults, count)` yields:

```json
{
"status": "success",
"data": [
	{
	"id": 1,
	"name": "fluffy"
	},
	{
	"id": 2,
	"name": "thiggy"
	}
	],
"count": 2
}
```

Note: `count` represents the total number of matching records from the query, not the number of records returned in the result set.  That number is governed by `limit`.