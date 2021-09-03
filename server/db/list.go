package db

import "go-cache/internal"

func lPush(db *DB, key string, value interface{}) *internal.Payload {
	resCode := db.Data.Exec(key, func(s *shared) interface{} {
		v, ok := s.data[key]
		if !ok {
			temp := make([]interface{}, 1, 16)
			temp[0] = value
			s.data[key] = temp
			return 1
		}
		temp := append([]interface{}{value}, v.([]interface{})...)
		s.data[key] = temp
		return 1
	})

	p := &internal.Payload{
		Value: resCode,
	}
	return p
}

func lPop(db *DB, key string, value interface{}) *internal.Payload {

	resCode := db.Data.Exec(key, func(s *shared) interface{} {
		v, ok := s.data[key]
		if !ok {
			return key + " is not exist"
		}
		length := len(v.([]interface{}))
		if length == 0 {
			return "nil"
		}
		res := v.([]interface{})[0]
		s.data[key] = v.([]interface{})[1:]
		return res
	})
	p := &internal.Payload{
		Value: resCode,
	}
	return p
}

func rPush(db *DB, key string, value interface{}) *internal.Payload {

	resCode := db.Data.Exec(key, func(s *shared) interface{} {
		v, ok := s.data[key]
		if !ok {
			temp := make([]interface{}, 1, 16)
			temp[0] = value
			s.data[key] = temp
			return 1
		}
		temp := append(v.([]interface{}), value)
		s.data[key] = temp
		return 1
	})

	p := &internal.Payload{
		Value: resCode,
	}
	return p
}

func rPop(db *DB, key string, value interface{}) *internal.Payload {

	resCode := db.Data.Exec(key, func(s *shared) interface{} {
		v, ok := s.data[key]
		if !ok {
			return key + " is not exist"
		}
		length := len(v.([]interface{}))
		if length == 0 {
			return "nil"
		}
		res := v.([]interface{})[length-1]
		s.data[key] = v.([]interface{})[:length-1]
		return res
	})
	p := &internal.Payload{
		Value: resCode,
	}
	return p
}

func init() {

	RegisterCommand("lpush", lPush)
	RegisterCommand("lpop", lPop)
	RegisterCommand("rpush", rPush)
	RegisterCommand("rpop", rPop)
}
