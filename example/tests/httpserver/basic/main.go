//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	ip   string
	port int
	key  string
	cert string
	addr string
	dump bool

	filePath string
)

func init() {

	flag.StringVar(&ip, "i", "", "ip")
	flag.IntVar(&port, "p", 8080, "port")
	flag.BoolVar(&dump, "dump", false, "dump request")
	flag.StringVar(&filePath, "filepath", "./uploaded/", "path to upload files to")

	flag.Parse()

	addr = fmt.Sprintf("%v:%v", ip, strconv.Itoa(port))
}
func main() {

	// 라우터, gorilla mux를 쓴다
	router := mux.NewRouter()

	router.HandleFunc("/test/hello", HandleBasic)
	router.HandleFunc("/test/echo", HandleEcho)   //.Methods(http.MethodPost)
	router.HandleFunc("/test/echo2", HandleEcho2) //.Methods(http.MethodPost)
	router.HandleFunc("/test/echo3", HandleEcho3) //.Methods(http.MethodPost)
	router.HandleFunc("/test/echo4", HandleEcho4) //.Methods(http.MethodPost)
	router.HandleFunc("/test/file/upload", UploadFile)

	// http 서버 생성
	httpServer := &http.Server{
		Addr:         addr,             // listen 할 주소(ip:port)
		WriteTimeout: 30 * time.Second, // 서버 > 클라이언트 응답
		ReadTimeout:  30 * time.Second, // 클라이언트 > 서버 요청
		Handler:      router,           // mux다
	}

	log.Fatal(httpServer.ListenAndServe())
}

func HandleBasic(w http.ResponseWriter, req *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpReq))
	}

	_, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// log.Println(string(b))
	w.Write([]byte("hello"))
}

func HandleEcho(w http.ResponseWriter, req *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpReq))
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// log.Println(string(body))
	w.Write(body)
}

func HandleEcho2(w http.ResponseWriter, req *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpReq))
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// log.Println(string(body))
	w.Write(append([]byte("2"), body...))
}
func HandleEcho3(w http.ResponseWriter, req *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpReq))
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// log.Println(string(body))
	w.Write(append([]byte("3"), body...))
}
func HandleEcho4(w http.ResponseWriter, req *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpReq))
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// log.Println(string(body))
	w.Write(append([]byte("4"), body...))
}

// UploadFile 파일을 업로드 하기 위한 핸들러 함수
func UploadFile(w http.ResponseWriter, r *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(r, true)
		fmt.Println(string(dumpReq))
	}
	var err error

	// multipart/form-data로 파싱된 요청본문을 최대 1메가까지 메모리에 저장하도록 한다.
	// r.ParseMultipartForm(1 << 20)
	r.ParseMultipartForm(1 * 1024)

	// FormFile returns the first file for the provided form key.
	// FormFile calls ParseMultipartForm and ParseForm if necessary.
	// 첫번째 파일 데이터와 헤더를 반환한다. ParseMultipartForm과 ParseForm을 호출할 수 있다는데 언제인지는 모르겠다.
	file, header, err := r.FormFile("file_to_upload")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Uploaded File: %+v, File Size: %+v, MIME Header: %+v\n",
		header.Filename, header.Size, header.Header)

	// filePath 디렉토리가 없으면 만들기
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// 경로와 파일명 붙이기
	filePathName := filepath.Join(filePath, header.Filename)

	// 파일 만들기
	f, err := os.Create(filePathName)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer f.Close()

	// 멀티파트 파일 받아서 읽기 위함
	reader := bufio.NewReader(file)

	// 어디까지 읽었는지 보기 위함, 결국엔 사이즈랑 같아야 함
	var offset int64 = 0

	// reader로부터 4096 바이트씩 읽을 것임
	rb := make([]byte, 4096)
	for {
		size, err := reader.Read(rb) // rb에 집어넣기
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		// n, err := f.WriteAt(rb[:size], offset)
		n, err := f.Write(rb[:size])
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		offset += int64(n)
	}
	log.Printf("file size: %v, %v", header.Size, offset)
	w.Write([]byte("success"))
}
