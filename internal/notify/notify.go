package notify

type Notify interface {
	SendNotify(msg string) error
}
