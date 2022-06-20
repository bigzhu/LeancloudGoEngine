package functions

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/leancloud/go-sdk/leancloud"
)

var remoteClient *leancloud.Client
var client *leancloud.Client

func createRemoteClient() *leancloud.Client {
	options := &leancloud.ClientOptions{
		AppID:     "RFO2ogo9Li9HhyPaYDQloRUL-MdYXbMMI",
		AppKey:    "NHW1gCDPBExOhEVMxPnzBabE",
		MasterKey: "FOW0ngVSritVU5BWTBGl4Rhj",
		ServerURL: "https://bigzhu.avosapps.us",
	}

	return leancloud.NewClient(options)
}

func init() {
	// create client, class function from http
	remoteClient = createRemoteClient()
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
	//captions, err := client.Run("captions", map[string]string{"uri": uri})
	//captions := ""
	//err := client.RPC("captions", map[string]string{"uri": uri}, &captions)
	captions, err := remoteClient.Run("captions", map[string]string{"uri": uri})
	if err != nil {
		panic(err)
	}
	content, err := format(captions.(string))
	if err != nil {
		return nil, err
	}

	//return content, nil
	Sentences, err := AnalysisArticle(content)
	if err != nil {
		return nil, err
	}

	article := Article{}
	article.Owner = req.CurrentUser
	article.Sentences = Sentences

	videoInfoRemote, err := remoteClient.Run("videoInfo", map[string]string{"uri": uri})
	if err != nil {
		return nil, err
	}
	videoInfoJson, err := json.Marshal(videoInfoRemote)
	if err != nil {
		return nil, err
	}
	videoInfo := VideoInfo{}
	json.Unmarshal([]byte(videoInfoJson), &videoInfo)

	article.Thumbnail = videoInfo.Thumbnails[len(videoInfo.Thumbnails)-1].URL
	article.Title = videoInfo.Title
	//article.ChannelName = videoInfo.Channel.Name
	article.Channel = videoInfo.Channel.ID
	article.Avatar = videoInfo.Avatar
	article.Youtube = uri
	article.WordCount = computeArticleWordCount(article.Sentences)

	if _, err := client.Class("Article").Create(&article); err != nil {
		return nil, err
	}
	return article, err
}

type Article struct {
	leancloud.Object
	Owner     *leancloud.User `json:"owner"`
	Sentences []Sentence      `json:"sentences"`
	WordCount int             `json:"wordCount"`
	Title     string          `json:"title"`
	Youtube   string          `json:"youtube"`
	Avatar    string          `json:"avatar"`
	Channel   string          `json:"channel"`
	Thumbnail string          `json:"thumbnail"`
}

type VideoInfo struct {
	Avatar        string  `json:"avatar"`
	AverageRating float64 `json:"averageRating"`
	Channel       struct {
		ID   string `json:"id"`
		Link string `json:"link"`
		Name string `json:"name"`
	} `json:"channel"`
	Description string   `json:"description"`
	ID          string   `json:"id"`
	Keywords    []string `json:"keywords"`
	Link        string   `json:"link"`
	PublishDate string   `json:"publishDate"`
	Thumbnails  []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"thumbnails"`
	Title      string `json:"title"`
	UploadDate string `json:"uploadDate"`
	ViewCount  struct {
		Text string `json:"text"`
	} `json:"viewCount"`
}

// 计算文章有多少个单词
func computeArticleWordCount(sentences []Sentence) int {
	wordUnique := make(map[string]string)
	wordReg := regexp.MustCompile("[a-zA-Z]+")
	for _, sentence := range sentences {
		for _, word := range sentence.Words {
			// 替换换行, 去除空格, 转小写
			lowWord := strings.ToLower(strings.TrimSpace(word))
			if wordReg.MatchString(lowWord) && len(lowWord) > 1 {
				wordUnique[lowWord] = ""
			}
		}
	}
	return len(wordUnique)
}
