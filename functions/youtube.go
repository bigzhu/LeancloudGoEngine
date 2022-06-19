package functions

import (
	"errors"

	"github.com/leancloud/go-sdk/leancloud"
)

var client *leancloud.Client

func init() {
	// create client, class function from http
	client = leancloud.NewEnvClient()
	leancloud.Engine.Define("youtube", youtube)
}

func youtube(req *leancloud.FunctionRequest) (interface{}, error) {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid params")
	}

	uri, ok := params["uri"].(string)
	if !ok {
		return nil, errors.New("invalid params")
	}

	//captions, err := leancloud.Engine.Run("captions", map[string]string{"uri": uri})
	captions, err := client.Run("captions", map[string]string{"uri": uri})
	if err != nil {
		panic(err)
	}
	content, err := format(captions.(string))
	if err != nil {
		return nil, err
	}

	//return content, nil
	Sentences, err := AnalysisArticle(content)

	return Sentences, nil
}
