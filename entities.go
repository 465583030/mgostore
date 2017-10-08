package mgostore

import "time"

/*
This configuration struct determines the connection to the mongo DB.
*/
type MongoConfig struct {
	// Define a comma separated array of Servers here
	Servers string
	// Define the DB Name for which the corresponding model would connect to
	DBName  string
	Timeout time.Duration
	// Specifies if the connection is SSL. This is important to use a tls.Dial function in that case due to limitations of mgo package
	IsSSL bool
}

/*
A model interface represents a storable entity
Two methods need to be defined for this.
*/
type Model interface {
	// This method should return the collection name in the mongo DB where this model will be stored
	CollectionName() string

	// Should return a valid MongoConfig which contains the configuration values to connect to the DB
	DBConfig() *MongoConfig
}

/*
A models interface represents a storable entity list of entities
Two methods need to be defined for this.
*/
type Models interface {
	// This method should return the collection name in the mongo DB where this model will be stored
	CollectionName() string

	// Should return a valid MongoConfig which contains the configuration values to connect to the DB
	DBConfig() *MongoConfig
}
