package gocontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {

	background := context.Background()
	fmt.Println(background)

	todo := context.TODO()
	fmt.Println(todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()

	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")
	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")
	contextF := context.WithValue(contextE, "f", "F")
	contextG := context.WithValue(contextD, "g", "G")
	contextH := context.WithValue(contextG, "h", "H")

	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)
	fmt.Println(contextH)
}

func CreateCounterWithLeak() chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			destination <- counter
			counter++
		}
	}()
	return destination
}

// go routine leak (never stop go routine)
func TestContextWithoutCancel(t *testing.T) {
	fmt.Println("Total Goroutines ", runtime.NumGoroutine())
	destination := CreateCounterWithLeak()
	for n := range destination {
		fmt.Println("Counter ", n)
		if n == 10 {
			break
		}
	}
	fmt.Println("Total Goroutines ", runtime.NumGoroutine())
}

func CreateCounterWithoutLeak(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
			}
		}
	}()
	return destination
}

func TestContextWithCancel(t *testing.T) {
	fmt.Println("Total Goroutines ", runtime.NumGoroutine())
	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)
	destination := CreateCounterWithoutLeak(ctx)
	for n := range destination {
		fmt.Println("Counter ", n)
		if n == 10 {
			break
		}
	}
	cancel()
	time.Sleep(2 * time.Second)
	fmt.Println("Total Goroutines ", runtime.NumGoroutine())
}

/*
Selain itu juga ada yang namanya context Timeout dan Deadline ini
hanya untuk melengkapi saja agar bisa berjaga2 jika ada leak
goroutine yang berjalan tanpa pengetahuan kita

*/
