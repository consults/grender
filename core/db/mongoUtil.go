package db

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"grender/core/configReader"
	"time"
)

type MongoUtil struct {
	client        *mongo.Client
	db            *mongo.Database
	RenderHandler *mongo.Collection
}

func (m *MongoUtil) Disconnect() {
	if err := m.client.Disconnect(context.TODO()); err != nil {
		log.Println("断开MongoDB失败！")
	}
}

func (m *MongoUtil) Connect(cfg configReader.MongoCfg) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Url))
	if err != nil {
		panic(err)
	}
	m.client = client
	m.db = m.client.Database(cfg.DB)
	m.RenderHandler = m.db.Collection(cfg.RenderCol)
}

func (m *MongoUtil) InsertOne(url, html string, xpath string, render bool) {
	now := time.Now()
	ts := now.Format("2006-01-02 15:04:05")
	doc := bson.M{"url": url, "html": html, "ts": ts, "xpath": xpath, "render": render}
	_, err := m.RenderHandler.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Println(".....", err)
		log.Println("写入HTML失败！对应URL：", url)
	}
}
