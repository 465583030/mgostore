package mgostore

import (
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
)

/*
These are functional tests and requite that you have an instance of
mongodb locally with the db test
golang's tests run linearly and since its a compiled language, have longer tests with each testing
more functionalities.
*/

func TestCreate(t *testing.T) {
	m := &mockModel{}
	t.Log("When no connection can be established")
	setTestEnvVars()
	os.Setenv("MONGODB_SERVERS", "invalid_server")
	err := Create(m)
	assert.Equal(t,
		"no reachable servers",
		err.Error(),
		"Expected not reachable servers error")

	setTestEnvVars()
	t.Log("When connection can be established")
	m.EncryptedField1 = "crypto text"
	m.PlainTextField = "plain text"
	err = Create(m)
	t.Log(err)

	tc := testMongoCollection()
	// Make sure to drop the entire collection after the test is run
	defer tc.DropCollection()

	assert.Equal(t, nil, err, "Expected no error")
	assert.NotEqual(t, nil, m.ID, "Expected object Id to be generated")
	assert.Equal(t, "crypto text", m.EncryptedField1, "Expected encrypted field to be plain text")
	assert.Equal(t, "plain text", m.PlainTextField, "Expected plain text field")
	err = tc.FindId(m.ID).One(m)
	assert.Equal(t, nil, err, "Expected record to be saved")
	//assert.NotEqual(t, "crypto text", m.EncryptedField1, "Expected field to be encrypted in the DB")
}

func TestDelete(t *testing.T) {
	id := bson.NewObjectId()
	m := &mockModel{ID: id}
	t.Log("When no connection can be established")
	setTestEnvVars()
	os.Setenv("MONGODB_TIMEOUT", "30000")
	os.Setenv("MONGODB_SERVERS", "invalid_server")
	err := Destroy(m)
	assert.Equal(t,
		"no reachable servers",
		err.Error(),
		"Expected not reachable servers error")

	setTestEnvVars()
	t.Log("When connection can be established")
	tc := testMongoCollection()
	// Make sure to drop the entire collection after the test is run
	defer tc.DropCollection()

	t.Log("When record does not exist")
	err = Destroy(m)
	assert.Equal(t, string("not found"), err.Error(), "Not found error expected")

	t.Log("When record exists")

	tc.Insert(m)
	err = Destroy(m)
	assert.Nil(t, err)

	err = tc.FindId(id).One(m)
	assert.Equal(t, ErrRecordNotFound, err, "Expected not found error")
}

func TestUpdate(t *testing.T) {
	id := bson.NewObjectId()
	m := &mockModel{ID: id}
	t.Log("When no connection can be established")
	setTestEnvVars()
	os.Setenv("MONGODB_TIMEOUT", "30000")
	os.Setenv("MONGODB_SERVERS", "invalid_server")
	err := Update(m)
	assert.Equal(t,
		"no reachable servers",
		err.Error(),
		"Expected not reachable servers error")

	setTestEnvVars()
	t.Log("When connection can be established")
	tc := testMongoCollection()
	// Make sure to drop the entire collection after the test is run
	defer tc.DropCollection()

	tc.Insert(m)

	m.EncryptedField1 = "crypto text"
	m.PlainTextField = "plain text"

	err = Update(m)
	assert.Equal(t, nil, err, "Expected no error")
	assert.Equal(t,
		"crypto text",
		m.EncryptedField1,
		"Expected encrypted field to be decrypted")
	assert.Equal(t,
		"plain text",
		m.PlainTextField,
		"Expected plain text field to be not encrypted")
	tc.FindId(id).One(m)
	// assert.NotEqual(t,
	// 	"crypto text",
	// 	m.EncryptedField1,
	// 	"Expected encrypted field to be encrypted")
	assert.Equal(t,
		"plain text",
		m.PlainTextField,
		"Expected plain text field to be not encrypted")

}

