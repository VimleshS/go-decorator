package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

func main() {
	c := cacheWrapper(fib, &sync.Map{})
	w := wrapper(c, log.New(os.Stdout, "Debug! ", 1))

	fmt.Println(w(100000))
	fmt.Println(w(200000))
	fmt.Println(w(100000))
}

type fnc func(n int) float64

func cacheWrapper(fn fnc, m *sync.Map) fnc {
	return func(n int) float64 {
		if v, ok := m.Load(n); ok {
			return v.(float64)
		}
		res := fn(n)
		m.Store(n, res)
		return res
	}
}

func wrapper(fn fnc, logger *log.Logger) fnc {
	return func(n int) (res float64) {
		defer func(t time.Time) {
			logger.Printf("----> Res is %v, After %v \n", res, time.Since(t))
		}(time.Now())

		return fn(n)
	}
}

func fib(n int) float64 {
	ch := make(chan float64)

	for k := 0; k <= n; k++ {
		go func(ch chan float64, t float64) {
			ch <- 4 * math.Pow(-1, t) / (2*t + 1)
		}(ch, float64(k))
	}

	res := 0.0
	for k := 0; k <= n; k++ {
		res += <-ch
	}
	return res

}
