package hzmgo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

type Collection struct {
	coll *mongo.Collection
}

func (coll *Collection) FindId(ctx context.Context, id interface{}, ret interface{}, opts ...*options.FindOneOptions) error {
	err := coll.coll.FindOne(ctx, bson.M{
		"_id": id,
	}, opts...).Decode(ret)
	return err
}

func (coll *Collection) DeleteId(ctx context.Context, id interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ret, err := coll.coll.DeleteOne(ctx, bson.M{
		"_id": id,
	}, opts...)
	return ret, err
}

func (coll *Collection) Find(ctx context.Context, filter interface{}, ret interface{}, opts ...*options.FindOptions) error {
	cur, err := coll.coll.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	err = cur.All(ctx, ret)
	return err
}

func (coll *Collection) FindKeys(ctx context.Context, filter interface{}, ret interface{}, opts ...*options.FindOptions) error {
	cur, err := coll.coll.Find(ctx, bson.M{
		"_id": bson.M{
			"$in": filter,
		},
	}, opts...)
	if err != nil {
		return err
	}
	err = cur.All(ctx, ret)
	return err
}

func (coll *Collection) FindOne(ctx context.Context, filter interface{}, ret interface{}, opts ...*options.FindOneOptions) error {
	err := coll.coll.FindOne(ctx, filter, opts...).Decode(ret)
	return err
}

func (coll *Collection) UpdateId(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ret, err := coll.coll.UpdateOne(ctx, bson.M{
		"_id": id,
	}, update, opts...)
	return ret, err
}

func (coll *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ret, err := coll.coll.UpdateOne(ctx, filter, update, opts...)
	return ret, err
}

func (coll *Collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	ret, err := coll.coll.UpdateMany(ctx, filter, update, opts...)
	return ret, err
}

func (coll *Collection) Upsert(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	opts := options.UpdateOptions{}
	ret, err := coll.coll.UpdateOne(ctx, filter, update, opts.SetUpsert(true))
	return ret, err
}

func (coll *Collection) InsertOne(ctx context.Context, document interface{}, opt... *options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return coll.coll.InsertOne(ctx, document, opt...)
}

func (coll *Collection) InsertMany(ctx context.Context, document interface{}, opt... *options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	list := make([]interface{}, 0)
	for i := 0; i < reflect.ValueOf(document).Len(); i++ {
		data := reflect.ValueOf(document).Index(i).Interface()
		list = append(list, data)
	}
	return coll.coll.InsertMany(ctx, list, opt...)
}

func (coll *Collection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ret, err := coll.coll.DeleteMany(ctx, filter, opts...)
	return ret, err
}

func (coll *Collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	ret, err := coll.coll.DeleteOne(ctx, filter, opts...)
	return ret, err
}

func (coll *Collection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	ret, err := coll.coll.CountDocuments(ctx, filter, opts...)
	return ret, err
}

func (coll *Collection) ListObj(ctx context.Context, filter interface{}, ret interface{}, limit,
offset int64, sort ...interface{}) (ListObj, error) {
	total, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return ListObj{}, err
	}
	findOptions := options.Find()

	if limit == 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	if len(sort) > 0 {
		findOptions.SetSort(sort[0])
	}

	if err != nil {
		return ListObj{}, err
	}
	err = coll.Find(ctx, filter, ret, findOptions)
	if err != nil {
		return ListObj{}, err
	}
	list := reflect.ValueOf(ret).Elem().Interface()
	return ListObj{
		total,
		list,
	}, err
}

func (coll *Collection) Search(ctx context.Context, f interface{}, ret interface{}, limit,
	offset int64, sort ...interface{}) (ListObj, error) {
	filter := f.(bson.M)
	search := make(bson.M)
	for key, v := range filter {
		if reflect.TypeOf(v).Kind() == reflect.String && key != "dealer_id" {
			search[key] = bson.M{
				//"$regex": v.(string),    //区分大小写
				"$regex" : primitive.Regex{Pattern: v.(string), Options: "i"},   //不区分大小写
			}
		} else {
			search[key] = v
		}
	}

	total, err := coll.CountDocuments(context.TODO(), search)
	if err != nil {
		return ListObj{}, err
	}
	findOptions := options.Find()

	if limit == 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	findOptions.SetLimit(limit)
	findOptions.SetSkip(offset)
	if len(sort) > 0 {
		findOptions.SetSort(sort[0])
	}

	if err != nil {
		return ListObj{}, err
	}
	err = coll.Find(ctx, search, ret, findOptions)
	if err != nil {
		return ListObj{}, err
	}
	list := reflect.ValueOf(ret).Elem().Interface()
	return ListObj{
		total,
		list,
	}, err
}