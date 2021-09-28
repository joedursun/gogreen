# Green

Keep your environment green! The green package provides a way to add runtime checks for your environment variables (AKA env vars). Its purpose is to provide guarantees around environment variable values and to help standardize configuration via env vars. If a struct field is marked as `required` then if it is not found in either the environment or the `.env` file `green.LoadEnv` will panic.

Features:
  - require env vars at runtime
  - provide defaults for env vars
  - load env vars from .env file

Order of precedence:
  1. environment variables
  2. .env file variables



## Example Usage

```go
// Use struct tags to inform green to track these fields. It's fine
// to use alongside JSON struct tags as well.
type MyEnv struct {
	Hostname   string `json:"-" green:"HOSTNAME,required"`
	Database   string `json:"database" green:"DATABASE,default=myuser"`
	AuthToken  string `green:"AUTH_TOKEN"`
}

// assuming HOSTNAME is present in the environment and nothing else
env := green.UnmarshalENV(MyEnv{})
env.Hostname // same as os.Getenv("HOSTNAME")
env.Database // "myuser"
env.AuthToken // ""
```
