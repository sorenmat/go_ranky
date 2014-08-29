package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/emicklei/go-restful"
	"github.com/sorenmat/ranky/matchservice"
	"github.com/sorenmat/ranky/playerservice"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {
	fmt.Println("Starting suite")
	restful.Add(playerservice.New(playerservice.mongoplayerrepository{}))
	restful.Add(matchservice.New())
	go http.ListenAndServe(":8080", nil)

}

func createPlayer(name string) string {
	res, _ := doPut("http://localhost:8080/players", "{\"Id\": \"\", \"Name\": \""+name+"\"}")
	user := playerservice.User{}
	json.Unmarshal([]byte(res), &user)
	fmt.Printf("Res %s\n", res)
	return user.Id

}
func createTwoPlayerMatch(playerOne string, playerTwo string, scoreOne int8, scoreTwo int8) int {
	m := matchservice.Match{PlayerOne: playerOne, PlayerTwo: playerTwo, ScoreOne: scoreOne, ScoreTwo: scoreTwo}
	json, _ := json.Marshal(m)
	_, status := doPut("http://localhost:8080/matches", string(json))
	return status
}

func createFourPlayerMatch(playerOne string, playerTwo string, playerThree string, playerFour string, scoreOne int8, scoreTwo int8) int {
	m := matchservice.Match{
		PlayerOne:   playerOne,
		PlayerTwo:   playerTwo,
		PlayerThree: playerThree,
		PlayerFour:  playerFour,
		ScoreOne:    scoreOne,
		ScoreTwo:    scoreTwo}
	json, _ := json.Marshal(m)
	_, status := doPut("http://localhost:8080/matches", string(json))
	return status
}

func (s *MySuite) TestCreatingPlayersAndAMatch(c *C) {
	sorenId := createPlayer("Soren")
	joeId := createPlayer("Joe")

	jackId := createPlayer("Jack")
	jillId := createPlayer("Jill")

	c.Assert(createTwoPlayerMatch(sorenId, joeId, 10, 2), Equals, 200)
	c.Assert(createTwoPlayerMatch(sorenId, joeId, 10, 6), Equals, 200)

	c.Assert(createFourPlayerMatch(sorenId, joeId, jackId, jillId, 10, 7), Equals, 200)
}

func (s *MySuite) TestCreateSamePlayerTwice(c *C) {
	//createPlayer("Soren")
	//createPlayer("Soren")
}

func doPut(url string, body string) (string, int) {
	client := &http.Client{}
	request, err := http.NewRequest("PUT", url, strings.NewReader(body))
	request.Header.Add("Content-type", "application/json")
	request.ContentLength = int64(len(body))
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return "", response.StatusCode
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		return string(contents), response.StatusCode
	}
}
