package matchservice

type MatchRepository interface {
	SaveMatch(match Match)
	FindMatch(id string) Match
	FindAllMatches() []Match
}
