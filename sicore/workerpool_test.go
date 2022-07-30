package sicore_test

import (
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-wonk/si/sicore"
	"github.com/stretchr/testify/assert"
)

var (
	myjobCompleted int64 = 0
)

type myJob struct {
	jobName string
}

func (j *myJob) Execute() (any, error) {
	// log.Println(j.jobName)
	atomic.AddInt64(&myjobCompleted, 1)
	return j.jobName, errors.New("unknown error")
}

type mySlowJob struct {
	jobName string
}

func (j *mySlowJob) Execute() (any, error) {
	// log.Println(j.jobName)
	time.Sleep(500 * time.Millisecond)
	return j.jobName, errors.New("unknown error")
}

func TestWorkerPool_Basic(t *testing.T) {

	log.Println("start")

	p := sicore.NewWorkerPool(5, 128)
	p.Start()

	for i := 0; i < 256; i++ {
		p.Queue() <- &myJob{fmt.Sprintf("hey %d", i)}
	}

	go func() {
		// time.Sleep(1500 * time.Millisecond)

		log.Println("finishing")
		p.Finish()
	}()
	p.Wait()

	assert.EqualValues(t, 256, myjobCompleted)
	log.Println("finished")
}

func TestWorkerPool_Ignore(t *testing.T) {

	log.Println("start")

	p := sicore.NewWorkerPool(5, 128)
	p.Start()

	for i := 0; i < 256; i++ {
		p.QueueIgnore(&mySlowJob{fmt.Sprintf("hey %d", i)})
	}

	go func() {
		// time.Sleep(1500 * time.Millisecond)

		log.Println("finishing")
		p.Finish()
	}()
	p.Wait()

	assert.NotEqualValues(t, 256, myjobCompleted)
	log.Println("finished")
}
func TestWorkerPool_ResultsAndErrors(t *testing.T) {

	log.Println("start")

	p := sicore.NewWorkerPoolWithResultsAndErrors(5, 128)
	p.Start()

	go func() {
		p.ErrorzReady()
		defer p.ErrorzDone()
		for err := range p.Errorz() {
			if err == nil {
				continue
			}
			log.Println(err)
		}
	}()

	go func() {
		p.ResultsReady()
		defer p.ResultsDone()
		for res := range p.Results() {
			if res == nil {
				continue
			}
			log.Println(res)
		}
	}()

	for i := 0; i < 256; i++ {
		p.Queue() <- &myJob{fmt.Sprintf("hey %d", i)}
	}

	go func() {
		// time.Sleep(1500 * time.Millisecond)

		log.Println("finishing")
		p.Finish()
	}()
	p.Wait()

	assert.EqualValues(t, 256, myjobCompleted)
	log.Println("finished")
}
