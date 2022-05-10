package objects

type Level struct {
	id      int
	secret  string
	counter int
	init    [256]byte
	flag    string
}
