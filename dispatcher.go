package main

var dispatchFuncs chan func()

func CallLater(f func()) {
	dispatchFuncs <- f
}

func RunDispatcher() {
	dispatchFuncs = make(chan func(), 64);
	for {
		f := <-dispatchFuncs;
		if f == nil {
			return
		}
		f();
	}
}
