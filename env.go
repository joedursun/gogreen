package gogreen

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var tagParseRegExp = regexp.MustCompile(`,\s*default=(?P<default>[\w!@#$%^&*()]+)`)
var envFileLineFormat = regexp.MustCompile(`\w+=\w+`)

type Environment interface {
	EnvFileLocation() string
}

type FieldTag struct {
	Required   bool
	EnvVarName string
	FieldName  string
	Default    string
}

func parseTag(field reflect.StructField) (ft FieldTag) {
	tag := field.Tag.Get("green")
	if tag == "" {
		return
	}

	ft.FieldName = field.Name
	ft.Required, _ = regexp.MatchString("required", tag)
	ft.EnvVarName = strings.Split(tag, ",")[0]
	defaults := tagParseRegExp.FindStringSubmatch(tag)
	if len(defaults) > 1 {
		ft.Default = defaults[1]
	}
	return
}

/*
loadEnvFile reads `filename` into a map by parsing each line
represented with the `key=val` syntax. Lines starting with `#`
are treated as comments, but if `#` appears after the first character
*/
func LoadEnvFile(filename string) (env map[string]string) {
	env = make(map[string]string)
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := bytes.TrimSpace(s.Bytes())
		if bytes.HasPrefix(line, []byte("#")) {
			continue
		}

		if isKeyVal := envFileLineFormat.Match(line); !isKeyVal {
			continue
		}

		keyVal := bytes.Split(line, []byte("="))
		name, val := keyVal[0], keyVal[1]
		env[string(name)] = string(val)
	}

	return
}

func getStructFields(env Environment) []reflect.StructField {
	val := reflect.ValueOf(env)
	ifc := reflect.Indirect(val)

	return reflect.VisibleFields(ifc.Type())
}

/*
LoadEnv accepts a struct with the `green` field tag defined
for its fields and returns a map of environment variables with
defaults if specified. If a field is marked as required but
not found in the environment then LoadEnv panics.
*/
func LoadEnv(env Environment) (results map[string]string) {
	results = LoadEnvFile(env.EnvFileLocation())
	fields := getStructFields(env)

	for _, field := range fields {
		ft := parseTag(field)
		val := os.Getenv(ft.EnvVarName)

		if len(val) == 0 {
			if ft.Required {
				panic(fmt.Sprintf("%s not found in environment", ft.EnvVarName))
			} else if ft.Default != "" {
				results[ft.EnvVarName] = ft.Default
			}
			continue
		}
		results[ft.EnvVarName] = val
	}

	return
}

/*
UnmarshalENV accepts a struct with the `green` field tag defined
for its fields and assigns values from the environment or the defaults
to `env`.
Since os.Getenv() always returns a string we leave conversion to other
data types up to the caller. Any field whose type isn't a string
will be skipped as will any field that doesn't have a `green` struct tag.
*/
func UnmarshalENV(env Environment) (err error) {
	defer func() {
		if r := recover(); r != nil {
			envType := reflect.TypeOf(env).String()
			err = fmt.Errorf("must provide pointer to struct but given %s", envType)
		}
	}()

	envStruct := reflect.ValueOf(env).Elem()
	if !envStruct.CanAddr() {
		return errors.New("argument not addressable")
	}

	values := LoadEnv(env)
	fields := getStructFields(env)

	for _, field := range fields {
		if field.Type.String() != "string" {
			continue
		}

		ft := parseTag(field)
		val, found := values[ft.EnvVarName]
		if !found {
			continue
		}

		fieldVal := envStruct.FieldByName(ft.FieldName)
		fieldVal.Set(reflect.ValueOf(val))
	}
	return
}
