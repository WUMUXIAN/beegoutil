package httpclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/astaxie/beego/logs"
)

var (
	transport http.RoundTripper
	timeout   time.Duration
)

func init() {
	// disable the CA verification
	tlsConfig := &tls.Config{InsecureSkipVerify: true}
	transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	timeout = time.Second * 10
}

// Get sends a http get request to the url with the auth header.
func Get(url string, authHeader string) (int, []byte) {
	code, _, response := GetWithHeaderInResult(url, authHeader)
	return code, response
}

// GetWithHeaderInResult sends a http get request and return the response together with headers.
func GetWithHeaderInResult(url string, authHeader string) (int, http.Header, []byte) {
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->\n", "start"))
	logs.Alert(fmt.Sprintf("[HTTP GET: %s]", url))

	r, _ := http.NewRequest("GET", url, nil)
	r.Header.Add("Authorization", authHeader)
	logs.Alert(fmt.Sprintf("[HTTP HEADERS: %s]", r.Header))

	client := &http.Client{Transport: transport, Timeout: timeout}
	resp, err := client.Do(r)

	if err != nil {
		logs.Alert(fmt.Sprintf("[Request Error --> %v]", err))
		return http.StatusInternalServerError, nil, nil
	}

	defer resp.Body.Close()
	jsonBytes, _ := ioutil.ReadAll(resp.Body)
	logs.Alert(fmt.Sprintf("[Status Code --> %d]", resp.StatusCode))
	logs.Alert(fmt.Sprintf("[Response Headers --> %v]", resp.Header))
	logs.Alert(fmt.Sprintf("[Response --> %s]", string(jsonBytes)))
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->", "end"))
	return resp.StatusCode, resp.Header, jsonBytes
}

// Post sends a http post request to the url with post body and auth header.
func Post(url string, requestBody []byte, authHeader string) (int, []byte) {
	code, jsonBytes, _ := PostRaw(url, requestBody, "text/plain", nil, authHeader)
	return code, jsonBytes
}

// PostWithHeaderResult sends a http post request to the url with post body and auth header and get response together with response header.
func PostWithHeaderResult(url string, requestBody []byte, authHeader string, headers map[string]string) (int, []byte, map[string][]string) {
	return PostRaw(url, requestBody, "text/plain", headers, authHeader)
}

// PostMultipart sends a multi-part form post request to the url.
func PostMultipart(url string, files map[string][]byte, fields map[string][]byte, authHeader string) (int, []byte, map[string][]string) {
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->\n", "start"))
	logs.Alert(fmt.Sprintf("[HTTP POST Multipart: %s]", url))
	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	for fieldName, field := range fields {
		part, err := writer.CreateFormField(fieldName)
		if err != nil {
			fmt.Println(err.Error())
		}
		part.Write(field)
	}
	for fileName, fileContent := range files {
		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			fmt.Println(err.Error())
		}
		part.Write(fileContent)
	}

	writer.Close()

	// logs.Alert(fmt.Sprintln("Multi-Part Request Body:", string(body.Bytes())))

	return PostRaw(url, body.Bytes(), writer.FormDataContentType(), nil, authHeader)
}

// PostRaw posts a raw byte array to server.
func PostRaw(url string, requestBody []byte, contentType string, requestHeaders map[string]string, authHeader string) (int, []byte, map[string][]string) {
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->\n", "start"))
	logs.Alert(fmt.Sprintf("[HTTP POST: %s]", url))

	r, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	r.Header.Set("Content-Type", contentType)

	if authHeader != "" {
		r.Header.Add("Authorization", authHeader)
	}

	for name, value := range requestHeaders {
		r.Header.Add(name, value)
	}

	logs.Alert(fmt.Sprintf("[HTTP HEADERS: %s]", r.Header))

	client := &http.Client{Transport: transport, Timeout: timeout}
	resp, err := client.Do(r)

	if err != nil {
		logs.Alert(fmt.Sprintf("[Request Error --> %v]", err))
		return http.StatusInternalServerError, nil, nil
	}

	defer resp.Body.Close()
	jsonBytes, _ := ioutil.ReadAll(resp.Body)
	logs.Alert(fmt.Sprintf("[Status Code --> %d]", resp.StatusCode))
	logs.Alert(fmt.Sprintf("[Response Headers --> %v]", resp.Header))
	logs.Alert(fmt.Sprintf("[Response --> %s]", string(jsonBytes)))
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->", "end"))
	return resp.StatusCode, jsonBytes, resp.Header
}

