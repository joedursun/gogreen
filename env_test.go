package green_test

import (
	"os"
	"path/filepath"
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

func TestLoadEnvFile(t *testing.T) {
	expected := map[string]string{
		"FOO":      "bar",
		"USERNAME": "guest",
		"TOKEN":    "abc$@H9876;",
		"Hello":    "World!",
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	filename := filepath.Join(wd, ".env.example")
	res := green.LoadEnvFile(filename)
	for key, val := range res {
		expectedVal, found := expected[key]

		if !found {
			t.Errorf("Expected key %s to be present", key)
			continue
		}

		if expectedVal != val {
			t.Errorf("Expected %s but got %s", expectedVal, val)
		}
	}
}
