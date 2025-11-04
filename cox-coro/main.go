package main

import "fmt"

func New[In, Out any](calculation func(In, func(Out) In) Out) (resumer func(In) Out) {
	cresume := make(chan In)
	cyield := make(chan Out)

	resume := func(in In) (out Out) {
		cresume <- in
		return <-cyield
	}

	yield := func(out Out) (in In) {
		cyield <- out
		return <-cresume
	}

	go func() {
		cyield <- calculation(<-cresume, yield)
	}()

	return resume
}

func counter() func(bool) int {
	return New(func(more bool, yield func(int) bool) int {
		for i := 2; more; i++ {
			fmt.Printf("passing %d from 1\n", i)
			more = yield(i)
		}
		return 0
	})
}

func filter(p int, next func(bool) int) (filtered func(bool) int) {
	return New(func(more bool, yield func(int) bool) int {
		for more {
			n := next(true)
			if n%p != 0 {
				fmt.Printf("passing %d from %d\n", n, p)
				more = yield(n)
			}
		}
		return next(false)
	})
}

func main() {
	next := counter()
	for i := 0; i < 10; i++ {
		p := next(true)
		fmt.Println(p)
		next = filter(p, next)
	}
	next(false)
}
