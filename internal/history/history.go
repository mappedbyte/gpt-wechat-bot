package history

type History interface {
	SetUserHistory(text any)
	GetHistory() []any
	ClearUserHistory()
}
