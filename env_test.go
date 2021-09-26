package green_test

import (
	"testing"

	"github.com/joedursun/green"
)

type TestEnv struct {
	Hostname        string `json:"-" green:"GREEN_TEST_HOSTNAME,required"`
	Database        string `json:"database" green:"DATABASE,default=myuser"`
	EmptyVal        string `green:"EMPTY_VAL"`
	unexportedField string
}

func testPanicString(t *testing.T, expected string) {
	if r := recover(); r != nil {
		if r != expected {
			t.Errorf("Received unexpected panic message: %s", r)
		}
	} else {
		t.Error("Expected to receive panic")
	}
}

func TestLoadEnv(t *testing.T) {
	defer testPanicString(t, "GREEN_TEST_HOSTNAME not found in environment")

	tenv := TestEnv{unexportedField: "TEST"}
	res := green.LoadEnv(tenv)
	if res["Database"] != "myuser" {
		t.Error("Expected Database to be present")
	}

	if res["EmptyVal"] != "" {
		t.Errorf("Expected Database to be an empty string but received %s", res["EmptyVal"])
	}

	if res["unexportedField"] != "TEST" {
		t.Errorf("Expected unexportedField to be \"TEST\" but received %s", res["EmptyVal"])
	}
}