func TestFind(t *testing.T) {
	id := bson.NewObjectId()
	m := &mockModel{ID: id}
	t.Log("When no connection can be established")
	setTestEnvVars()
	os.Setenv("MONGODB_TIMEOUT", "30000")
	os.Setenv("MONGODB_SERVERS", "invalid_server")
	err := Find(m)
	assert.Equal(t,
		"no reachable servers",
		err.Error(),
		"Expected not reachable servers error")
	setTestEnvVars()
	t.Log("When connection can be established")
	tc := testMongoCollection()
	// Make sure to drop the entire collection after the test is run
	defer tc.DropCollection()
	t.Log("When record does not exist")
	err = Find(m)
	assert.Equal(t, ErrRecordNotFound, err, "Expected not found error")

	t.Log("When record exists")
	err = tc.Insert(m)
	mNew := &mockModel{ID: id}
	err = Find(mNew)
	assert.Nil(t, err)
	assert.Equal(t, id, mNew.ID, "Invalid record returned")
}

func TestFindBy(t *testing.T) {
	m := &mockModel{}
	t.Log("When no connection can be established")
	setTestEnvVars()
	os.Setenv("MONGODB_TIMEOUT", "30000")
	os.Setenv("MONGODB_SERVERS", "invalid_server")
	err := FindBy(bson.M{"_id": "someid"}, m)
	assert.Equal(t,
		"no reachable servers",
		err.Error(),
		"Expected not reachable servers error")
	setTestEnvVars()
	t.Log("When connection can be established")
	tc := testMongoCollection()
	// Make sure to drop the entire collection after the test is run
	defer tc.DropCollection()
	id := bson.NewObjectId()
	m.ID = id
	m.NumField = 42
	t.Log("When record does not exist")
	err = FindBy(bson.M{"num_field": 42}, m)
	assert.Equal(t, ErrRecordNotFound, err, "Expected not found error")

	t.Log("When record exists")
	tc.Insert(m)
	m = &mockModel{ID: id}
	err = FindBy(bson.M{"num_field": 42}, m)
	assert.Nil(t, err)
	assert.Equal(t, id, m.ID, "Invalid record returned")
}

func TestFindMany(t *testing.T) {
	var models mockModels
	t.Log("When no connection can be established")
	setTestEnvVars()
	os.Setenv("MONGODB_TIMEOUT", "30000")
	os.Setenv("MONGODB_SERVERS", "invalid_server")
	err := FindMany(bson.M{"id": "someid"}, &models)
	assert.Equal(t,
		"no reachable servers",
		err.Error(),
		"Expected not reachable servers error")
	setTestEnvVars()
	t.Log("When connection can be established")
	tc := testMongoCollection()
	// Make sure to drop the entire collection after the test is run
	defer tc.DropCollection()
	var mIds []bson.ObjectId
	for i := 0; i < 3; i++ {
		m := &mockModel{NumField: 42, EncryptedField1: "encrypted text"}
		generateModelID(m)
		//EncryptFields(m)
		tc.Insert(m)
		mIds = append(mIds, m.ID)
	}
	err = FindMany(bson.M{"num_field": 42}, &models)
	assert.Equal(t, nil, err, "No error expected")
	for _, m := range models {
		assert.Equal(t, 42, m.NumField, "ID field does not match")
		// assert.NotEqual(t, "encrypted text", m.EncryptedField1, "Field returned is not encrypted")
		// DecryptFields(&m)
		assert.Equal(t, "encrypted text", m.EncryptedField1, "Field is encrypted wrongly")
	}

	t.Log("When limit is passed and skip isnt")
	err = FindMany(bson.M{"num_field": 42}, &models, -1, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(models), "Expected only one match")
	assert.Equal(t, mIds[0], models[0].ID)

	t.Log("When skip is passed and limit isnt")
	err = FindMany(bson.M{"num_field": 42}, &models, 1)
	assert.Equal(t, 2, len(models), "Expected only 2 matches")
	assert.Equal(t, mIds[1], models[0].ID)
	assert.Equal(t, mIds[2], models[1].ID)
}
