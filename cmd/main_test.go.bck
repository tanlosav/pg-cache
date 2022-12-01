package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/tanlosav/pg-cache/internal/cache/pgcache"
)

const user = "postgres"
const password = "password"
const name = "cache"
const host = "srv8-noteburnaw"

func TestHandler(t *testing.T) {
	var key = "1"
	var value1 = "value 1"
	var value2 = "value 2"
	var req *http.Request
	var recorder *httptest.ResponseRecorder
	var responseBody string

	cache := pgcache.NewCache(user, password, name, host)
	cache.Connect()

	handler := CacheHandler{
		Cache: cache,
	}

	router := httprouter.New()
	router.GET("/documents/:id", handler.get)
	router.POST("/documents/:id", handler.create)
	router.PUT("/documents/:id", handler.update)
	router.DELETE("/documents/:id", handler.delete)
	router.DELETE("/documents", handler.clean)

	// clean
	req, _ = http.NewRequest("DELETE", "/documents", nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusNoContent {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}

	// get none existent document
	req, _ = http.NewRequest("GET", "/documents/"+key, nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusNotFound {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}

	// create document
	req, _ = http.NewRequest("POST", "/documents/"+key, bytes.NewBuffer([]byte(`{"value":"`+value1+`"}`)))
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusCreated {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}

	// get document
	req, _ = http.NewRequest("GET", "/documents/"+key, nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}
	responseBody = recorder.Body.String()
	if responseBody != `{"value": "`+value1+`"}` {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong document: %q", responseBody)
	}

	// duplicate document
	req, _ = http.NewRequest("POST", "/documents/"+key, bytes.NewBuffer([]byte(`{"value":"`+value1+`"}`)))
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusConflict {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}

	// update document
	req, _ = http.NewRequest("PUT", "/documents/"+key, bytes.NewBuffer([]byte(`{"value":"`+value2+`"}`)))
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusNoContent {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}

	// get document
	req, _ = http.NewRequest("GET", "/documents/"+key, nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}
	responseBody = recorder.Body.String()
	if responseBody != `{"value": "`+value2+`"}` {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong document: %q", responseBody)
	}

	// delete document
	req, _ = http.NewRequest("DELETE", "/documents/"+key, nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusNoContent {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}

	// get document
	req, _ = http.NewRequest("GET", "/documents/"+key, nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusNotFound {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}

	// delete none existent document
	req, _ = http.NewRequest("DELETE", "/documents/"+key, nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusNoContent {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}

	// update none existent document
	req, _ = http.NewRequest("PUT", "/documents/"+key, bytes.NewBuffer([]byte(`{"value":"`+value1+`"}`)))
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusNoContent {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}

	// get document
	req, _ = http.NewRequest("GET", "/documents/"+key, nil)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong status: %d", recorder.Code)
	}
	responseBody = recorder.Body.String()
	if responseBody != `{"value": "`+value1+`"}` {
		fmt.Println(recorder.Body)
		t.Errorf("Wrong document: %q", responseBody)
	}
}
