package main

import "testing"

func TestClean(t *testing.T) {
	got := handleClean("This is a kerfuffle opinion I need to share with the world")
	want := "This is a **** opinion I need to share with the world"

	if got != want {
		t.Errorf("test failed, got %v, want %v", got, want)
	}
}