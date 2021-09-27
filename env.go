package green

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

var tagParseRegExp = regexp.MustCompile(`,\s*default=(?P<default>\w+)`)

const envFileLineFormat = `\w+=\w+`

type FieldTag struct {
	Required bool
	Name     string
	Default  string
}

func parseTag(tag string) (ft FieldTag) {
	ft.Required, _ = regexp.MatchString("required", tag)
	ft.Name = strings.Split(tag, ",")[0]
	defaults := tagParseRegExp.FindStringSubmatch(tag)
	if len(defaults) > 1 {
		ft.Default = defaults[1]
	}
	return
}

/*
loadEnvFile reads `filename` into a map by parsing each line
represented with the `key=val` syntax.
*/
func loadEnvFile(filename string) (env map[string]string) {
	env = make(map[string]string)
	varLines, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	s := bufio.NewScanner(strings.NewReader(string(varLines)))
	s.Split(bufio.ScanWords)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		if isKeyVal, _ := regexp.MatchString(line, envFileLineFormat); !isKeyVal {
			continue
		}

		keyVal := strings.Split(line, "=")
		name, val := keyVal[0], keyVal[1]
		env[name] = val
	}

	return
}

/*
LoadEnv accepts a struct with the `green` field tag defined
for its fields and returns a map of environment variables with
defaults if specified. If a field is marked as required but
not found in the environment then LoadEnv panics.

Example struct tag usage:
type MyEnv struct {
	Foo string `green:"FOO,default=myFooDefault"`
	Bar string `green:"BAR,required"`
}

env := green.LoadEnv(MyEnv{})
fmt.Printf("Foo = %s", env["Foo"])
*/
func LoadEnv(env interface{}) (results map[string]string) {
	results = make(map[string]string)

	wd, err := os.Getwd()
	if err != nil {
		return
	} else {
		filename := filepath.Join(wd, ".env")
		results = loadEnvFile(filename)
	}

	val := reflect.ValueOf(env)
	ifc := reflect.Indirect(val)

	for _, field := range reflect.VisibleFields(ifc.Type()) {
		ft := parseTag(field.Tag.Get("green"))
		val := os.Getenv(ft.Name)

		if len(val) == 0 {
			if ft.Required {
				panic(fmt.Sprintf("%s not found in environment", ft.Name))
			} else if ft.Default != "" {
				results[ft.Name] = ft.Default
			}
		}
	}

	return
}

/*
UnmarshalENV accepts a struct with the `green` field tag defined
for its fields and assigns values from the environment or the defaults
to `env`.
*/
// func UnmarshalENV(env interface{}) (err error) {
// 	return
// }
