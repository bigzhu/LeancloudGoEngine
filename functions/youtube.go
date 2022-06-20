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
	//找到,直接返回
	if err == nil {
		return article, nil
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

	article = createArticle(uri, req.CurrentUser, sentences, &videoInfo)
	if _, err := client.Class("Article").Create(&article); err != nil {
		return nil, err
	}
	return article, err
}
