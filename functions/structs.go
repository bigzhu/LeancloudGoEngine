package functions

import (
	"github.com/leancloud/go-sdk/leancloud"
)

type Article struct {
	leancloud.Object
	Owner     leancloud.User `json:"owner"`
	Sentences []Sentence     `json:"sentences"`
	WordCount int            `json:"wordCount"`
	Title     string         `json:"title"`
	Youtube   string         `json:"youtube"`
	Avatar    string         `json:"avatar"`
	Channel   string         `json:"channel"`
	Thumbnail string         `json:"thumbnail"`
}
type UserArticle struct {
	leancloud.Object
	Owner          leancloud.User `json:"owner"`
	Article        Article        `json:"article"`
	AcquiringCount int            `json:"acquiringCount"`
	IsFollowing    bool           `json:"isFollowing"`
}

// Sentence 存储数据库的句子
type Sentence struct {
	SeekTo string   `json:"seekTo"`
	Words  []string `json:"words"`
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
