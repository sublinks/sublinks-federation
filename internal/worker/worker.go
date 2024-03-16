package worker

type Worker interface {
	Process(msg []byte) error
}
