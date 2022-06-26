package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func newGetScoreRequset(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func TestHttp(t *testing.T) {

	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequset("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "20"

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)
	})

	t.Run("return Floyd's score ", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := "10"
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, got, want)
	})

	t.Run("return 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequset("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		got := response.Code
		want := http.StatusNotFound

		if got != want {
			t.Errorf("got status %d want %d", got, want)
		}

	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not correct status ,got %d ,want %d", got, want)
	}
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s' ,want '%s'", got, want)
	}
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("it returns accpted on Post", func(t *testing.T) {
		player := "Pepper"
		//request, _ := http.NewRequest(http.MethodPost, "/plauers/Pepper", nil)
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got '%s' want '%s'", store.winCalls[0], player)
		}
	})
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s,", name), nil)
	return req
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := InMemoryPlayerStore{}
	server := NewPlayerServer(&store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequset(player))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}

func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := NewPlayerServer(&store)

	t.Run("it returns 200 on/league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		var got []PlayerServer

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from server '%s' into slice of Player ,'%v'", response.Body, err)
		}

		assertStatus(t, response.Code, http.StatusOK)

	})
	t.Run("json test", func(t *testing.T) {
		wantedLeague := []Player{
			{"cleo", 32},
			{"Chirs", 20},
			{"Tiest", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)
		request := newLeagueRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		got := getLeagueFromREsponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)

		if response.Header().Get("content-type") != "appltcation/json" {
			t.Errorf("response got %v", response.HeaderMap)
		}
	})
}

func getLeagueFromREsponse(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from server '%s' into slice of Player, '%v'", body, err)
	}
	return

}

func assertLeague(t *testing.T, got, want []Player) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}
