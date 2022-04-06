package hzmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mgo struct {
	Db *mongo.Database
	TableMap map[string]string
}

func NewMgo(db *mongo.Database) *Mgo {
	tableMap := make(map[string]string)
	ret := &Mgo{
		db,
		tableMap,
	}
	return ret
}

func CreateMgo(url string, dbname string) (*Mgo, error) {
	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}

	db := client.Database(dbname)

	mgo := NewMgo(db)
	return mgo, nil
}

func (m *Mgo) SetTableMap(data map[string]string) {
	m.TableMap = data
}

func (m *Mgo) GetCollection(name string, opts ...*options.CollectionOptions) *Collection {
	tableName := name
	if dst, exists := m.TableMap[name]; exists == true {
		tableName = dst
	}
	ret := &Collection{
		m.Db.Collection(tableName, opts...),
	}
	return ret
}