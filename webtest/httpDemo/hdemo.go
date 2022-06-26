package main

import (
	"log"
	"net/http"
)

type StubPlayerStore struct {
	score    map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.score[name]
	return score
}

func (s *StubPlayerStore) GetLeague() []Player {
	return s.league
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func GetPlayerScore(name string) string {
	if name == "Pepper" {
		return "20"
	}

	if name == "Floyd" {
		return "10"
	}
	return ""
}

func main() {
	server := &PlayerServer{&InMemoryPlayerStore{}, nil}
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
