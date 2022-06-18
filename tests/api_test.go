package si_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-wonk/si"
	"github.com/go-wonk/si/siutils"
	"github.com/go-wonk/si/tests/testmodels"
)

func makeStuReq(num int) testmodels.StudentList {
	res := make(testmodels.StudentList, 0, num)
	for i := 0; i < num; i++ {
		res = append(res, testmodels.Student{ID: i})
	}
	return res
}
func makeStuRes(num int) testmodels.StudentList {
	res := make(testmodels.StudentList, 0, num)
	for i := 0; i < num; i++ {
		res = append(res, testmodels.Student{ID: i})
	}
	return res
}

var (
	stuReqSml = makeStuReq(8)
	stuResSml = makeStuRes(8)
	stuReqMed = makeStuReq(128)
	stuResMed = makeStuRes(128)
	stuReqLrg = makeStuReq(512)
	stuResLrg = makeStuRes(512)
)

func handleTestBasicSml(w http.ResponseWriter, r *http.Request) {
	var req testmodels.StudentList
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(b, &req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(&stuResSml); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
func handleTestReaderWriterSml(w http.ResponseWriter, r *http.Request) {
	var req testmodels.StudentList
	if err := si.DecodeJson(&req, r.Body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := si.EncodeJson(w, &stuResSml); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func handleTestBasicMed(w http.ResponseWriter, r *http.Request) {
	var req testmodels.StudentList
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(b, &req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(&stuResMed); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
func handleTestReaderWriterMed(w http.ResponseWriter, r *http.Request) {
	var req testmodels.StudentList
	if err := si.DecodeJson(&req, r.Body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := si.EncodeJson(w, &stuResMed); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func handleTestBasicLrg(w http.ResponseWriter, r *http.Request) {
	var req testmodels.StudentList
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(b, &req); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(&stuResLrg); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
func handleTestReaderWriterLrg(w http.ResponseWriter, r *http.Request) {
	var req testmodels.StudentList
	if err := si.DecodeJson(&req, r.Body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := si.EncodeJson(w, &stuResLrg); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}
func TestHttpHandlerSml(t *testing.T) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestReaderWriterSml)

	buf := bytes.NewBuffer([]byte(`[{"id":1,"email_address":"wonk@wonk.org","name":"wonk","borrowed":false,"book_id":23}]`))
	req, err := http.NewRequest("POST", "/test", buf)
	siutils.AssertNilFail(t, err)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	fmt.Println(rec)
}
