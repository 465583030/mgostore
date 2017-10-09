package mgostore

import (
	"crypto/tls"
	"net"

	mgo "gopkg.in/mgo.v2"
)

/*
newSession returns a new mgo session for use with a
set of operations.  One 'master' session per config.Servers
is implicitly managed by newSession.  The returned session
is a copy of the 'master' session.
*/
func newSession(config *MongoConfig) (*mgo.Session, error) {
	servers := config.Servers

	// Two levels of locks to allow multiple concurrent reads
	// but only one write at a time
	sessionRMux.RLock()
	defer sessionRMux.RUnlock()

	// check if session already exists
	s, ok := sessionMap[servers]
	if ok {
		return s.Copy(), nil
	}

	// Write lock
	sessionWMux.Lock()
	defer sessionWMux.Unlock()

	// Check if session had been created by another goroutine
	// in between previous check above and acquiring the
	// write lodk
	s, ok = sessionMap[servers]
	if ok {
		return s.Copy(), nil
	}

	// Now create new session
	dialInfo, err := mgo.ParseURL(servers)
	if err != nil {
		return nil, err
	}

	// Establish the TLS handshake manually
	// This needs to be done till mgo package fixes this issue
	if config.IsSSL {
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		}
	}
	dialInfo.Timeout = config.Timeout

	// Dial the new connection
	s, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	// Store new session to session Map
	sessionMap[servers] = s
	return s.Copy(), nil
}

/*
Fetches the Mongo Collection of the record you are looking for
*/
func fetchCollection(m Model, session *mgo.Session) *mgo.Collection {
	dbName := m.DBConfig().DBName
	collectionName := m.CollectionName()
	return session.DB(dbName).C(collectionName)
}
