package objects

type Level struct {
	Id      string    `bson:"id"`
	Secret  string    `bson:"secret"`
	Counter int       `bson:"counter"`
	Init    [256]byte `bson:"init"`
	Flag    string    `bson:"flag"`
}
