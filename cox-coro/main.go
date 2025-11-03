package main

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

// We can understand the first parameter as the initial value
func mult_by_two(val int, yield func(int) int) int {
	for val >= 0 {
		val = yield(2 * val)
	}
	println("leaving coro")
	return val
}

func main() {
	resume := New(mult_by_two)
	for i := 0; i < 3; i++ {
		println(resume(i))
	}
	resume(-1)
}
