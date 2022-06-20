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
func createArticle(uri string, user *leancloud.User, sentences []Sentence, videoInfo *VideoInfo) Article {
	article := Article{}
	article.Owner = user
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
