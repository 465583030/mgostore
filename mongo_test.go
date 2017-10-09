package mgostore

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newMongoSession(t *testing.T) {
	setTestEnvVars()

	t.Log("When the server is not reachable")
	os.Setenv("MONGODB_SERVERS", "invalid_server")
	os.Setenv("MONGODB_TIMEOUT", "30000")
	_, err := newSession(testMongoConfig())
	assert.Equal(t,
		"no reachable servers",
		err.Error(),
		"Expected not reachable servers error")

	t.Log("When the server is reachable")
	setTestEnvVars()
	s, err := newSession(testMongoConfig())
	assert.Equal(t, nil, err, "Expected no error")
	assert.NotEqual(t, nil, s, "Session returned is not empty")
}

func Test_fetchMongoCollection(t *testing.T) {
	setTestEnvVars()
	s, _ := newSession(testMongoConfig())
	m := &mockModel{}
	c := fetchCollection(m, s)
	assert.Equal(t, m.CollectionName(), c.Name, "Expected collection name to match")
}
