package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	tmp := pre(in, done)
	for _, stage := range stages {
		tmp = stage(pre(tmp, done))
	}
	return tmp
}

func pre(in, done In) Out {
	out := make(Bi, 100)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case x, ok := <-in:
				if ok {
					select {
					case out <- x:
					case <-done:
						return
					}
				} else {
					return
				}
			}
		}
	}()
	return out
}
