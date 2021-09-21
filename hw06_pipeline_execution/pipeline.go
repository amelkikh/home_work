package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := doneStage(in, done)
	for _, stage := range stages {
		if stage != nil {
			out = doneStage(stage(out), done)
		}
	}
	return out
}

func doneStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			// draining
			for range in {
			}
		}()

		for {
			select {
			case <-done:
				return
			default:
			}

			select {
			case <-done:
				return
			case d, ok := <-in:
				if !ok {
					return
				}
				out <- d
			}
		}
	}()
	return out
}
