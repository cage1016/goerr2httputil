# goerr2httphutil

> convert const to http error struct

`custom_error2/custom_error2.go`

```go
package custom_error2

type ErrorCode uint

// 400
const (
	BroadcastLangErrorQueryBroadcastLangRepoFindBadRequest ErrorCode = iota + 4001200
)

// 500
const (
	BroadcastLangErrorQueryBroadcastLangRepoFindServerInternalError ErrorCode = iota + 5001200 // custom error message
)

```

`custom_error2/custom_error2_custom_string.go`

```go
// Automatically-generated file. Do not edit
package custom_error2

//BroadcastLangErrorQueryBroadcastLangRepoFindBadRequest
type HTTPError4001200 struct {
	ApiVersion string `json:"apiVersion" example:"1.1"`
	Error      struct {
		Code    int    `json:"code" example:"4001200"`
		Message string `json:"message" example:""`
	} `json:"error"`
}

//BroadcastLangErrorQueryBroadcastLangRepoFindServerInternalError
type HTTPError5001200 struct {
	ApiVersion string `json:"apiVersion" example:"1.1"`
	Error      struct {
		Code    int    `json:"code" example:"5001200"`
		Message string `json:"message" example:"custom error message"`
	} `json:"error"`
}

```


## Install

```bash
go get github.com/cage1016/goerr2httputil
```

### Usage

```bash
# default will create custom_error2_custom_string.go in current folder
goerr2httputil -type custom_error2

# you could export custom_error2_custom_string.go to any path you want
goerr2httputil -type custom_error2 -export path-you-want

# you could define apiVersion
goerr2httputil -type custom_error2 -export path-you-want -version 1.1
```