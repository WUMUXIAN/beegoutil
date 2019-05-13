package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func Get(url string, authHeader string) (int, []byte) {
	code, _, response := GetWithHeaderInResult(url, authHeader)
	return code, response
}

func GetWithHeaderInResult(url string, authHeader string) (int, http.Header, []byte) {
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->\n", "start"))
	logs.Alert(fmt.Sprintf("[HTTP GET: %s]", url))

	r, _ := http.NewRequest("GET", url, nil)
	r.Header.Add("Authorization", authHeader)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	logs.Alert(fmt.Sprintf("[HTTP HEADERS: %s]", r.Header))

	jsonBytes := w.Body.Bytes()

	logs.Alert(fmt.Sprintf("[Status Code --> %d]", w.Code))
	logs.Alert(fmt.Sprintf("[Response Headers --> %v]", w.Header()))
	logs.Alert(fmt.Sprintf("[Response --> %s]", string(jsonBytes)))
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->", "end"))

	return w.Code, w.HeaderMap, jsonBytes
}

func Post(url string, request []byte, authHeader string) (int, []byte) {
	requestBody, _ := json.Marshal(request)
	code, jsonBytes, _ := httpPostRaw(url, requestBody, "text/plain", nil, authHeader)
	return code, jsonBytes
}

func PostWithHeaderResult(url string, request []byte, authHeader string, headers map[string]string) (int, []byte, map[string][]string) {
	requestBody, _ := json.Marshal(request)
	return httpPostRaw(url, requestBody, "text/plain", headers, authHeader)
}

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

	return httpPostRaw(url, body.Bytes(), writer.FormDataContentType(), nil, authHeader)
}

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

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	jsonBytes := w.Body.Bytes()

	logs.Alert(fmt.Sprintf("[Status Code --> %d]", w.Code))
	logs.Alert(fmt.Sprintf("[Response Headers --> %v]", w.HeaderMap))
	logs.Alert(fmt.Sprintf("[Response --> %s]", string(jsonBytes)))
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->", "end"))

	return w.Code, jsonBytes, w.HeaderMap
}

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

	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	jsonBytes := w.Body.Bytes()

	logs.Alert(fmt.Sprintf("[Status Code --> %d]", w.Code))
	logs.Alert(fmt.Sprintf("<------------------------- %s ------------------------->", "end"))

	return w.Code, jsonBytes, w.HeaderMap
}
