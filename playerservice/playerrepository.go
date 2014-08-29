package playerservice

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PlayerRepository interface {
	SaveUser(usr User) error
	FindAllUsers() []User
	FindUser(id string) (User, error)
	IsUserUnique(usr User) bool
}

type MongoRepository struct {
}

func (m MongoRepository) SaveUser(usr User) error {
	c, session := m.GetUserCollection()
	defer session.Close()

	err := c.Insert(&usr)
	return err
}

func (m MongoRepository) GetUserCollection() (*mgo.Collection, *mgo.Session) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("ranky").C("players")
	return c, session
}

func (m MongoRepository) FindAllUsers() []User {
	var results []User

	c, session := m.GetUserCollection()
	defer session.Close()
	c.Find(bson.M{}).All(&results)
	return results
}

func (m MongoRepository) isUserUnique(usr User) bool {
	var result User
	c, session := m.GetUserCollection()
	defer session.Close()

	c.Find(bson.M{"name": usr.Name}).One(&result)
	if result == (User{}) {
		return true
	} else {
		return false
	}

}
