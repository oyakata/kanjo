package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONWOrdCountHandler(t *testing.T) {
	v := "今川義元1560年桶狭間"
	location := fmt.Sprintf("http://example.com/count?format=json&text=%v", v)

	w := httptest.NewRecorder()
	// httptestにはなぜかNewRequestはないと言われる(go version 1.6.2)
	// undefined: httptest.NewRequest
	r, _ := http.NewRequest("GET", location, nil)

	JSONWordCountHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("JSONWordCountHandler=%v, got=%v", http.StatusOK, w.Code)
	}
	result := &WordCount{}
	json.Unmarshal([]byte(w.Body.String()), result)
	expected := WordCount{v, 12, 28}
	if *result != expected {
		t.Errorf("JSONWordCountHandler=%v, got=%v", expected, result)
	}
}
