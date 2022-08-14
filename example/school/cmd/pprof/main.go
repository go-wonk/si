package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/go-wonk/si/example/school/adaptor"
	"github.com/go-wonk/si/example/school/core"
	"github.com/go-wonk/si/sicore"
	"github.com/go-wonk/si/sihttp"
	_ "github.com/lib/pq"
)

var (
	ip   string
	port int
	addr string
	dump bool

	studentUsc   core.StudentUsecase
	borrowingUsc core.BorrowingUsecase
	bookUsc      core.BookUsecase

	defaultClient *http.Client
	client        *sihttp.Client
)

func init() {

	flag.StringVar(&ip, "i", "", "ip")
	flag.IntVar(&port, "p", 8082, "port")
	flag.BoolVar(&dump, "dump", false, "dump request")

	flag.Parse()

	addr = fmt.Sprintf("%v:%v", ip, strconv.Itoa(port))

}
func main() {

	connStr := "host=testpghost port=5432 user=test password=test123 dbname=testdb sslmode=disable connect_timeout=60"
	driver := "postgres"
	db, err := sql.Open(driver, connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	defaultClient = sihttp.DefaultInsecureClient()
	client = sihttp.NewClient(defaultClient)

	txBeginner := adaptor.NewTxBeginner(db)
	studentRepo := adaptor.NewPgStudentRepo(db)
	bookRepo := adaptor.NewPgBookRepo(db)
	borrowingRepo := adaptor.NewPgBorrowingRepo(db)

	studentUsc = core.NewStudentUsecaseImpl(txBeginner, studentRepo)
	borrowingUsc = core.NewBorrowingUsecaseImpl(txBeginner, borrowingRepo, studentRepo, bookRepo)
	bookUsc = core.NewBookUsecaseImpl(txBeginner, bookRepo)

	// 라우터, gorilla mux를 쓴다
	router := http.DefaultServeMux

	router.HandleFunc("/test/hello", HandleBasic)
	router.HandleFunc("/test/echo", HandleEcho)
	router.HandleFunc("/test/pprof", HandlePprof)
	router.HandleFunc("/test/gc", HandleGC)
	router.HandleFunc("/test/findall", HandleFindAllStudent)
	// router.HandleFunc("/test/repeat/findall", HandleRepeatFindAll)
	router.HandleFunc("/test/sendfiles", HandlerSendFile)
	router.HandleFunc("/test/repeat/sendfiles", HandlerRepeatSendFile)

	// http 서버 생성
	httpServer := &http.Server{
		Addr:         addr,             // listen 할 주소(ip:port)
		WriteTimeout: 30 * time.Second, // 서버 > 클라이언트 응답
		ReadTimeout:  30 * time.Second, // 클라이언트 > 서버 요청
		Handler:      router,           // mux다
	}

	go func() {
		log.Println(http.ListenAndServe(":56060", nil))
	}()

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

func HandleGC(w http.ResponseWriter, req *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpReq))
	}

	body := sicore.GetReader(req.Body)
	defer sicore.PutReader(body)

	_, err := body.ReadAll()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	runtime.GC()
	w.Write([]byte("manual gc done"))
}

func HandlePprof(w http.ResponseWriter, req *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpReq))
	}

	body := sicore.GetReader(req.Body, sicore.SetJsonDecoder())
	defer sicore.PutReader(body)

	var s core.Student
	err := body.Decode(&s)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = studentUsc.Add(s.EmailAddress, s.Name)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	list, err := studentUsc.FindAll()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	res := sicore.GetWriter(w, sicore.SetJsonEncoder())
	defer sicore.PutWriter(res)

	err = res.EncodeFlush(list)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func HandleFindAllStudent(w http.ResponseWriter, req *http.Request) {
	if dump {
		dumpReq, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dumpReq))
	}

	// read request body
	body := sicore.GetReader(req.Body, sicore.SetJsonDecoder())
	defer sicore.PutReader(body)

	var s core.Student
	err := body.Decode(&s)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//////////////////////////////////////////////////////////////////////////////////////

	// find all students
	list, err := studentUsc.FindAll()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//////////////////////////////////////////////////////////////////////////////////////

	// write to file
	f, err := os.OpenFile("./students.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer f.Close()
	fw := sicore.GetWriter(f)
	defer sicore.PutWriter(fw)
	for _, student := range list {
		_, err := fw.Write([]byte(student.EmailAddress + "," + student.Name + "\n"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
	if err = fw.Flush(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//////////////////////////////////////////////////////////////////////////////////////

	// write to client
	res := sicore.GetWriter(w, sicore.SetJsonEncoder())
	defer sicore.PutWriter(res)

	err = res.EncodeFlush(list)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

var (
	fileNames = []string{"account_black.png", "account-lock_black.png", "disconnected.png", "down_arrow.png", "down2_arrow.png", "expand_down.png", "expand_up.png"}
)

func HandlerSendFile(w http.ResponseWriter, req *http.Request) {
	for i, name := range fileNames {
		params := make(map[string]string)
		params["id"] = strconv.Itoa(i)
		_, err := client.RequestPostFile("http://127.0.0.1:8080/test/file/upload", nil, params, "file_to_upload", "./data/"+name)
		if err != nil {
			fmt.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("success"))
}

func HandlerRepeatSendFile(w http.ResponseWriter, req *http.Request) {
	for j := 0; j < 5000; j++ {

		for i, name := range fileNames {
			params := make(map[string]string)
			params["id"] = strconv.Itoa(i)
			_, err := client.RequestPostFile("http://127.0.0.1:8080/test/file/upload", nil, params, "file_to_upload", "./data/"+name)
			if err != nil {
				fmt.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		if j%100 == 0 {
			fmt.Println("gcing...")
			runtime.GC()
		}
		time.Sleep(100 * time.Millisecond)
	}

	w.Write([]byte("success"))
}
