package shutdown

type worker struct {
	priority int
	shutdown chan bool
	finished chan bool
}

func (w *worker) Stop() chan bool {
	return w.shutdown
}

func (w *worker) Finished() {
	w.finished <- true
}
