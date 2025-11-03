package main

func New[In, Out any](calculation func(In) Out) (resumer func(In) Out) {
	cin := make(chan In)
	cout := make(chan Out)

	resume := func(in In) (out Out) {
		cin <- in
		return <-cout
	}

	go func() {
		cout <- calculation(<-cin)
	}()

	return resume
}

func double(val int) (double_of_val int) {
	return 2 * val
}

func main() {
	doubling_coro := New(double)
	println(doubling_coro(2))
}
