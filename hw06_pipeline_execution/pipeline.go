package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make([]Out, len(stages)+1)
	out[0] = in
	for i, stage := range stages {
		select {
		case <-done:
			return nil
		default:
		}
		out[i+1] = stage(out[i])
	}
	return out[len(stages)]
}

/*func worker(name string, jobQueue In, jobDo Stage, brk In, bufferSize int) Out {
	var jq, jr interface{}
	jobResult := make(Bi, bufferSize)
	go func() {
		defer close(jobResult)
		jget := false
		for {
			select {
			case <-brk:
				return
			case jq, jget = <-jobQueue:
			}
			if jget {
				jr = jobDo(jq)
				select {
				case <-brk:
					return
				case jobResult <- jr:
				}
			}
		}
	}()
	return jobResult
}

func getFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}*/
