package redis

import (
	"github.com/garyburd/redigo/redis"
)

/**
*redis client wapper
*author liuyoubin@baidu.com
**/

type RedisError struct{
	Errno int
	Errmsg string
	Err error
}

type RedisClient struct{
	conn *redis.Conn
}

func (re *RedisError) Error() string{
	return "Errno:"+string(re.Errno)+",Errmsg:"+re.Errmsg+",Error:"+re.Err.Error()
}

func NewRedisClient() (redisClient *RedisClient){
	return &RedisClient{conn:nil}
}

func (redisClient *RedisClient) OpenConn(host string, port string) (err error){
	rs, err := redis.Dial("tcp", host+":"+port)
	
	
	if err != nil {
		return &RedisError{Errno: 1, Errmsg: "redis connect failed", Err: err}
	}
	
	//rs.Do("SELECT", 0)
	
	(*redisClient).conn = &rs
	return
}

func (redisClient RedisClient)CloseConn(){
	(*(redisClient).conn).Close()
}

func (redisClient RedisClient) Set(key string, value interface{}) (err error){
	_ ,err = (*redisClient.conn).Do("SET", key, value)
	if err != nil {
		return &RedisError{Errno: 1, Errmsg: "redis set failed", Err: err}
	}
	
	return
}

func (redisClient RedisClient) Get(key string) (result string, err error){
	value, err := (*redisClient.conn).Do("GET", key)
	if err != nil {
		return "", &RedisError{Errno:1, Errmsg:"redis get failed", Err:err}
	}
	
	if value == nil {
		return "", nil
	}
	
	result,err = redis.String(value, err)
	if err != nil {
		return "", &RedisError{Errno:1, Errmsg:"redis convert to string failed", Err:err}
	}
	return result, nil
}

func (redisClient RedisClient) Del(key string) (err error){
	_, err = (*redisClient.conn).Do("DEL", key)
	if err != nil {
		return &RedisError{Errno:1, Errmsg:"redis del failed", Err:err}
	}
	return nil
}

func (redisClient RedisClient) SetEx(key string, expire int, value interface{})(err error){
	_ ,err = (*redisClient.conn).Do("SETEX", key, expire, value)
	if err != nil {
		return &RedisError{Errno: 1, Errmsg: "redis setex failed", Err: err}
	}
	
	return
}

func (redisClient RedisClient) Expire(key string, expire int)(err error){
	_ ,err = (*redisClient.conn).Do("EXPIRE", key, expire)
	if err != nil {
		return &RedisError{Errno: 1, Errmsg: "redis expire failed", Err: err}
	}
	
	return
} 

func (redisClient RedisClient) Command(cmd string, args...interface{}) (value interface{}, err error){
	value, err = (*redisClient.conn).Do(cmd, args...)
	if err != nil {
		return nil, &RedisError{Errno:1, Errmsg:"command "+cmd+" failed", Err:err}
	}
	
	return value, err
}


