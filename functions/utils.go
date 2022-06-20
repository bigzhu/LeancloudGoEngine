package functions

import (
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
