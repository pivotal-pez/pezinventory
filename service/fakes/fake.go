package fakes

import (
	"encoding/json"
	"strconv"

	"github.com/pivotal-pez/cfmgo"

	"gopkg.in/mgo.v2"
)

//FakeNewCollectionDialer -
func FakeNewCollectionDialer(c interface{}) func(url, dbname, collectionname string) (col cfmgo.Collection, err error) {
	b, err := json.Marshal(c)
	if err != nil {
		panic("shit is broken")
	}

	return func(url, dbname, collectionname string) (col cfmgo.Collection, err error) {
		col = &FakeCollection{
			Data: b,
		}
		return
	}
}

//FakeCollection -
type FakeCollection struct {
	mgo.Collection
	Data  []byte
	Error error
}

//Close -
func (s *FakeCollection) Close() {

}

//Wake -
func (s *FakeCollection) Wake() {

}

//Find -- finds all records matching given selector
func (s *FakeCollection) Find(params cfmgo.Params, result interface{}) (count int, err error) {
	count = 2
	err = json.Unmarshal(s.Data, result)
	return
}

//FindAndModify -
func (s *FakeCollection) FindAndModify(selector interface{}, update interface{}, result interface{}) (info *mgo.ChangeInfo, err error) {
	return
}

//UpsertID -
func (s *FakeCollection) UpsertID(id interface{}, result interface{}) (changInfo *mgo.ChangeInfo, err error) {
	return
}

//FindOne -
func (s *FakeCollection) FindOne(id string, result interface{}) (err error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	var col []interface{}
	err = json.Unmarshal(s.Data, &col)
	if err != nil {
		return
	}
	b, err := json.Marshal(col[i])
	if err != nil {
		return
	}
	err = json.Unmarshal(b, result)
	return
}
