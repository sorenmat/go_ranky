package matchservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sorenmat/ranky/playerservice"

	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Match struct {
	playerOne, playerTwo, playerThree, playerFour string
	scoreOne, scoreTwo                            int8
	playedAt                                      time.Time
}

func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/matches").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	service.Route(service.GET("/").To(AllMatches))
	service.Route(service.GET("/{user-id}").To(FindMatch))
	service.Route(service.PUT("").To(CreateMatch))
	service.Route(service.DELETE("/{user-id}").To(RemoveMatch))

	return service
}

func FindMatch(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")

	result := Match{}
	c, session := matchCollection()
	defer session.Close()

	err := c.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		response.WriteError(http.StatusInternalServerError, err)
	} else {

		response.WriteEntity(result)
	}
}

func AllMatches(request *restful.Request, response *restful.Response) {
	results := make([]Match, 0)

	c, session := matchCollection()
	defer session.Close()
	c.Find(bson.M{}).All(&results)
	response.WriteEntity(results)
}

func matchCollection() (*mgo.Collection, *mgo.Session) {
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

func SaveMatch(usr Match) {
	c, session := matchCollection()
	defer session.Close()

	err := c.Insert(usr)
	if err != nil {
		panic(err)
	}
}

func ValidateMatch(match Match) error {
	if match.playerOne == "" || match.playerTwo == "" {
		return errors.New("Needs atleast two players, to enter a match")
	}
	if match.scoreOne < 10 && match.scoreTwo < 10 {
		return errors.New("Exactly one player or team must have the score of 10")
	}
	if match.scoreOne == 10 && match.scoreTwo == 10 {
		return errors.New("Only one match may reach the score of 10")
	}
	return nil
}

func CreateMatch(request *restful.Request, response *restful.Response) {
	usr := Match{playerOne: "playerOne",
		playerTwo:   "playerTwo",
		playerThree: "playerThree",
		playerFour:  "playerFour",
		playedAt:    time.Now()}

	jsonUser, _ := json.Marshal(usr)
	fmt.Println("--> Initial usr: ", string(jsonUser))

	err := request.ReadEntity(&usr)
	ValidateMatch(usr)
	fmt.Println("Trying to find: ", usr.playerOne)
	pOne := playerservice.FindUser(usr.playerOne)
	fmt.Println("Player one: ", pOne.Name)
	//	fmt.Println("Usr: ", usr.Id)
	SaveMatch(usr)
	// here you would create the user with some persistence system
	if err == nil {
		response.WriteEntity(usr)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func RemoveMatch(request *restful.Request, response *restful.Response) {
	// here you would delete the user from some persistence system
}
