package matchservice

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sorenmat/ranky/playerservice"

	"github.com/emicklei/go-restful"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Match struct {
	PlayerOne, PlayerTwo, PlayerThree, PlayerFour string
	ScoreOne, ScoreTwo                            int8
	PlayedAt                                      time.Time
}

func (m Match) toString() string {
	json, _ := json.Marshal(m)
	return string(json)
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
	if match.PlayerOne == "" || match.PlayerTwo == "" {
		return errors.New("Needs atleast two players, to enter a match")
	}
	if match.ScoreOne < 10 && match.ScoreTwo < 10 {
		return errors.New("Exactly one player or team must have the score of 10")
	}
	if match.ScoreOne == 10 && match.ScoreTwo == 10 {
		return errors.New("Only one match may reach the score of 10")
	}
	return nil
}

func ValidatePlayersInMatch(match Match) {
	if match.PlayerOne == "" && match.PlayerTwo == "" {
		panic("A Match needs atleast two players")
	}
	if match.PlayerThree != "" && match.PlayerFour == "" {
		panic("A fourth player is needed")
	}

	if match.PlayerThree == "" && match.PlayerFour != "" {
		panic("A thrid player is needed")
	}

	playerservice.FindUser(match.PlayerOne)
	playerservice.FindUser(match.PlayerTwo)
	if match.PlayerThree != "" {
		playerservice.FindUser(match.PlayerThree)
		playerservice.FindUser(match.PlayerFour)

	}
}

func CreateMatch(request *restful.Request, response *restful.Response) {
	match := Match{PlayerOne: "PlayerOne",
		PlayerTwo:   "PlayerTwo",
		PlayerThree: "playerThree",
		PlayerFour:  "playerFour",
		PlayedAt:    time.Now()}

	//jsonUser, _ := json.Marshal(usr)
	//fmt.Println("--> Initial usr: ", string(jsonUser))

	err := request.ReadEntity(&match)
	match.PlayedAt = time.Now()
	ValidateMatch(match)
	ValidatePlayersInMatch(match)
	//	fmt.Printf("JSON Match payload: %s\n", match.toString())
	//	fmt.Println("Trying to find: ", match.PlayerOne)

	SaveMatch(match)
	// here you would create the user with some persistence system
	if err == nil {
		response.WriteEntity(match)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func RemoveMatch(request *restful.Request, response *restful.Response) {
	// here you would delete the user from some persistence system
}
