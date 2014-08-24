package matchservice

import "testing"

//	"github.com/sorenmat/ranky/matchservice"

func TestValidateMatch(t *testing.T) {
	match := Match{playerOne: "Soren", playerTwo: "Joe", scoreOne: 5, scoreTwo: 10}
	ValidateMatch(match)
}

func TestOnlyOnePlayerFails(t *testing.T) {
	match := Match{playerOne: "Soren", scoreOne: 5, scoreTwo: 10}
	err := ValidateMatch(match)
	if err == nil || err.Error() != "Needs atleast two players, to enter a match" {
		panic(err)
	}
}

func TestOnlyNoScoreReached10(t *testing.T) {
	match := Match{playerOne: "Soren", playerTwo: "Joe", scoreOne: 5, scoreTwo: 9}
	err := ValidateMatch(match)
	if err == nil || err.Error() != "Exactly one player or team must have the score of 10" {
		panic(err)
	}
}
