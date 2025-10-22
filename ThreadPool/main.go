package main

import (
	"log"
	"net"
	"time"
)

type Job struct {
	conn net.Conn
}

type Worker struct {
	id      int
	jobChan chan Job
}

type Pool struct {
	JobQueue chan Job
	Workers  []*Worker
}

func CreateWorker(id int, jobChan chan Job) *Worker {
	return &Worker{
		id:      id,
		jobChan: jobChan,
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobChan {
			log.Printf("Worker %d is handling job from %s", w.id, job.conn.RemoteAddr())
			handleConnection(job.conn)
		}
	}()
}

func CreatePool(numberOfWorkers int) *Pool {
	return &Pool{
		JobQueue: make(chan Job),
		Workers:  make([]*Worker, numberOfWorkers),
	}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.Workers); i++ {
		worker := CreateWorker(i, p.JobQueue)
		p.Workers[i] = worker
		p.Workers[i].Start()
	}
}
func (p *Pool) AddJob(conn net.Conn) {
	p.JobQueue <- Job{conn: conn}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1000)
	conn.Read(buf)
	time.Sleep(1 * time.Second)
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\ntest\r\n"))
}

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// 1 pool with 2 threads
	pool := CreatePool(2)
	pool.Start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		//go handleConnection(conn)
		pool.AddJob(conn)
	}
}
