package gogreen_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/joedursun/gogreen"
)

type TestRequiredEnv struct {
	Hostname string `json:"-" green:"GREEN_TEST_HOSTNAME,required"`
}

func (e TestRequiredEnv) EnvFileLocation() string {
	return ""
}

type TestEnv struct {
	Database     string `json:"database" green:"DATABASE,default=myuser"`
	EmptyVal     string `green:"EMPTY_VAL"`
	SpecialChars string `green:"SPECIAL_CHARS,default=special!$@#"`
}

func (e TestEnv) EnvFileLocation() string {
	return ""
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
	tests := []struct {
		varName string
		varVal  string
	}{
		{"DATABASE", "myuser"},
		{"EMPTY_VAL", ""},
		{"SPECIAL_CHARS", "special!$@#"},
	}

	tenv := TestEnv{}
	res := gogreen.LoadEnv(tenv)

	for _, tt := range tests {
		actual := res[tt.varName]
		expected := tt.varVal
		if actual != expected {
			t.Errorf("expected %s to be %s but got %s", tt.varName, expected, actual)
		}
	}

	defer testPanicString(t, "GREEN_TEST_HOSTNAME not found in environment")

	env := TestRequiredEnv{}
	gogreen.LoadEnv(env) // this should panic and be caught by the deferred test condition
}

func TestLoadEnvFile(t *testing.T) {
	tests := []struct {
		VarName     string
		ExpectedVal string
	}{
		{"FOO", "bar"},
		{"USERNAME", "guest"},
		{"TOKEN", "#abc$@H9876;"},
		{"Hello", "World!"},
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	filename := filepath.Join(wd, ".env.example")
	res, err := gogreen.LoadEnvFile(filename)
	if err != nil {
		t.Error("Unable to read env file: ", err)
	}

	for _, tt := range tests {
		if res[tt.VarName] != tt.ExpectedVal {
			t.Errorf("expected %s but got %s", tt.ExpectedVal, res[tt.VarName])
		}
	}
}

type TestUnmarshalEnv struct {
	GreenRequiredField string `green:"GREEN_REQUIRED_FIELD,required"`
	Foo                string `green:"FOO,default=foo"`
	Username           string `green:"USERNAME,default=guestDefault"`
	EmptyVal           string `green:"EMPTY_VAL"`
	SomeInt            int    `green:"SOME_INT"` // meant to be ignored; defined here to ensure no errors are thrown
	unexportedField    string // meant to be ignored; defined here to ensure no errors are thrown
}

func (e TestUnmarshalEnv) EnvFileLocation() string {
	return "./env.example"
}

func TestUnmarshalENV(t *testing.T) {
	os.Setenv("GREEN_REQUIRED_FIELD", "required_value")
	te := TestUnmarshalEnv{}
	err := gogreen.UnmarshalENV(&te)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		Actual   string
		Expected string
	}{
		{te.Username, "guestDefault"},
		{te.Foo, "foo"},
		{te.GreenRequiredField, "required_value"},
	}

	if te.SomeInt != 0 {
		t.Errorf("Expected SomeInt = 0 but got '%d'", te.SomeInt)
	}

	for _, tt := range tests {
		if tt.Actual != tt.Expected {
			t.Errorf("expected %s but got %s", tt.Expected, tt.Actual)
		}
	}
}

func TestUnmarshalENVBadInput(t *testing.T) {
	err := gogreen.UnmarshalENV(TestUnmarshalEnv{})
	if err == nil {
		t.Error("expected to receive error from UnmarshalENV")
	}

	if err.Error() != "must provide pointer to struct but given gogreen_test.TestUnmarshalEnv" {
		t.Errorf("received unexpected error message: '%s'", err.Error())
	}
}

type ExampleEnv struct {
	Token    string `green:"Token,default=abcdef1234!"`
	Username string `green:"USERNAME,default=guestDefault"`
	Hostname string `green:"GREENTEST_HOSTNAME,required"`
}

func (e ExampleEnv) EnvFileLocation() string {
	return ""
}

func ExampleUnmarshalENV() {
	/*
		type ExampleEnv struct {
			Token    string `green:"TOKEN,default=abcdef1234!"`
			Username string `green:"USERNAME,default=guestDefault"`
			Hostname string `green:"GREENTEST_HOSTNAME,required"`
		}
	*/
	os.Setenv("GREENTEST_HOSTNAME", "localhost")
	env := ExampleEnv{}
	gogreen.UnmarshalENV(&env)

	fmt.Printf("Token: %s\n", env.Token)
	fmt.Printf("Username: %s\n", env.Username)
	fmt.Printf("Hostname: %s\n", env.Hostname)
	// Output:
	// Token: abcdef1234!
	// Username: guestDefault
	// Hostname: localhost
}
