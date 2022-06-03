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
	"runtime"
	"strconv"
	"time"

	"github.com/go-wonk/si/example/school/adaptor"
	"github.com/go-wonk/si/example/school/core"
	"github.com/go-wonk/si/sicore"
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
)

func init() {

	flag.StringVar(&ip, "i", "", "ip")
	flag.IntVar(&port, "p", 8080, "port")
	flag.BoolVar(&dump, "dump", false, "dump request")

	flag.Parse()

	addr = fmt.Sprintf("%v:%v", ip, strconv.Itoa(port))

}
func main() {

	connStr := "host=172.16.130.144 port=5432 user=test password=test123 dbname=testdb sslmode=disable connect_timeout=60"
	driver := "postgres"
	db, err := sql.Open(driver, connStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

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

	// http 서버 생성
	httpServer := &http.Server{
		Addr:         addr,             // listen 할 주소(ip:port)
		WriteTimeout: 30 * time.Second, // 서버 > 클라이언트 응답
		ReadTimeout:  30 * time.Second, // 클라이언트 > 서버 요청
		Handler:      router,           // mux다
	}

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	log.Fatal(httpServer.ListenAndServe())
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
