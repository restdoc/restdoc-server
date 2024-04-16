package redispool

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/glog"

	"restdoc/config"

	Models "github.com/restdoc/restdoc-models"
)

var CacheDb *redis.Client
var SessionDb *redis.Client

const SESSION_PREFIX = "restdoc_session"

func Init() {

	if CacheDb == nil {
		CacheDb = redis.NewClient(&redis.Options{
			Addr:         config.DefaultConfig.CacheRedis.Address,
			Password:     config.DefaultConfig.CacheRedis.Password,
			DialTimeout:  1 * time.Second,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
			PoolSize:     config.DefaultConfig.CacheRedis.PoolSize,
			PoolTimeout:  3 * time.Second,
			DB:           config.DefaultConfig.CacheRedis.DB,
		})
		_, err := CacheDb.Ping().Result()
		if err != nil {
			glog.Error("cache redis db connect error", err)
		}
	}

	if SessionDb == nil {
		SessionDb = redis.NewClient(&redis.Options{
			Addr:         config.DefaultConfig.SessionRedis.Address,
			Password:     config.DefaultConfig.SessionRedis.Password,
			DialTimeout:  1 * time.Second,
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
			PoolSize:     config.DefaultConfig.SessionRedis.PoolSize,
			PoolTimeout:  3 * time.Second,
			DB:           config.DefaultConfig.SessionRedis.DB,
		})
		_, err := SessionDb.Ping().Result()
		if err != nil {
			glog.Error("session redis db connect error", err)
		}
	}
}

func GetSession(session_id string) (Models.Session, error) {

	var session Models.Session
	key := SESSION_PREFIX + session_id
	data, err := SessionDb.Get(key).Result()
	if err != nil {
		return session, err
	} else {
		err = json.Unmarshal([]byte(data), &session)
		if err != nil {
			glog.Error("unmarshal session error", err)
			return session, err
		} else {
			session.Login = "true"
			return session, nil
		}
	}
}

func SetSession(session_id string, session Models.Session, expire int) error {

	key := SESSION_PREFIX + session_id
	data, err := json.Marshal(session)
	if err != nil {
		glog.Error("session marshal error", err)
		return err
	}
	err = SessionDb.Set(key, data, time.Duration(expire)*time.Second).Err()
	if err != nil {
		glog.Error("set session error", err)
		return err
	} else {
		return nil
	}
}

func DeleteSession(session_id string) error {

	key := "session_" + session_id
	err := SessionDb.Del(key).Err()
	if err != nil {
		glog.Error("delete session error", err)
		return err
	} else {
		return nil
	}
}

func GetResendState(user string) (string, error) {

	key := "resend_" + user
	data, err := CacheDb.Get(key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		glog.Errorf("get resend state error %v", err)
		return "", err
	} else {
		return data, nil
	}
}

func SetResendState(user string, expire int) error {

	key := "resend_" + user
	err := CacheDb.Set(key, "true", time.Duration(expire)*time.Second).Err()
	if err != nil {
		glog.Errorf("set resend key error %v", err)
		return err
	} else {
		return nil
	}
}

func SetResetpasswordState(key string, email string, expire int) error {

	k := "resetpassword_" + key
	err := CacheDb.Set(k, email, time.Duration(expire)*time.Second).Err()
	if err != nil {
		glog.Error("set reset key error", err)
		return err
	} else {
		return nil
	}
}

func GetResetpasswordState(key string) (string, error) {

	k := "resetpassword_" + key
	data, err := CacheDb.Get(k).Result()
	if err != nil {
		glog.Error("set reset key error", err)
		return "", err
	} else {
		return data, nil
	}
}

func DeleteResetpasswordState(key string) error {

	k := "resetpassword_" + key
	err := SessionDb.Del(k).Err()
	if err != nil {
		glog.Error("delete session error", err)
		return err
	} else {
		return nil
	}
}

func GetGtSession(session_id string) (Models.Session, error) {

	var session Models.Session
	key := "gtsession_" + session_id
	data, err := SessionDb.Get(key).Result()
	if err != nil {
		return session, err
	} else {
		err = json.Unmarshal([]byte(data), &session)
		if err != nil {
			glog.Error("unmarshal session error", err)
			return session, err
		} else {
			session.Login = "true"
			return session, nil
		}
	}
}

func SetGtSession(session_id string, session map[string]string, expire int) error {

	key := "session_" + session_id
	data, err := json.Marshal(session)
	if err != nil {
		glog.Error("session marshal error", err)
		return err
	}
	err = SessionDb.Set(key, data, time.Duration(expire)*time.Second).Err()
	if err != nil {
		glog.Error("set session error", err)
		return err
	} else {
		return nil
	}
}

func DeleteGtSession(session_id string) error {

	key := "session_" + session_id
	err := SessionDb.Del(key).Err()
	if err != nil {
		glog.Error("delete session error", err)
		return err
	} else {
		return nil
	}
}

func IncrUserLabel(user_id string, label_id string, number int64) (int64, error) {

	key := "c_" + user_id
	count, err := SessionDb.HIncrBy(key, label_id, number).Result()
	if err != nil {
		return count, err
	} else {
		return count, nil
	}
}

func GetUserLabelCount(user_id string) (map[string]string, error) {

	key := "c_" + user_id
	result, err := SessionDb.HGetAll(key).Result()
	if err != nil {
		return result, err
	} else {
		return result, nil
	}
}
