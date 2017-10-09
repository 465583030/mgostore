package mgostore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestFetchModelIdVal(t *testing.T) {
	id := bson.NewObjectId()
	m := &mockModel{ID: id}
	assert.Equal(t, id, fetchModelIDVal(m), "Returned value does not match the model id")
}

func TestGenerateModelId(t *testing.T) {
	m := &mockModel{}
	generateModelID(m)
	assert.NotNil(t, m.ID)
}
