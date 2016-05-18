package server

import (
	"bytes"
	"deadrop/api"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	conf := InitServer()
	cm := conf.upMap

	if n := api.LenChan(cm); n != 0 {
		t.Errorf("Map containing more/less elements than it should: %d elem", api.LenChan(cm))
	}

	for i := 0; i < 10; i++ {
		csHandler := makeHandler(create, conf)
		req, _ := http.NewRequest("GET", "http://localhost:9090/create", nil)
		w := httptest.NewRecorder()
		csHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("[%d] No Stash GET [Create] response: %v", i, http.StatusOK)
		}
	}

	if n := api.LenChan(cm); n != 10 {
		t.Errorf("Map containing more/less elements than it should: %d elem", api.LenChan(cm))
	}
}

// Test does not work yet, as POST is not 100% defined yet
func TestUpload(t *testing.T) {
	conf := InitServer()
	w, req := createGet("http://localhost:9090/upload", conf)
	if w == nil {
		t.Errorf("Error creating GET [Create] request")
	}
	if w.Code != http.StatusOK {
		t.Errorf("[Create] GET response error: %v", w.Code)
	}

	token, err := getToken(w)
	if err != nil {
		t.Errorf("[Upload] Invalid token format")
	}

	csHandler := makeHandler(upload, conf)
	req, err = uploadPost("test1.txt", "http://localhost:9090/upload", token)
	if err != nil {
		t.Errorf("Error creating POST [Upload] request")
	}

	w = httptest.NewRecorder()
	// Wait for parent /create request to complete
	time.Sleep(50 * time.Millisecond)
	csHandler.ServeHTTP(w, req)

	resp_body := w.Body.String()
	fmt.Println(resp_body)

	// TODO: Do some confirmation with resp_body
	if w.Code != http.StatusOK {
		t.Errorf("Response error [Upload]: %v", w.Code)
	}

	ttoken = token
}

func TestFinalize(t *testing.T) {
	conf := InitServer()
	var jsonStr = []byte(`{"Token":"52359c633f1eae96ac7e600a9a4a885b","Lifetime":60,"Files":[{"Fname":"foo.txt","Size":100,"Type":"txt","Download":10},{"Fname":"bar.txt","Size":50,"Type":"txt","Download":5}]}`)
	w, _ := finalizePOST("http://localhost:9090/finalize", jsonStr, conf)

	if w.Code != http.StatusOK {
		t.Errorf("Response error [Finalize]: %v", w.Code)
	}

	stash := api.NewEmptyStash()
	stash.Token = ttoken
	stash.Lifetime = 999
	file := api.NewEmptyStashFile()
	file.Fname = "test1.txt"
	file.Download = 1
	stash.Files = append(stash.Files, file)
	json, _ := json.Marshal(stash)
	jsonStr = []byte(json)
	w, _ = finalizePOST("http://localhost:9090/finalize", jsonStr, conf)

	if w.Code != http.StatusOK {
		t.Errorf("Response error [Finalize]: %v", w.Code)
	}
}

var ttoken string = ""

func TestStashDownload(t *testing.T) {
	if ttoken == "" {
		t.Errorf("[Upload] failure, invalid token")
	}

	conf := InitServer()
	csHandler := makeHandler(download, conf)

	req, err := http.NewRequest("GET", "http://localhost:9090/download/"+ttoken, nil)
	if err != nil {
		t.Errorf("Error creating GET [Download] request")
	}

	w := httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response error [Download]: %v", w.Code)
	}

	// TODO: Fix json response
	// var jsStash stash
	// var body []byte
	// if body, err = ioutil.ReadAll(w.Body); err != nil {
	// 	t.Errorf("[Upload] Invalid token format")
	// }
	// err = json.Unmarshal(body, &jsStash)
	// if err != nil {
	// 	fmt.Println("error:", err)
	// }
	// tfilename = jsStash.Files[0].Fname
	// fmt.Println(tfilename)
}

var tfilename string = ""

func TestFileDownload(t *testing.T) {
	if ttoken == "" {
		t.Errorf("[Upload] failure, invalid token")
	}

	conf := InitServer()
	csHandler := makeHandler(download, conf)

	req, err := http.NewRequest("GET", "http://localhost:9090/download/"+ttoken+"/"+tfilename, nil)
	if err != nil {
		t.Errorf("Error creating GET [Download] request")
	}

	w := httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Response error [Download]: %v", w.Code)
	}

	// TODO: Check if file exists!
}

func createGet(targetUrl string, conf *Configuration) (*httptest.ResponseRecorder, *http.Request) {
	csHandler := makeHandler(create, conf)
	req, err := http.NewRequest("GET", "http://localhost:9090/create", nil)
	if err != nil {
		return nil, nil
	}
	w := httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)

	return w, req
}

func uploadPost(filename string, targetUrl string, token string) (*http.Request, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	err := bodyWriter.WriteField("token", token)
	if err != nil {
		fmt.Println("error creating field")
		return nil, err
	}

	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}

	_, err = io.Copy(fileWriter, f)
	if err != nil {

		return nil, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest("POST", targetUrl, bodyBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

func finalizePOST(targetUrl string, jsonStr []byte, conf *Configuration) (*httptest.ResponseRecorder, *http.Request) {
	csHandler := makeHandler(finalize, conf)
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, nil
	}
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)
	return w, req
}

func getToken(w *httptest.ResponseRecorder) (string, error) {
	type jsonToken struct {
		Token string
	}
	var jsToken jsonToken
	//var body []byte
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &jsToken)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}
	token := jsToken.Token
	fmt.Println(token)
	return token, nil
}
