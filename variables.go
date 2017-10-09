package mgostore

import (
	"sync"

	mgo "gopkg.in/mgo.v2"
)

var (
	sessionMap  = make(map[string]*mgo.Session)
	sessionWMux sync.Mutex
	sessionRMux sync.RWMutex
)