// Put sends a put request to the url
func Put(url string, requestBody []byte, requestHeaders map[string]string, authHeader string) (int, []byte, map[string][]string) {
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->\n", "start"))
	logs.Alert(fmt.Sprintf("[HTTP PUT: %s]", url))

	r, _ := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))

	if authHeader != "" {
		r.Header.Add("Authorization", authHeader)
	}

	for name, value := range requestHeaders {
		r.Header.Add(name, value)
	}

	client := &http.Client{Transport: transport, Timeout: timeout}
	resp, err := client.Do(r)

	if err != nil {
		logs.Alert(fmt.Sprintf("[Request Error --> %v]", err))
		return http.StatusInternalServerError, nil, nil
	}

	defer resp.Body.Close()
	jsonBytes, _ := ioutil.ReadAll(resp.Body)
	logs.Alert(fmt.Sprintf("[Status Code --> %d]", resp.StatusCode))
	logs.Alert(fmt.Sprintf("[Response Headers --> %v]", resp.Header))
	logs.Alert(fmt.Sprintf("[Response --> %s]", string(jsonBytes)))
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->", "end"))
	return resp.StatusCode, jsonBytes, resp.Header
}

// Delete sends a delete request to the given url
func Delete(url string, authHeader string) (int, []byte) {
	code, _, response := DeleteWithHeaderInResult(url, authHeader)
	return code, response
}

// DeleteWithHeaderInResult sends a http delete request and return the response together with headers.
func DeleteWithHeaderInResult(url string, authHeader string) (int, http.Header, []byte) {
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->\n", "start"))
	logs.Alert(fmt.Sprintf("[HTTP DELETE: %s]", url))

	r, _ := http.NewRequest("DELETE", url, nil)
	r.Header.Add("Authorization", authHeader)
	logs.Alert(fmt.Sprintf("[HTTP HEADERS: %s]", r.Header))

	client := &http.Client{Transport: transport, Timeout: timeout}
	resp, err := client.Do(r)

	if err != nil {
		logs.Alert(fmt.Sprintf("[Request Error --> %v]", err))
		return http.StatusInternalServerError, nil, nil
	}

	defer resp.Body.Close()
	jsonBytes, _ := ioutil.ReadAll(resp.Body)
	logs.Alert(fmt.Sprintf("[Status Code --> %d]", resp.StatusCode))
	logs.Alert(fmt.Sprintf("[Response Headers --> %v]", resp.Header))
	logs.Alert(fmt.Sprintf("[Response --> %s]", string(jsonBytes)))
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->", "end"))
	return resp.StatusCode, resp.Header, jsonBytes
}

// Patch sends a patch request to the url
func Patch(url string, requestBody []byte, requestHeaders map[string]string, authHeader string) (int, []byte, map[string][]string) {
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->\n", "start"))
	logs.Alert(fmt.Sprintf("[HTTP PATCH: %s]", url))

	r, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))

	if authHeader != "" {
		r.Header.Add("Authorization", authHeader)
	}

	for name, value := range requestHeaders {
		r.Header.Add(name, value)
	}

	client := &http.Client{Transport: transport, Timeout: timeout}
	resp, err := client.Do(r)

	if err != nil {
		logs.Alert(fmt.Sprintf("[Request Error --> %v]", err))
		return http.StatusInternalServerError, nil, nil
	}

	defer resp.Body.Close()
	jsonBytes, _ := ioutil.ReadAll(resp.Body)
	logs.Alert(fmt.Sprintf("[Status Code --> %d]", resp.StatusCode))
	logs.Alert(fmt.Sprintf("[Response Headers --> %v]", resp.Header))
	logs.Alert(fmt.Sprintf("[Response --> %s]", string(jsonBytes)))
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->", "end"))
	return resp.StatusCode, jsonBytes, resp.Header
}
