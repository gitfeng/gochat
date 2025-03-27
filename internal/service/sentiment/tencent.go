package sentiment

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
)

var (
	nlpClient    *nlp.Client
	keywordsOnce sync.Once
	sentimentMap = map[string]int{
		"positive": 1,
		"neutral":  0,
		"negative": -1,
	}
)

func init() {
	keywordsOnce.Do(func() {
		viper.SetConfigName("feedback_keywords")
		viper.AddConfigPath("./config")
		if err := viper.ReadInConfig(); err != nil {
			log.Printf("无法读取关键词配置，使用默认值: %v", err)
			return
		}

		secretId := viper.GetString("tencentcloud.secret_id")
		secretKey := viper.GetString("tencentcloud.secret_key")
		var err error
		nlpClient, err = nlp.NewClientWithSecretId(secretId, secretKey, "ap-guangzhou")
		if err != nil {
			log.Printf("无法初始化腾讯云客户端: %v", err)
			return
		}
	})
}

func AnalyzeSentiment(comment string) (int, error) {

	req := nlp.NewAnalyzeSentimentRequest()
	req.Text = common.StringPtr(comment)
	resp, err := nlpClient.AnalyzeSentiment(req)
	if err != nil {
		return 0, err
	}

	val, exists := sentimentMap[*resp.Response.Sentiment]
	if !exists {
		return 0, fmt.Errorf("unknown sentiment: %s", *resp.Response.Sentiment)
	}
	return val, nil
}
