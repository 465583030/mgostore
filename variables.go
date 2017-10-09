package mgostore

import (
	"errors"
	"sync"

	mgo "gopkg.in/mgo.v2"
)

var (
	sessionMap  = make(map[string]*mgo.Session)
	sessionWMux sync.Mutex
	sessionRMux sync.RWMutex
)

// All error variables here

var ErrInvalidId = errors.New("invalid id")
var ErrRecordNotFound = mgo.ErrNotFound
var ErrMongoCollectionNotFetched = errors.New("mongo collection not fetched")
var ErrMissingCryptoSecret = errors.New("missing crypto secret")
