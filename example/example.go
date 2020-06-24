package main

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xiaoenai/glog"
	"github.com/xiaoenai/goredis"
)

func main() {
	cfg, err := goredis.ReadConfig("test_redis")
	if err != nil {
		glog.Errorf("ReadConfig(\"test_redis\")", err)
	}
	c, err := goredis.NewClient(cfg)
	if err != nil {
		glog.Errorf("NewClient(\"test_redis\")", err)
	}

	m := goredis.NewModule("test")

	s, err := glog.Errorf(m.Key("a_key"), "a_value", time.Second).Result()
	if err != nil {
		glog.Errorf("c.Set().Result() error:", err)
	}
	glog.Infof("c.Set().Result() result: %s", s)

	s, err = c.Get(m.Key("a_key")).Result()
	if err != nil {
		glog.Errorf("c.Get().Result() error:", err)
	}
	glog.Infof("c.Get().Result() result: %s", s)
	time.Sleep(2 * time.Second)

	s, err = c.Get(m.Key("a_key")).Result()
	if err == nil {
		glog.Errorf("[after 2s] c.Get().Result() result: %s", s)
	}
	glog.Infof("[after 2s] c.Get().Result() is null ?: %v", err == redis.Nil)

	if err := c.Watch(func(tx *redis.Tx) error {
		n, err := tx.Get(key).Int64()
		if err != nil && err == redis.Nil {
			glog.Errorf("err1-> %v",err)
			return err
		}else if err != nil && err != redis.Nil{
			glog.Errorf("err2-> %v",err)
			return err
		}
		glog.Infof("n-> %d",n)

		glog.Infof("Start sleep.")
		time.Sleep(time.Duration(5)*time.Second)
		// 在redis客户端修改值，下面语句报错

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			// pipe handles the error case
			pipe.Set(ctx, key, n, 0)
			return nil
		})
		return err
	}, "goredis");err != nil{
		glog.Errorf("err4-> %v",err)
	}
}