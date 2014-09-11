package matchservice

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoMatchRepository struct {
}

func (m MongoMatchRepository) SaveMatch(match Match) {

	c, session := m.MatchCollection()
	defer session.Close()

	err := c.Insert(match)
	if err != nil {
		panic(err)
	}

}

func (m MongoMatchRepository) FindMatch(id string) Match {
	result := Match{}
	c, session := m.MatchCollection()
	defer session.Close()

	c.Find(bson.M{"id": id}).One(&result)
	return result
}

func (m MongoMatchRepository) FindAllMatches() []Match {
	results := make([]Match, 0)

	c, session := m.MatchCollection()
	defer session.Close()
	c.Find(bson.M{}).All(&results)
	return results
}

func (m MongoMatchRepository) MatchCollection() (*mgo.Collection, *mgo.Session) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	//	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("ranky").C("matches")
	return c, session
}
