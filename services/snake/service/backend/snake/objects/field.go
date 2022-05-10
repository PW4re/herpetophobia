package objects

type Level struct {
	Id      int
	Secret  string
	Counter int
	Init    [256]byte
	Flag    string
}
