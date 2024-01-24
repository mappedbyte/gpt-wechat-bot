package gpt

type Gpt interface {
	Chat(msg []string) (string, error)
}
