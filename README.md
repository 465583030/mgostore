# mgostore
An ORM like adaptor to store structs in mongo DB. This heavily uses the external mongo library [mgo](http://labix.org/mgo)

## Description
This package is designed in such a way that it could be used easily to have operations on models (structs which can be stored as documents in mongo).

We would create a model called `MyAwesomeModel`
Create a struct with the field. It should compulsorily have one field, viz.,
`ID` which is the `\_id` in the mongoDB storage. This is of type `bson.ObjectId`, so you will need to parse it into Hex string when you get it as a string from GET requests.

For fields which you do not want to show up in the json responses have the tag `json:"-"`
Binary json tags would be used to unmarshall this struct to the mongoDb document.
eg, field ActivatedDate with `bson:"activated_date"` will be stored in mongo with key activated_date
for this field. If not specified, it automatically lower cases it, so it will then become activateddate.

For fields which you would like to be stored encrypted, simply add the tag `encrypt="aes"`
Right now, only this option is supported for encryption

```go
type MyAwesomeModel struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	MyAwesomeField   string `json:"my_awesome_field" bson:"my_awesome_field"`
	AnEncryptedField string `json:"an_encrypted_field" bson:"an_encrypted_field" encrypt:"aes"`
}
```
The `json` tag helps in marshaling and unmarshaling from and to the json responses and requests. If not specified then it will be directly the name of the field. So,
```go
type SomeModel struct {
	SomeField string
}
=>
{
  "SomeField": ..
}
```
The `bson` tag similarly helps in marshaling and unmarshaling to mongoDB storage. If not specified, it is by default lowercased.

You need to have a method CollectionName() string on your struct.
This should simply return the name of the collection in mongoDB
```go
func (_ *MyAwesomeModel) CollectionName() string {
	return "my_awesome_models"
}
```
You also need to specify the config which the model should use to connect to the MongoDB. This gives you the flexibility to have multiple DB connections in the same application. This can be done by simply adding the method `DBConfig()`.
```go
func (_ *MyAwesomeModel) DBConfig() *mgostore.MongoConfig {
	return &mgostore.MongoConfig{
		Servers: "localhost",
		DBName: "test",
		Timeout: 500,
    		CryptoConfig: &mgostore.CryptoConfig{AESSecret: []byte("SOMEKEYSECRET")},
	}
}
```
Now you can store Your model to mongoDB by
```go
mam := &MyAwesomeModel{MyAwesomeField: "some value1", AnEncryptedField: "shhhhhh"}

// Store to DB
mgostore.Create(mam)
```

Similarly you can perform other operations

```go
mam := &MyAwesomeModel{Id: oId, MyAwesomeField: "some value1", AnEncryptedField: "shhhhhh"}

// update in mongoDB for the record with id 1234
mgostore.Update(mam)

// Find model from the DB
mam := &MyAwesomeModel{Id: oId}
mgostore.Find(mam)

// Find one record based on some where clause from the DB
whereClause := bson.M{"my_awesome_field": "some val"}
mam := &MyAwesomeModel
mgostore.FindBy(whereClause, mam)

// Find multiple records based on some where clause from the DB
type models []MyAwesomeModel
var models MyAwesomeModels
whereClause := bson.M{"my_awesome_field": "some val"}
mgostore.FindMany(whereClause, &models)

// Delete from storage
mgostore.Destroy(mam)
```
If you want nested documents then the `mgo` package used requires the tag `bson:",inline"`. Consider the following example

```go
type Person struct {
	Id            bson.ObjectId  `json:"_id" bson:"_id,omitempty"`
	Name          string  `json:"name"`
	HomeAddress   Address `json:"home_address" bson:"home_address,inline"`
	OfficeAddress Address `json:"office_address" bson:"office_address,inline"`
}

type Address struct {
	AddressLine1 string `json:"address_line_1" bson:"address_line_1"`
	AddressLine2 string `json:"address_line_2" bson:"address_line_2"`
	City         string `json:"city"`
}

```

## Testing
First make sure you have mongoDB running on your machine.
The project is maintained in [`govendor`](https://github.com/kardianos/govendor). 

Run the tests for this package by running
```
govendor test +local
```
