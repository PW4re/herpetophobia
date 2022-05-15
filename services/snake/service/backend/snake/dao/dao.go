package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"snake/db"
	"snake/objects"
)

func GetMap(id int) objects.Level {
	res := db.Get(db.DbName, db.ColName, bson.M{"id": id})
	var level objects.Level
	_ = res.Decode(&level)
	return level
}

func SaveMap(level objects.Level) {
	_, _ = db.InsertDoc(db.DbName, db.ColName, level)
}

func IncCounter(level objects.Level) {
	filter := bson.D{{"id", level.Id}}
	update := bson.D{{"$inc", bson.D{{"counter", 1}}}}
	_, _ = db.UpdateDoc(db.DbName, db.ColName, filter, update)
}

func ListId(limit int64, offset int64) objects.Ids {
	opts := options.Find().SetProjection(bson.D{{"id", 1}}).SetLimit(limit)
	results, _ := db.List("snake", "level", bson.D{}, opts)
	var listId []int
	for _, result := range results {
		mRes := result.Map()
		listId = append(listId, mRes["id"].(int))
	}
	return objects.Ids{Ids: listId}
}
