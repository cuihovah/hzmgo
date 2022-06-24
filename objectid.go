package hzmgo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

func GetObjectIDs(list []string) ([]interface{}, error) {
	ret := make([]interface{}, len(list))
	var err error
	for i, v := range list {
		d, e := primitive.ObjectIDFromHex(v)
		if e != nil {
			err = e
		}
		ret[i] = d
	}
	return ret, err
}

func NewKeys() []primitive.ObjectID {
	ret := make([]primitive.ObjectID, 0)
	return ret
}

func GetObjectIDsByInterface(list []interface{}) ([]interface{}, error) {
	ret := make([]interface{}, len(list))
	var err error
	for i, v := range list {
		d, e := primitive.ObjectIDFromHex(v.(string))
		if e != nil {
			err = e
		}
		ret[i] = d
	}
	return ret, err
}

func GetBsonMap(list interface{}) bson.M {
	ret := make(bson.M)
	for i := 0; i < reflect.ValueOf(list).Len(); i++ {
		data := reflect.ValueOf(list).Index(i).FieldByName("Id")
		var key string
		if data.Kind() == reflect.String {
			key = data.String()
		} else {
			key = data.Interface().(primitive.ObjectID).Hex()
		}
		ret[key] = reflect.ValueOf(list).Index(i).Interface()
	}
	return ret
}
