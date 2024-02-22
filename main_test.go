package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoadServerConfig(t *testing.T) {
	testCases := []struct {
		name                 string
		serverConfigFile     string
		expectedServerConfig *ServerConfig
		expectedError        error
	}{
		{
			name:             "ShouldReturnResponsesDefinedInTheConfigFile",
			serverConfigFile: "./testdata/config.yaml",
			expectedServerConfig: &ServerConfig{
				Host: "0.0.0.0",
				Port: 8080,
				Responses: &Responses{
					{
						Path:       "/status",
						StatusCode: 200,
						Response:   "OK",
						Headers: http.Header{
							contentTypeHeader: {
								textPlainHeader,
							},
						},
					},
					{
						Path:       "/hello",
						StatusCode: 200,
						Response:   "{\"message\": \"hello\"}",
						Headers: http.Header{
							contentTypeHeader: {
								applicationJsonHeader,
							},
							"Server": {
								"nginx",
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:             "ShouldReturnResponsesWithSamePath",
			serverConfigFile: "./testdata/configWithSamePaths.yaml",
			expectedServerConfig: &ServerConfig{
				Host: "localhost",
				Port: 8080,
				Responses: &Responses{
					{
						Path:       "/status",
						StatusCode: 200,
						Response:   "OK",
						Headers: http.Header{
							contentTypeHeader: {
								textPlainHeader,
							},
						},
					},
					{
						Path:       "/status",
						StatusCode: 200,
						Response:   "{\"message\": \"hello\"}",
						Headers: http.Header{
							contentTypeHeader: {
								applicationJsonHeader,
							},
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name:             "ShouldReturnResponsesWithDefaultValues",
			serverConfigFile: "./testdata/configWithDefaultValues.yaml",
			expectedServerConfig: &ServerConfig{
				Host: "0.0.0.0",
				Port: 8080,
				Responses: &Responses{
					{
						Path: "/",
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualServerConfig, actualError := loadServerConfig(testCase.serverConfigFile)
			if testCase.expectedError != actualError {
				t.Errorf("error does not match: expected = %v, actual = %v", testCase.expectedError, actualError)
			}

			fmt.Printf("expected: %v\n", testCase.expectedServerConfig.Responses)
			fmt.Printf("actual: %v\n", actualServerConfig.Responses)

			if !reflect.DeepEqual(testCase.expectedServerConfig, actualServerConfig) {
				t.Errorf("ServerConfig does not match: expected = %v, actual = %v", testCase.expectedServerConfig, actualServerConfig)
			}
		})
	}
}

// reflect.DeepEqual では http.Header (= http.Header) が等しいか否かを正しく判定できなかった
func areHeadersEqual(a, b http.Header) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if w, ok := b[k]; !ok || !reflect.DeepEqual(v, w) {
			return false
		}
	}
	return true
}

// error の == は二つのエラーのポインタが等しいか否かを判定する
// そのため, fmt.Errorf() で作成した二つのエラーはエラーメッセージはひとしくてもポインタとしては異なる
// ゆえに, fmt.Errorf() で作成した二つのエラーのエラーメッセージが等しいか否かは単純な == では判定できない
func areErrorMessagesEqual(a, b error) bool {
	// (0, 0)
	if a == nil && b == nil {
		return true
	}

	// (0, 1)
	if a == nil && b != nil {
		return false
	}

	// (1, 0)
	if a != nil && b == nil {
		return false
	}

	// (1, 1)
	// エラーメッセージが異なる
	if a.Error() != b.Error() {
		return false
	}

	// エラーメッセージが等しい
	return true
}

func TestCreateHandler(t *testing.T) {
	testCases := []struct {
		name          string
		response      *Response
		expectedError error
	}{
		{
			name: "ShouldReturnHandlerThatRespondsResponseInResponseWhenContentTypeHeaderValueIsTextPlainHeader",
			response: &Response{
				Path:       "/status",
				StatusCode: 200,
				Response:   "OK",
				Headers: http.Header{
					contentTypeHeader: {
						textPlainHeader,
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "ShouldReturnHandlerThatRespondsResponseInResponseWhenContentTypeHeaderValueIsApplicationJsonHeader",
			response: &Response{
				Path:       "/hello",
				StatusCode: 200,
				Response:   "{\"message\":\"hello\"}",
				Headers: http.Header{
					contentTypeHeader: {
						applicationJsonHeader,
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "ShouldReturnHandlerThatRespondsResponseInResponseWhenStatusCodeIsNot200",
			response: &Response{
				Path:       "/error",
				StatusCode: 500,
				Response:   "Internal Server Error",
				Headers: http.Header{
					contentTypeHeader: {
						textPlainHeader,
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "ShouldReturnHandlerThatRespondsResponseInResponseWhenResponseHasMultipleHeaders",
			response: &Response{
				Path:       "/multipleHeaders",
				StatusCode: 200,
				Response:   "Multiple Headers",
				Headers: http.Header{
					contentTypeHeader: {
						textPlainHeader,
					},
					"Server": {
						"nginx",
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "ShouldReturnErrorWhenContentTypeHeaderValueIsNotIncludedInAvailableContentTypeHeaderValues",
			response: &Response{
				Path:       "/UnAvailableContentTypeHeaderValue",
				StatusCode: 0,
				Response:   "",
				Headers: http.Header{
					contentTypeHeader: {
						"text/html",
					},
				},
			},
			expectedError: fmt.Errorf("the Content-Type header value is not available: text/html"),
		},
		{
			name: "ShouldReturnErrorWhenContentTypeHeaderValueIsNotIncludedInAvailableContentTypeHeaderValues",
			response: &Response{
				Path:       "/UnAvailableContentTypeHeaderValue",
				StatusCode: 0,
				Response:   "",
				Headers: http.Header{
					contentTypeHeader: {
						"text/html",
					},
				},
			},
			expectedError: fmt.Errorf("the Content-Type header value is not available: text/html"),
		},
	}

	for _, testCase := range testCases {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		_, engine := gin.CreateTestContext(w)
		request := httptest.NewRequest(http.MethodGet, testCase.response.Path, nil)

		handler, actualError := createHandler(testCase.response)
		if !areErrorMessagesEqual(testCase.expectedError, actualError) {
			t.Errorf("error does not match: expected = %v, actual = %v", testCase.expectedError, actualError)
		}

		// expectedError != nil の場合, 以降のテストは実施しない
		if testCase.expectedError != nil {
			continue
		}

		engine.GET(testCase.response.Path, handler)
		engine.ServeHTTP(w, request)

		expectedStatusCode := testCase.response.StatusCode
		actualStatusCode := w.Code
		if expectedStatusCode != actualStatusCode {
			t.Errorf("status code does not match: expected = %v, actual = %v", expectedStatusCode, actualStatusCode)
		}

		expectedResponseBody := testCase.response.Response
		actualResponse := w.Result()
		actualResponseBody, err := io.ReadAll(actualResponse.Body)
		actualResponse.Body.Close()
		if err != nil {
			t.Errorf("the actual response is not read: %s", err)
		}
		if expectedResponseBody != string(actualResponseBody) {
			t.Errorf("response body does not match: expected = %v, actual = %v", expectedResponseBody, string(actualResponseBody))
		}

		expectedHeaders := testCase.response.Headers
		actualHeaders := w.Header()
		fmt.Printf("ex: %v\n", expectedHeaders)
		fmt.Printf("ac: %v\n", actualHeaders)
		if !areHeadersEqual(expectedHeaders, actualHeaders) {
			t.Errorf("headers do not match: expected = %v, actual = %v", expectedHeaders, actualHeaders)
		}
	}
}
