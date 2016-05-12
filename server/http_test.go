package server

import (
	"bytes"
	"deadrop/api"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	"regexp"
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
			t.Errorf("[%d] No Stash GET response: %v", i, http.StatusOK)
		}
	}

	if n := api.LenChan(cm); n != 10 {
		t.Errorf("Map containing more/less elements than it should: %d elem", api.LenChan(cm))
	}
}

// Test does not work yet, as POST is not 100% defined yet
func TestUpload(t *testing.T) {
	conf := InitServer()

	csHandler := makeHandler(create, conf)
	req, err := http.NewRequest("GET", "http://localhost:9090/create", nil)
	if err != nil {
		t.Errorf("Error creating GET /create request")
	}
	w := httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("[Create] GET response: %v", w.Code)
	}
	body := w.Body.String()

	// Removing '"' characters from string
	reg := regexp.MustCompile(`"([^"]*)"`)
	res := reg.ReplaceAllString(body, "${1}")
	token := string(res)
	fmt.Println(token)
	
	csHandler = makeHandler(upload, conf)
	req, err = postFile("test1.txt", "http://localhost:9090/upload", token)
	if err != nil {
		t.Errorf("Error creating POST /upload request")
	}

	w = httptest.NewRecorder()
	// Wait for parent /create request to complete
	time.Sleep(50 * time.Millisecond)
	csHandler.ServeHTTP(w, req)

	resp_body := w.Body.String()
	fmt.Println(resp_body)

	// TODO: Do some confirmation with resp_body
}


func TestFinalize(t *testing.T) {
	conf := InitServer()

	csHandler := makeHandler(finalize, conf)

	var jsonStr = []byte(`{"Token":"hfsiehfsiehf983989wrhiuhsi","Lifetime":60,"Files":[{"Fname":"blaj.txt","Size":100,"Type":"txt","Download":10},{"Fname":"blaj.txt","Size":100,"Type":"txt","Download":10}]}`)

	req, _ := http.NewRequest("POST", "http://localhost:9090/finalize", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	csHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("No Stash POST response: %v", http.StatusOK)
	}
}

func TestStashDownload(t *testing.T) {

}

func TestFileDownload(t *testing.T) {

}

func postFile(filename string, targetUrl string, token string) (*http.Request, error) {
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
