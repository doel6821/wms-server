package helpers

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"wms-server/constants"
	"wms-server/helpers/models"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/kelseyhightower/envconfig"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

// Env ..
type Env struct {
	DebugClient bool   `envconfig:"DEBUG_CLIENT" default:"true"`
	Timeout     string `envconfig:"TIMEOUT" default:"60s"`
	RetryBad    int    `envconfig:"RETRY_BAD" default:"1"`
}

var (
	httpEnv Env
)

// HTTPMethodGet ...
const (
	HTTPMethodGet      = "GET"
	HTTPMethodPost     = "POST"
	HTTPMethodPut      = "PUT"
	HTTPMethodDelete   = "DELETE"
	HTTPMethodPostForm = "POST_FORM"
)

func init() {
	if err := envconfig.Process("HTTP", &httpEnv); err != nil {
		fmt.Println("Failed to get HTTP env:", err)
	}
}

// HTTPGet func
func HTTPGet(url string, header http.Header, to string) (gorequest.Response, []byte, error) {
	request := gorequest.New()
	request.SetDebug(httpEnv.DebugClient)
	timeout, err := time.ParseDuration(to)
	if err != nil {
		timeout, _ = time.ParseDuration(httpEnv.Timeout)
	}
	reqagent := request.Get(url)
	reqagent.Header = header
	resp, body, errs := reqagent.
		Timeout(timeout).
		// Retry(httpEnv.RetryBad, time.Second).
		End()
	if errs != nil {
		return resp, []byte(body), errs[0]
	}
	return resp, []byte(body), nil
}

