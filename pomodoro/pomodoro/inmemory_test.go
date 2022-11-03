package pomodoro_test

import (
	"testing"

	"achristie.net/pomodoro/pomodoro"
	"achristie.net/pomodoro/pomodoro/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	return repository.NewInMemoryRepo(), func() {}
}
