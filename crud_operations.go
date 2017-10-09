package mgostore

import "gopkg.in/mgo.v2/bson"

/*
Add all storage related methods here.
Contains basic methods for CRUD operations on models.
*/

/*
Stores the struct to DB
This will update this model with all its attributes in the DB
*/
func Update(m Model) error {
	id := fetchModelIDVal(m)
	session, err := newSession(m.DBConfig())
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		return err
	}
	c := fetchCollection(m, session)
	if c == nil {
		return ErrMongoCollectionNotFetched
	}
	// values need to be encrypted
	encryptFields(m)

	c.UpdateId(id, bson.M{"$set": m})
	if err != nil {
		return err
	}
	// Fetch the saved value from storage
	Find(m)
	return nil
}

/*
Returns the Struct from the DB
For this to work, the model should be initialized with the correct value of Id for which to lookup in DB
*/
func Find(m Model) error {
	id := fetchModelIDVal(m)
	// if !bson.IsObjectIdHex(id) {
	// 	return errors.New("invalid id")
	// }
	session, err := newSession(m.DBConfig())
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		return err
	}
	c := fetchCollection(m, session)
	if c == nil {
		return ErrMongoCollectionNotFetched
	}
	if err = c.FindId(id).One(m); err != nil {
		return err
	}
	decryptFields(m)
	return nil
}

/*
Delete from the DB
For this to work, the model should be initialized with the correct value of Id for which to lookup in DB
*/
func Destroy(m Model) error {
	id := fetchModelIDVal(m)
	session, err := newSession(m.DBConfig())
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		return err
	}
	c := fetchCollection(m, session)
	if c == nil {
		return ErrMongoCollectionNotFetched
	}
	if err := c.Remove(bson.M{"_id": id}); err != nil {
		return err
	}
	return nil
}

/*
Create the model in DB
*/
func Create(m Model) error {
	session, err := newSession(m.DBConfig())
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		return err
	}
	encryptFields(m)
	generateModelID(m)
	c := fetchCollection(m, session)
	if c == nil {
		return ErrMongoCollectionNotFetched
	}
	if err = c.Insert(m); err != nil {
		return err
	}

	// Fetch stored values after saving
	Find(m)

	return nil
}

/*
FindBy returns a record of a model interface by a where clause passed to it.
It expects an input of the whereClause as the shortened bson.M format.
Check here : https://godoc.org/gopkg.in/mgo.v2/bson#M
*/
func FindBy(whereClause bson.M, m Model) error {
	session, err := newSession(m.DBConfig())
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		return err
	}
	c := fetchCollection(m, session)
	if c == nil {
		return ErrMongoCollectionNotFetched
	}
	if err = c.Find(whereClause).One(m); err != nil {
		return err
	}
	decryptFields(m)
	return nil
}

/*
FindBy returns many records of a model interface by a where clause passed to it.
It expects the string collectionName it should look into
It expects an input of the whereClause as the shortened bson.M format.
Check here : https://godoc.org/gopkg.in/mgo.v2/bson#M

This method has an issue with Decryption. So Decrypt the models manually.
var models []MyAwesomeModel
FindMany("my_awesome_models", bson.M{"some_field": "some_field_value"}, &models)

In case limit and skip is called, then please pass them in the options parameter in the order skip, limit
eg, to have no skips but to have a limit of 10

FindMany("my_awesome_models", bson.M{"some_field": "some_field_value"}, &models, -1, 10)
*/
func FindMany(whereClause bson.M, models Models, options ...int) error {
	session, err := newSession(models.DBConfig())
	if session != nil {
		defer session.Close()
	}
	if err != nil {
		return err
	}
	c := session.DB(models.DBConfig().DBName).C(models.CollectionName())
	if c == nil {
		return ErrMongoCollectionNotFetched
	}
	q := c.Find(whereClause)
	limit := -1
	skip := -1
	if len(options) > 0 {
		skip = options[0]
		if len(options) > 1 {
			limit = options[1]
		}
	}
	if skip > 0 {
		q.Skip(skip)
	}
	if limit > 0 {
		q.Limit(limit)
	}
	if err = q.All(models); err != nil {
		return err
	}
	return nil
}
