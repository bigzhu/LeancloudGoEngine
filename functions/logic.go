package functions

import (
	"encoding/json"
	"errors"

	"github.com/leancloud/go-sdk/leancloud"
)

func getUri(req *leancloud.FunctionRequest) (string, error) {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return "", errors.New("invalid params")
	}

	uri, ok := params["uri"].(string)
	if !ok {
		return "", errors.New("invalid params")
	}
	return uri, nil
}

func getSentences(uri string) ([]Sentence, error) {
	captions, err := remoteClient.Run("captions", map[string]string{"uri": uri})
	if err != nil {
		return nil, err
	}
	content, err := format(captions.(string))
	if err != nil {
		return nil, err
	}
	//return content, nil
	sentences, err := AnalysisArticle(content)
	if err != nil {
		return nil, err
	}
	return sentences, nil
}

func getVideoInfo(uri string) (VideoInfo, error) {
	videoInfo := VideoInfo{}
	videoInfoRemote, err := remoteClient.Run("videoInfo", map[string]string{"uri": uri})
	if err != nil {
		return videoInfo, err
	}
	videoInfoJson, err := json.Marshal(videoInfoRemote)
	if err != nil {
		return videoInfo, err
	}
	json.Unmarshal([]byte(videoInfoJson), &videoInfo)
	return videoInfo, err
}
func buildArticle(uri string, user *leancloud.User, sentences []Sentence, videoInfo *VideoInfo) Article {
	article := Article{}
	article.Owner = *user
	article.Sentences = sentences
	article.Thumbnail = videoInfo.Thumbnails[len(videoInfo.Thumbnails)-1].URL
	article.Title = videoInfo.Title
	//article.ChannelName = videoInfo.Channel.Name
	article.Channel = videoInfo.Channel.ID
	article.Avatar = videoInfo.Avatar
	article.Youtube = uri
	article.WordCount = computeArticleWordCount(article.Sentences)
	return article
}

func bindUserArticle(user *leancloud.User, article *Article) (UserArticle, error) {
	// bind UserArticle relation
	userArticle := UserArticle{}
	//check exist
	err := client.Class("UserArticle").NewQuery().EqualTo("owner", user).EqualTo("article", article).First(&userArticle)
	//找到,直接返回
	if err == nil {
		// 存在
		if userArticle.IsFollowing != true {
			//更新为关注
			userArticle.IsFollowing = true
			err = client.Class("UserArticle").ID(userArticle.ID).Update(&userArticle)
		}
		return userArticle, err
	}
	//不是这个错,说明查询出问题了
	if err.Error() != "no matched object found" {
		return userArticle, err
	}
	//bind
	userArticle.Article = *article
	userArticle.Owner = *user
	userArticle.IsFollowing = true
	userArticle.AcquiringCount = 1
	_, err = client.Class("UserArticle").Create(&userArticle)

	return userArticle, err
}
