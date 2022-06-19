package functions

import (
	"errors"

	"github.com/leancloud/go-sdk/leancloud"
)

func init() {
	leancloud.Engine.Define("hello", hello)
}

func hello(req *leancloud.FunctionRequest) (interface{}, error) {
	params, ok := req.Params.(map[string]string)

	if !ok {
		return nil, errors.New("invalid params")
	}

	return map[string]string{
		"hello": params["name"],
	}, nil
}
