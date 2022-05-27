package sicore_test

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
)

var (
	// onlinetest = os.Getenv("ONLINE_TEST")
	onlinetest = "1"

	listener       net.Listener
	listenerClosed bool
)

func startTcpServer(waitChannel chan bool) error {
	//통신 방식과 포트값을 전달해 리스너 객체 생성
	var err error
	listener, err = net.Listen("tcp", ":10000")

	//예외처리
	if err != nil {
		log.Printf("fail to bind address; err: %v", err)
		return err
	}
	// defer listener.Close()

	listenerClosed = false
	// log.Printf("## 프로그램 시작")
	waitChannel <- true
	//메인 루프
	for {
		//연결 대기
		connection, err := listener.Accept()
		//연결 실패
		if err != nil {
			if listenerClosed {
				break
			}
			log.Printf("Accept failed: %v", err)
			continue
		}
		// log.Printf("client connected: %v", connection.RemoteAddr())

		//각 연결에 대한 처리를 고루틴으로 실행
		go func() {
			buffer := make([]byte, 1000) //버퍼

			totalReceived := 0
			received := make([]byte, 0, 1000001)
			//다 받을때까지 반복하며 읽음
			for {
				//입력
				count, err := connection.Read(buffer)
				if nil != err {
					//입력이 종료되면
					if io.EOF == err {
						// log.Printf("연결 종료: %v", connection.RemoteAddr().String())
					} else {
						log.Printf("수신 실패: %v", err)
					}
					return
				}

				totalReceived += count
				if count > 0 {
					//받아온 길이만큼 슬라이스를 잘라서 출력
					data := buffer[:count]
					// log.Println(string(data[:14]))
					received = append(received, data...)
				}

				lenStr := string(received[:7])
				lenProt, _ := strconv.ParseInt(lenStr, 10, 64)
				if int(lenProt) == totalReceived {
					// log.Println("writing...")
					connection.Write(received)
				}
			}
		}()
	}
	return nil
}

func setup() error {
	// waitChannel := make(chan bool)
	// go startTcpServer(waitChannel)

	// for range waitChannel {
	// 	break
	// }
	if onlinetest == "1" {
	}

	return nil
}

func shutdown() {
	if listener != nil {
		listenerClosed = true
		listener.Close()
	}
}

func TestMain(m *testing.M) {

	err := setup()
	if err != nil {
		fmt.Println(err)
		shutdown()
		os.Exit(1)
	}

	exitCode := m.Run()

	shutdown()
	os.Exit(exitCode)
}
