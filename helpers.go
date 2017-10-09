package mgostore

import (
	"reflect"

	"gopkg.in/mgo.v2/bson"
)

func fetchModelIDVal(m Model) interface{} {
	return reflect.ValueOf(m).Elem().FieldByName("ID").Interface()
}

func generateModelID(m Model) {
	s := reflect.ValueOf(m).Elem()
	f := s.FieldByName("ID")
	id := bson.NewObjectId()
	f.Set(reflect.ValueOf(id))
}
