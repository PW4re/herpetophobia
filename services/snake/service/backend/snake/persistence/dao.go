package persistence

const (
	DB_NAME         = "snake"
	COLLECTION_NAME = "level"
)

type LevelDao struct{}

func (*LevelDao) Get(id int) {
	//level, err := db.Get(DB_NAME, COLLECTION_NAME, bson.M{"id": id})

}

func (*LevelDao) List() {

}

func (*LevelDao) Ids() {

}

func (*LevelDao) Save() {

}
