package main

import (
	"reflect"
	"testing"
)

func TestInMemorySessionRepository(t *testing.T) {
	session1 := Session{
		ID:     "sess1",
		UserID: "user1",
		Realm:  "users",
	}

	sessions := []Session{
		session1,
	}
	repo := NewInMemorySessionRepository(sessions)

	newSession, err := repo.CreateSession(Session{
		UserID: "user2",
		Realm:  "users",
	})
	assertError(t, err)
	asserSessionLength(t, *repo, 2)
	assertNotEmpty(t, newSession.ID)

	gotSession, ok := repo.GetSessionByID("sess1")
	assertOk(t, ok)
	if !reflect.DeepEqual(session1, gotSession) {
		t.Fatalf("sessions are not equal, got %v, want %v", gotSession, session1)
	}

	_, okBadID := repo.GetSessionByID("bad id")
	if okBadID {
		t.Fatalf("Should return session not found")
	}

	//Find by userID
	gotSessions := repo.FindByUserId("users", "user2")

	if len(gotSessions) != 1 {
		t.Fatalf("Shood find at least 1 session")
	}

	if gotSessions[0].ID != newSession.ID {
		t.Fatalf("Session IDs should match")
	}

	//Delete session
	err = repo.DeleteSession(newSession.ID)
	assertError(t, err)
	asserSessionLength(t, *repo, 1)

	err = repo.DeleteSession("bad id")
	if err == nil {
		t.Fatalf("should return error")
	}
}

func assertNotEmpty(t *testing.T, got string) {
	t.Helper()
	if got == "" {
		t.Fatalf("string is empty")
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertOk(t *testing.T, ok bool) {
	t.Helper()
	if !ok {
		t.Fatalf("bad result")
	}
}

func asserSessionLength(t *testing.T, repo InMemorySessionRepository, want int) {
	t.Helper()
	if len(repo) != want {
		t.Fatalf("bad session repo length got %v want %v", len(repo), want)
	}
}
