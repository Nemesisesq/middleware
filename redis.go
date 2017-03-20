package middleware

import (
	"context"
	"net/http"

	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"gopkg.in/redis.v5"
)

type RedisAccessor struct {
	client redis.Client
	addr   string
	pass   string
	db     int
}

func NewRedisAccessor(addr, pass string, db int) (*RedisAccessor, error) {
	logrus.WithFields(logrus.Fields{
		"address":  addr,
		"password": pass,
		"db":       db,
	}).Info("This is the information to connect to redis")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	logrus.Info(addr)

	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("redis %v", pong)

	return &RedisAccessor{*client, addr, pass, db}, nil
}

func (ca *RedisAccessor) Set(request *http.Request, client redis.Client) context.Context {
	//gcontext.Set(request, "client", *client)
	return context.WithValue(request.Context(), "redis_client", &client)
}

type redisClient struct {
	rca RedisAccessor
}

func NewRedisClient(RedisAccessor RedisAccessor) *redisClient {
	return &redisClient{RedisAccessor}
}

func (ca *redisClient) Middleware() negroni.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		ctx := ca.rca.Set(request, ca.rca.client)
		next(writer, request.WithContext(ctx))
}
}
