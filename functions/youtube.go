package functions

import (
	"github.com/leancloud/go-sdk/leancloud"
)

func init() {
	remoteClient = createRemoteClient()
	client = leancloud.NewEnvClient()
	leancloud.Engine.Define("youtube", youtube)
}

func youtube(req *leancloud.FunctionRequest) (interface{}, error) {
	uri, err := getUri(req)
	if err != nil {
		return nil, err
	}

	sentences, err := getSentences(uri)
	if err != nil {
		return nil, err
	}
	videoInfo, err := getVideoInfo(uri)
	if err != nil {
		return nil, err
	}

	article := createArticle(uri, req.CurrentUser, sentences, &videoInfo)
	if _, err := client.Class("Article").Create(&article); err != nil {
		return nil, err
	}
	return article, err
}
