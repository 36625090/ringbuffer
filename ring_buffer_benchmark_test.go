package ringbuffer

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func BenchmarkRingBuffer_Sync(b *testing.B) {
	rb := New(1024)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Write(data)
		rb.Read(buf)
	}
}

func TestRingBuffer_AsyncRead(b *testing.T) {
	rb := New(16)
	done := make(chan byte, 1)
	var betters []byte
	for i := 'A'; i <= 'Z'; i++ {
		betters = append(betters, byte(i))
	}

	go func() {
		buf := make([]byte, 2)
		var data []byte

		for len(done) == 0 {
			n, err := rb.Read(buf)
			if err == nil && n == 2 {
				data = append(data, buf...)
			}
			time.Sleep(time.Millisecond * 10)
		}
		b.Log("r1", len(data))
	}()

	go func() {
		buf := make([]byte, 2)
		var data []byte
		for len(done) == 0 {
			if n, err := rb.Read(buf); err == nil && n == 2 {
				data = append(data, buf...)
			}
			time.Sleep(time.Millisecond * 10)
		}
		b.Log("r2", len(data))
	}()

	//b.ResetTimer()
	t := 0
	for i := 0; i < 900; i++ {
		rb.BlockingWrite([]byte{betters[i%26]})
		t++
	}
	time.Sleep(time.Second * 10)
	b.Log("write", t, rb.buf)
	done <- 0
	time.Sleep(time.Second * 10)
}

func BenchmarkRingBuffer_AsyncWrite(b *testing.B) {
	rb := New(os.Getpagesize() / 2)
	data := []byte(strings.Repeat("a", 512))
	buf := make([]byte, 512)

	go func() {
		for {
			rb.Write(data)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Read(buf)
	}
}
func BenchmarkRingBuffer_AsyncWrite2(b *testing.B) {
	rb := New(1024)
	done := make(chan byte, 1)
	var betters [][]byte
	for i := 'A'; i <= 'Z'; i++ {
		betters = append(betters, []byte{byte(i)})
	}
	//b.StartTimer()
	//defer b.StopTimer()
	for i := 0; i < 1; i++ {
		go func() {
			buf := make([]byte, 2)
			//var data []byte

			for len(done) == 0 {
				rb.Read(buf)
				//time.Sleep(time.Millisecond * 1)
			}
			//log.Println("r1", len(data))
		}()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb.Write(betters[0])
	}
	//time.Sleep(time.Second * 10)
	//log.Println("write", t)
	//done <- 0
	//time.Sleep(time.Second * 10)
}
func BenchmarkRingBuffer_AsyncChan(b *testing.B) {
	//rb := New(1024)
	done := make(chan byte, 1)
	var betters [][]byte
	for i := 'A'; i <= 'Z'; i++ {
		betters = append(betters, []byte{byte(i)})
	}
	ch := make(chan []byte, 4096)
	//b.StartTimer()
	//defer b.StopTimer()
	for i := 0; i < 4; i++ {
		go func() {
			buf := make([]byte, 2)
			//var data []byte
			for len(done) == 0 {
				b, ok := <-ch
				if ok {
					buf = append(buf, b...)
				}
				//time.Sleep(time.Millisecond * 1)
				log.Println(len(done))
			}
			log.Println("r1", len(buf), buf)
		}()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ch <- betters[0]
	}
	time.Sleep(time.Second * 5)
	log.Println("write", len(ch))
	done <- 1
	time.Sleep(time.Second * 10)
}