// HTTPPost func
func HTTPPost(url string, jsondata interface{}, to string) (gorequest.Response, []byte, error) {
	request := gorequest.New()
	request.SetDebug(httpEnv.DebugClient)
	timeout, err := time.ParseDuration(to)
	if err != nil {
		timeout, _ = time.ParseDuration(httpEnv.Timeout)
	}
	if url[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Post(url)
	reqagent.Header.Set("Content-Type", "application/json")
	resp, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		// Retry(httpEnv.RetryBad, time.Second).
		End()
	if errs != nil {
		return resp, []byte(body), errs[0]
	}
	return resp, []byte(body), nil
}

// HTTPPostWithHeader func
func HTTPPostWithHeader(url string, jsondata interface{}, header http.Header, to string) (gorequest.Response, []byte, error) {
	request := gorequest.New()
	request.SetDebug(httpEnv.DebugClient)
	timeout, err := time.ParseDuration(to)
	if err != nil {
		timeout, _ = time.ParseDuration(httpEnv.Timeout)
	}
	reqagent := request.Post(url)
	reqagent.Header = header
	resp, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		// Retry(httpEnv.RetryBad, time.Second).
		End()
	if errs != nil {
		return resp, []byte(body), errs[0]
	}
	return resp, []byte(body), nil
}

// HTTPPostFormWithHeader func
func HTTPPostFormWithHeader(url string, jsondata interface{}, header http.Header, to string) (gorequest.Response, []byte, error) {
	request := gorequest.New()
	request.SetDebug(httpEnv.DebugClient)
	timeout, err := time.ParseDuration(to)
	if err != nil {
		timeout, _ = time.ParseDuration(httpEnv.Timeout)
	}
	reqagent := request.Post(url)
	reqagent.Header = header
	resp, body, errs := reqagent.
		Type("multipart").
		SendFile("client_credentials", "grand_type").
		Timeout(timeout).
		// Retry(httpEnv.RetryBad, time.Second).
		End()
	if errs != nil {
		return resp, []byte(body), errs[0]
	}
	return resp, []byte(body), nil
}

// HTTPPutWithHeader func
func HTTPPutWithHeader(url string, jsondata interface{}, header http.Header, to string) (gorequest.Response, []byte, error) {
	request := gorequest.New()
	request.SetDebug(httpEnv.DebugClient)
	timeout, err := time.ParseDuration(to)
	if err != nil {
		timeout, _ = time.ParseDuration(httpEnv.Timeout)
	}
	if url[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Put(url)
	reqagent.Header = header
	resp, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		// Retry(httpEnv.RetryBad, time.Second).
		End()
	if errs != nil {
		return resp, []byte(body), errs[0]
	}
	return resp, []byte(body), nil
}

// HTTPDeleteWithHeader func
func HTTPDeleteWithHeader(url string, jsondata interface{}, header http.Header, to string) (gorequest.Response, []byte, error) {
	request := gorequest.New()
	request.SetDebug(httpEnv.DebugClient)
	timeout, err := time.ParseDuration(to)
	if err != nil {
		timeout, _ = time.ParseDuration(httpEnv.Timeout)
	}
	if url[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Delete(url)
	reqagent.Header = header
	resp, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		// Retry(httpEnv.RetryBad, time.Second).
		End()
	if errs != nil {
		return resp, []byte(body), errs[0]
	}
	return resp, []byte(body), nil
}

// SendHTTPRequest ..
func SendHTTPRequest(ctx context.Context, method, url string, header http.Header, body interface{}, timeout string, logs *logrus.Logger) (response gorequest.Response, data []byte, err error) {
	start := time.Now()

	header.Add("X-Trace-ID", ctx.Value(constants.TRANSACTION_ID).(string))

	switch method {
	case HTTPMethodGet:
		response, data, err = HTTPGet(url, header, timeout)
	case HTTPMethodPost:
		response, data, err = HTTPPostWithHeader(url, body, header, timeout)
	case HTTPMethodPut:
		response, data, err = HTTPPutWithHeader(url, body, header, timeout)
	case HTTPMethodDelete:
		response, data, err = HTTPDeleteWithHeader(url, body, header, timeout)
	case HTTPMethodPostForm:
		response, data, err = HTTPPostFormWithHeader(url, body, header, timeout)
	}

	var rawMsg json.RawMessage = data
	defer func() {
		hLogs := models.LogModels{
			METHOD:   method,
			PATH:     url,
			HEADER:   header,
			CLIENTIP: "",
			REQUEST:  body,
			RESPONSE: rawMsg,
			DURATION: time.Since(start).Milliseconds(),
		}
		logs.WithContext(ctx).WithFields(structs.Map(hLogs)).Infof(constants.LOG_OTHER)
	}()

	return response, data, err
}

// PostFormDataV2 ..
func PostFormDataV2(uri string, params map[string]string, auth, path, paramName string) (*http.Request, error) {
	var err error
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(paramName, filepath.Base(path))
		if err != nil {
			return nil, err
		}
		io.Copy(part, file)
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", auth)
	return req, err
}

// SendHTTP ...
func SendHTTP(ctx context.Context, log *logrus.Logger, method, URI string, body interface{}, header http.Header, timeout string) (response []byte, err error) {

	start := time.Now()
	header["X-Trace-ID"] = []string{ctx.Value(constants.TRANSACTION_ID).(string)}
	var req *http.Request

	to, err := time.ParseDuration(timeout)
	if err != nil {
		log.WithContext(ctx).Error(Errotype{}.Parse("duration", err))
		to = time.Second * 60
	}

	client := &http.Client{
		Timeout: to,
	}

	switch method {
	case "GET":
		req, err = sendHTTPGET(ctx, log, URI, header)
	case "POST":
		req, err = sendHTTPPOST(ctx, log, URI, body, header)
	case "FORM":
		req, err = sendHTTPPOSTFORM(ctx, log, URI, body, header)
	default:
		req, err = sendHTTPDefault(ctx, log, method, URI, body, header)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.WithContext(ctx).Error(Errotype{}.GetData("response", err))
		return nil, fmt.Errorf("request failed: %v", err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithContext(ctx).Error(Errotype{}.Read("response", err))
		return nil, fmt.Errorf("request failed: %v", err)
	}

	reqBody, _ := json.Marshal(body)
	headerByte, _ := json.Marshal(header)
	var rawMsg json.RawMessage = resBody

	defer func() {
		hLogs := models.LogModels{
			METHOD:   method,
			PATH:     URI,
			HEADER:   string(headerByte),
			CLIENTIP: "",
			REQUEST:  string(reqBody),
			RESPONSE: rawMsg,
			DURATION: time.Since(start).Milliseconds(),
		}
		log.WithContext(ctx).WithFields(structs.Map(hLogs)).Infof(constants.LOG_OTHER)
		resp.Body.Close()
	}()
	return resBody, nil
}

func sendHTTPGET(ctx context.Context, log *logrus.Logger, URI string, header http.Header) (res *http.Request, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		log.WithContext(ctx).Error(Errotype{}.SendData(URI, err))
		return nil, fmt.Errorf("%v", Errotype{}.SendData("request", err))
	}
	req.Header = header
	return req, err
}

func sendHTTPDefault(ctx context.Context, log *logrus.Logger, method, URI string, body interface{}, header http.Header) (res *http.Request, err error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("%v", Errotype{}.Marshal("body", err))
	}
	bodyData := bytes.NewBuffer(bodyJSON)
	req, err := http.NewRequestWithContext(ctx, method, URI, bodyData)
	if err != nil {
		log.WithContext(ctx).Error(Errotype{}.SendData(URI, err))
		return nil, fmt.Errorf("%v", Errotype{}.CreateData("request", err))
	}
	req.Header = header
	return req, err
}

func sendHTTPPOST(ctx context.Context, log *logrus.Logger, URI string, body interface{}, header http.Header) (res *http.Request, err error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("%v", Errotype{}.Marshal("body", err))
	}
	bodyData := bytes.NewBuffer(bodyJSON)
	req, err := http.NewRequestWithContext(ctx, "POST", URI, bodyData)

	if err != nil {
		log.WithContext(ctx).Error(Errotype{}.SendData(URI, err))
		return nil, fmt.Errorf("%v", Errotype{}.SendData("request", err))
	}
	req.Header = header
	req.Header.Set(constants.HEAD_CONTENT_TYPE, "application/json")
	return req, err
}

func sendHTTPPOSTFORM(ctx context.Context, log *logrus.Logger, URI string, body interface{}, header http.Header) (res *http.Request, err error) {

	data, _ := convertToFormData(body)
	bodyData := strings.NewReader(data.Encode())
	req, err := http.NewRequestWithContext(ctx, "POST", URI, bodyData)

	if err != nil {
		log.WithContext(ctx).Error(Errotype{}.SendData(URI, err))
		return nil, fmt.Errorf("%v", Errotype{}.SendData("request", err))
	}
	req.Header = header
	req.Header.Set(constants.HEAD_CONTENT_TYPE, "application/x-www-form-urlencoded")
	return req, err
}

// convertToFormData ...
func convertToFormData(data interface{}) (url.Values, error) {
	formData := url.Values{}

	// Memeriksa tipe data
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Map {
		for _, key := range val.MapKeys() {
			keyStr := fmt.Sprintf("%v", key)
			value := val.MapIndex(key).Interface()
			// Menyusun form data
			formData.Set(keyStr, fmt.Sprintf("%v", value))
		}
		return formData, nil
	}
	return nil, fmt.Errorf("unsupported data type")
}
