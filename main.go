package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

const defaultServerConfigPath = "config.yaml"

type Response struct {
	Path       string      `yaml:"path"`
	StatusCode int         `yaml:"statusCode"`
	Response   string      `yaml:"response"`
	Headers    http.Header `yaml:"headers"`
}

type Responses = []Response

type ServerConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Responses *Responses
}

const (
	contentTypeHeader     = "Content-Type"
	applicationJsonHeader = "application/json"
	textPlainHeader       = "text/plain"
)

var availableContentTypeHeaderValues = map[string]struct{}{
	applicationJsonHeader: {},
	textPlainHeader:       {},
}

func loadServerConfig(serverConfigFile string) (*ServerConfig, error) {
	configData, err := os.ReadFile(serverConfigFile)
	if err != nil {
		return nil, err
	}

	var serverConfig *ServerConfig
	if err = yaml.Unmarshal(configData, &serverConfig); err != nil {
		return nil, err
	}

	return serverConfig, nil
}

func createHandler(response *Response) (func(c *gin.Context), error) {
	contentTypeHeaderValues := response.Headers["Content-Type"]
	// Content-Type ヘッダーの値が availableContentTypeHeaderValues に含まれない場合はエラー
	for _, contentTypeHeaderValue := range contentTypeHeaderValues {
		if _, ok := availableContentTypeHeaderValues[contentTypeHeaderValue]; !ok {
			return nil, fmt.Errorf("the Content-Type header value is not available: %s", contentTypeHeaderValue)
		}
	}
	return func(c *gin.Context) {
		// ヘッダーの設定
		// c.JSON などでレスポンスを書き込む前にヘッダーを設定しなければならなかった
		// c.JSON などの後にヘッダーを設定するとヘッダーが出力されなかった
		for key, values := range response.Headers {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		statusCode := response.StatusCode
		// Content-Type ヘッダーの値ごとにレスポンスを設定
		for _, contentTypeHeaderValue := range contentTypeHeaderValues {
			if contentTypeHeaderValue == applicationJsonHeader {
				// 文字列になっている response.Response を map に変換
				// どのようなデータが response に入るのか不明なため,
				// 型を interface{} とする
				var jsonResponse interface{}
				if err := json.Unmarshal([]byte(response.Response), &jsonResponse); err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, err)
				}
				c.JSON(statusCode, jsonResponse)
			} else if contentTypeHeaderValue == textPlainHeader {
				c.String(statusCode, response.Response)
			}
		}
	}, nil
}

// 同じ path を engine.GET で指定するとエラーになる
// たとえば, engine.GET("/status", handler1) と engine.GET("/status", handler2) を実行するとエラーになる
// エラー内容: panic: handlers are already registered for path '/status'
// i.e., testdata/configWithSamePaths.yaml でサーバーを起動するとエラーになる
func runServer(serverConfigFile string) error {
	serverConfig, err := loadServerConfig(serverConfigFile)
	if err != nil {
		return err
	}
	fmt.Printf("serverConfig.Responses: %v\n", serverConfig.Responses)

	engine := gin.Default()
	for i := 0; i < len(*serverConfig.Responses); i++ {
		response := (*serverConfig.Responses)[i]
		handler, err := createHandler(&response)
		if err != nil {
			return err
		}
		engine.GET(response.Path, handler)
	}

	address := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
	engine.Run(address)
	return nil
}

func main() {
	// TODO: cobra に変更
	serverConfigFile := flag.String("config", defaultServerConfigPath, "a server config file path")
	flag.Parse()

	if err := runServer(*serverConfigFile); err != nil {
		log.Fatal(err.Error())
	}
}
