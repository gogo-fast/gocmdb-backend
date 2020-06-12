package utils

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

// 由于头像上传没有在注册用户的时候进行，库中的avatar是空的，导致反序列化的时候会报错（不能讲 null 反序列化到 string 类型），
// 所以这里需要使用 sql.NullString 代替 string 类型。这样null序列化后的是 {String: "", Valid: false}
// 为了只得到String类型的结果，我们需要重自定义类型NullString并重写 Scan, MarshalJSON, UnmarshalJSON 三个方法。

/*
implement one
*/
type NullString sql.NullString

func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		Logger.Error(err)
		return err
	}

	if reflect.TypeOf(value) == nil {
		*ns = NullString{s.String, false}
	} else {
		*ns = NullString{s.String, true}
	}
	return nil
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil
	return err
}

/*
implement two
*/
//type NullString struct {
//	sql.NullString
//}
//
//func(v *NullString) MarshalJSON() ([]byte, error) {
//	if v.Valid {
//		return json.Marshal(v.String)
//	} else {
//		return json.Marshal(nil)
//	}
//}
//
//func (v NullString) UnmarshalJSON(data []byte) error {
//	var s *string
//	if err := json.Unmarshal(data, &s); err != nil {
//		return err
//	}
//	if s != nil {
//		v.Valid = true
//		v.String = *s
//	} else {
//		v.Valid = false
//	}
//	return nil
//}