package db

import (
	"go-cache/internal"
	"strconv"
)

func set(d *DB, key string, value interface{}) *internal.Payload {

	resCode := d.Data.Set(key, value)

	res := &internal.Payload{
		Value: strconv.Itoa(resCode),
	}
	return res
}

func get(d *DB, key string, v interface{}) *internal.Payload {
	value, ok := d.Data.Get(key)

	res := &internal.Payload{}
	if ok {
		res.Value = value
	} else {
		res.Value = "the key is not existed"
	}

	return res
}

func del(d *DB, key string, v interface{}) *internal.Payload {

	resCode := d.Data.Del(key)
	res := &internal.Payload{
		Value: strconv.Itoa(resCode),
	}
	return res
}

func init() {
	RegisterCommand("set", set)
	RegisterCommand("get", get)
	RegisterCommand("delete", del)

}
