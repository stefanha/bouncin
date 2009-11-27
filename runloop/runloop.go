package runloop

var funcs chan func()

func CallLater(f func()) {
	funcs <- f
}

func Quit() {
	CallLater(nil)
}

func Run() {
	funcs = make(chan func(), 64);
	for {
		f := <-funcs;
		if f == nil {
			return
		}
		f();
	}
}
