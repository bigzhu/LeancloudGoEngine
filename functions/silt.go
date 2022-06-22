package functions

import (
	"errors"
	"regexp"
	"strings"

	//"github.com/jdkato/prose/v2"
	"github.com/bigzhu/prose"
	"github.com/leancloud/go-sdk/leancloud"
)

func init() {
	leancloud.Engine.Define("slit", silt)
}

func silt(req *leancloud.FunctionRequest) (interface{}, error) {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid params")
	}

	content, ok := params["content"].(string)
	if !ok {
		return nil, errors.New("invalid params")
	}

	content, err := format(content)
	if err != nil {
		return nil, err
	}

	//return content, nil
	Sentences, err := AnalysisArticle(content)
	return Sentences, err
}

func format(content string) (string, error) {
	// 特殊字符的替换
	content = strings.ReplaceAll(content, "&gt;", ">")
	content = strings.ReplaceAll(content, "<i>", " ")
	content = strings.ReplaceAll(content, "</i>", " ")

	// 删除所有换行
	content = strings.ReplaceAll(content, "\n", " ")
	// 替换所有不间断空格 \u00A0
	content = strings.ReplaceAll(content, "\u00A0", " ")
	// 多个空格替换为一个
	BlankReg := regexp.MustCompile("\\s{2,}")
	content = BlankReg.ReplaceAllString(content, " ")
	// 按句子分割
	articlesSentences := ""
	doc, err := prose.NewDocument(content)
	if err != nil {
		return "", err
	}
	sentences := doc.Sentences()
	atLeastSentenceNumber := int(len(doc.Tokens()) / 30)
	if len(sentences) < atLeastSentenceNumber {
		//if false {
		// 按照时间轴进行换行
		startExp := regexp.MustCompile("00[0-9]+\\.[0-9]+00")
		content = startExp.ReplaceAllStringFunc(content, addLineBreakToStartTime)
		// no need break at first line
		content = strings.Replace(content, "\n", "", 1)
	} else {
		// 按照文法句子进行换行
		for _, sent := range sentences {
			articlesSentences += sent.Text + "\n"
		}
		content = articlesSentences
	}
	return content, nil
}

func addLineBreakToStartTime(start string) string {
	return "\n\n" + start
}

// AnalysisArticle spilt article to sentence with words
func AnalysisArticle(article string) (sentences []Sentence, err error) {
	starTimeReg := regexp.MustCompile("00[0-9]+.[0-9]+00")
	doc, err := prose.NewDocument(article)
	if err != nil {
		return
	}
	tokens := doc.Tokens()
	sentence := Sentence{}
	for _, tok := range tokens {
		if starTimeReg.MatchString(tok.Text) {
			if sentence.SeekTo == "" {
				sentence.SeekTo = tok.Text
			}
		} else {
			sentence.Words = append(sentence.Words, tok.Text)
		}
		// 换行就要另起一句
		if tok.Text == "\n" {
			sentences = append(sentences, sentence)
			sentence = Sentence{}
		}
	}
	if len(sentence.Words) > 0 {
		sentences = append(sentences, sentence)
	}
	return
}
