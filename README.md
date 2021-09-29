# GoGreen

[![Go Report Card](https://goreportcard.com/badge/github.com/joedursun/gogreen)](https://goreportcard.com/report/github.com/joedursun/gogreen)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/joedursun/gogreen)](https://pkg.go.dev/github.com/joedursun/gogreen)
[![Build Status](https://travis-ci.org/joedursun/gogreen.svg?branch=main)](https://travis-ci.org/joedursun/gogreen)

Keep your environment green! The green package provides a way to add runtime checks for your environment variables (AKA env vars). Its purpose is to provide guarantees around environment variable values and to help standardize configuration via env vars. If a struct field is marked as `required` then if it is not found in either the environment or the `.env` file `green.LoadEnv` will panic.

Features:
  - require env vars at runtime
  - provide defaults for env vars
  - load env vars from .env file

Order of precedence:
  1. environment variables
  2. .env file variables

### Motivation

Many apps load some portion of their configuration from environment variables and often have sprawling logic to check for their presence and provide defaults. By centralizing environment lookup and validation we get a single struct that guarantees all necessary variables are accounted for and can be tested against.



## Example Usage

```go
package main

import(
  "github.com/joedursun/gogreen"
)

// Use struct tags to inform green to track these fields. It's fine
// to use alongside JSON struct tags as well.
type MyEnv struct {
	Hostname   string `json:"-" green:"HOSTNAME,required"`
	Database   string `json:"database" green:"DATABASE,default=myuser"`
	AuthToken  string `green:"AUTH_TOKEN"`
}

// assuming HOSTNAME is present in the environment and nothing else
myenv := MyEnv{}
env := gogreen.UnmarshalENV(&myenv)
env.Hostname // same as os.Getenv("HOSTNAME")
env.Database // "myuser"
env.AuthToken // ""
```
