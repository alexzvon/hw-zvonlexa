package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	circleCh := in
	for _, stage := range stages {
		circleCh = func(inCircle In) Out {
			joinCh := make(chan interface{})

			go func() {
				defer close(joinCh)

				for {
					select {
					case <-done:
						return
					case v, ok := <-inCircle:
						if !ok {
							return
						}
						select {
						case <-done:
							return
						case joinCh <- v:
						}
					}
				}
			}()

			return stage(joinCh)
		}(circleCh)
	}

	return circleCh
}
