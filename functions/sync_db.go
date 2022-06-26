package functions

import (
	"fmt"

	"github.com/bigzhu/gobz/modelsbz"
	"github.com/leancloud/go-sdk/leancloud"
)

// Learned 记录学过的单词
type Learned struct {
	modelsbz.Base
	Word    string `gorm:"unique_index:unique_word_user_id" json:"word" binding:"required"`
	Count   int    `gorm:"default:0" json:"count"`                                    // 学过多少次
	Learned bool   `gorm:"default: false" json:"learned"`                             //是否学会
	UserID  int    `gorm:"unique_index:unique_word_user_id;not null" json:"user_id" ` // 用户ID
}

func init() {
	client = leancloud.NewEnvClient()
	leancloud.Engine.Define("sync_db", sync_db)
}

func sync_db(req *leancloud.FunctionRequest) (interface{}, error) {
	modelsbz.Connect()
	learneds := []Learned{}
	err := modelsbz.DB.Where("user_id = 1 AND learned=false").Find(&learneds).Error
	if err != nil {
		return nil, err
	}
	for _, l := range learneds {
		acquiringWord := AcquiringWord{
			Word:  l.Word,
			Count: l.Count,
			Done:  &l.Learned,
			Owner: *req.CurrentUser,
		}

		fmt.Printf("%v %v \n", l.Word, l.Learned)
		if _, err := client.Class("AcquiringWord").Create(&acquiringWord); err != nil {
			return nil, err
		}
	}
	return nil, err
}
