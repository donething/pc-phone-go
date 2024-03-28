package main

import "testing"

func Test_generateRandomToken(t *testing.T) {
	auth, err := generateRandomToken(32)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Authorization:%s\n", auth)
}
