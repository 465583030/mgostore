package mgostore

import (
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Add the common methods for tests here

func setTestEnvVars() {
	// override .env values to always point to the localhost instance
	os.Setenv("MONGODB_SERVERS", "localhost")
	os.Setenv("MONGODB_NAME", "test")
	os.Setenv("MONGODB_TIMEOUT", "30000")
}

func testMongoCollection() *mgo.Collection {
	s, _ := newSession(testMongoConfig())
	m := &mockModel{}
	return fetchCollection(m, s)
}

type mockModel struct {
	Id              bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	EncryptedField1 string        `json:"encrypted_field1" bson:"encrypted_field1" encrypt:"aes"`
	PlainTextField  string        `json:"plain_text_field" bson:"plain_text_field"`
	NumField        int           `json:"num_field" bson:"num_field" encrypt:"aes"`
	EncryptedField2 string        `json:"encrypted_field2" bson:"encrypted_field2" encrypt:"invalid_type"`
}

type mockModels []mockModel

func (m *mockModel) CollectionName() string {
	return "mock_models"
}

func (m mockModels) CollectionName() string {
	return "mock_models"
}

func (m *mockModel) DBConfig() *MongoConfig {
	return testMongoConfig()
}

func (m mockModels) DBConfig() *MongoConfig {
	return testMongoConfig()
}

var testMongoConfig = func() *MongoConfig {

	return &MongoConfig{Servers: os.Getenv("MONGODB_SERVERS"),
		DBName: "wellhive_test",
		Timeout: func() time.Duration {
			intTimeout := 100
			return time.Duration(intTimeout) * time.Millisecond
		}(),
	}
}

var testServers = func() string {
	return os.Getenv("MONGODB_SERVERS")
}
