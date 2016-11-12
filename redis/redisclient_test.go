package redis

import (
    "testing"
    "fmt"
    "github.com/garyburd/redigo/redis"
)

func TestSet(t *testing.T) {
	redisClient := NewRedisClient()
	
	err := redisClient.OpenConn("127.0.0.1", "6379")
	if err != nil {
		t.Error("open conn failed:" + err.Error())
	}
	
	err = redisClient.Set("a", "b")
	if err != nil {
		t.Error("set failed" + err.Error())
	}
	
	value ,err := redisClient.Get("a")
	if err != nil {
		t.Error("get failed" + err.Error())
	}
	
	value = value.(string)
	fmt.Println(value)
	if value != "b" {
		t.Error("value is not ok");
	}
	
	err = redisClient.SetEx("c", 10, "d")
	if err != nil {
		t.Error("setex failed" + err.Error())
	}
	
	redisClient.CloseConn()
}

func TestCommand(t *testing.T) {
	redisClient := NewRedisClient()
	
	err := redisClient.OpenConn("127.0.0.1", "6379")
	if err != nil {
		t.Error("open conn failed:" + err.Error())
	}
	
	value, err := redisClient.Command("GET", "a")
	
	if err != nil {
		t.Error("command failed:" + err.Error())
	}
	
	value , err = redis.String(value, err)
	fmt.Println(value)
}


