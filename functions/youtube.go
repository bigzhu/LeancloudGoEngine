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
	//check if youtube video is exists
	article := Article{}
	err = client.Class("Article").NewQuery().EqualTo("youtube", uri).First(&article)
	if err == nil {
		return article.ID, err
	}

	//不是这个错,说明查询出问题了
	if err.Error() != "no matched object found" {
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

	article = buildArticle(uri, req.CurrentUser, sentences, &videoInfo)
	_, err = client.Class("Article").Create(&article)
	if err != nil {
		return nil, err
	}

	return article.ID, err
}
