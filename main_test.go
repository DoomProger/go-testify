package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerWhenStatusOKBodyNotEmpty(t *testing.T) {

	reqCityCoount := 4
	reqCity := "moscow"
	reqUri := fmt.Sprintf("/cafe?count=%d&city=%s", reqCityCoount, reqCity)

	req := httptest.NewRequest("GET", reqUri, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	expectedStatusCode := http.StatusOK

	require.Equal(t, expectedStatusCode, responseRecorder.Code)
	assert.NotEmpty(t, responseRecorder.Body.String())

}
func TestMainHandlerWhenWrongCity(t *testing.T) {

	reqCityCoount := 10
	reqCity := "moscow1"
	reqUri := fmt.Sprintf("/cafe?count=%d&city=%s", reqCityCoount, reqCity)

	req := httptest.NewRequest("GET", reqUri, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	expectedStatusCode := http.StatusBadRequest
	badRequestWrongCity := `wrong city value`

	assert.Equal(t, expectedStatusCode, responseRecorder.Code)
	assert.Contains(t, badRequestWrongCity, responseRecorder.Body.String())

}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	reqCityCoount := 10
	reqCity := "moscow"
	reqUri := fmt.Sprintf("/cafe?count=%d&city=%s", reqCityCoount, reqCity)

	req := httptest.NewRequest("GET", reqUri, nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	responseBody := responseRecorder.Body.String()
	gotCountCafe := len(strings.Split(responseBody, ","))
	if gotCountCafe != totalCount {
		t.Errorf("expected cafe count: %d, got %d", totalCount, gotCountCafe)
	}

}
