package si_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-wonk/si/siutils"
)

func BenchmarkHttpHandlerReaderWriterTiny(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestReaderWriterTiny)

	data := stuReqTiny.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}

func BenchmarkHttpHandlerReaderWriterCopiedTiny(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestReaderWriterCopiedTiny)

	data := stuReqTiny.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}
func BenchmarkHttpHandlerBasicTiny(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestBasicTiny)

	data := stuReqTiny.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}

func BenchmarkHttpHandlerReaderWriterSml(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestReaderWriterSml)

	data := stuReqSml.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}

func BenchmarkHttpHandlerReaderWriterCopiedSml(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestReaderWriterCopiedSml)

	data := stuReqSml.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}

func BenchmarkHttpHandlerBasicSml(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestBasicSml)

	data := stuReqSml.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}

func BenchmarkHttpHandlerReaderWriterMed(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestReaderWriterMed)

	data := stuReqMed.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}

func BenchmarkHttpHandlerReaderWriterCopiedMed(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestReaderWriterCopiedMed)

	data := stuReqMed.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}

func BenchmarkHttpHandlerBasicMed(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestBasicMed)

	data := stuReqMed.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}

func BenchmarkHttpHandlerReaderWriterLrg(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestReaderWriterLrg)

	data := stuReqLrg.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}

func BenchmarkHttpHandlerReaderWriterCopiedLrg(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestReaderWriterCopiedLrg)

	data := stuReqLrg.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}

func BenchmarkHttpHandlerBasicLrg(b *testing.B) {
	router := http.NewServeMux()
	router.HandleFunc("/test", handleTestBasicLrg)

	data := stuReqLrg.String()

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer([]byte(data))
		req, err := http.NewRequest("POST", "/test", buf)
		siutils.AssertNilFailB(b, err)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// fmt.Println(rec)
	}
}
