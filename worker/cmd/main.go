package main

import (
	"context"
	"fmt"
	"gihub.com/mefourr/tgdevob/worker/internal/handler/worker"
	"gihub.com/mefourr/tgdevob/worker/internal/kafka"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"log/slog"
	//"log/slog"
	//"strconv"
)

const (
	brokers = "localhost:9092" // TODO: retrieve from sys env
	topic   = "tg_requests"
	group   = "example"
)

type UserHash struct {
	Id          int               `redis:"id"`
	TgUserId    string            `redis:"tg_user_id"`
	UserName    string            `redis:"tg_user_name"`
	FirstName   string            `redis:"tg_user_first_name"`
	LastName    string            `redis:"tg_user_last_name"`
	UpdateId    int               `redis:"update_id"`
	LastRequest *tgbotapi.Message `redis:"last_request"`
}

//func main() {
//	ctx := context.Background()
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//	//_ = rdb.FlushDB(ctx).Err()
//
//	var u UserHash
//	var key string = "user:" + "1234567"
//
//	all := rdb.HGetAll(ctx, key)
//
//	if all.Err() != nil {
//		panic(all.Err())
//	} else if len(all.Val()) == 0 {
//		fmt.Println("key not found")
//		u = UserHash{
//			Id:        0,
//			TgUserId:  "123456",
//			UserName:  "test",
//			FirstName: "test",
//			LastName:  "test",
//			UpdateId:  10,
//		}
//
//		if err := rdb.HSet(ctx, key, u).Err(); err != nil {
//			panic(err)
//		}
//	} else {
//		if err := all.Scan(&u); err != nil {
//			panic(err)
//		}
//		fmt.Printf("Loaded from Redis: %+v\n", u)
//	}
//}

//func main() {
//	ctx := context.Background()
//	rdb := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//	//_ = rdb.FlushDB(ctx).Err()
//
//	var u UserHash
//	key := "user:" + "123456"
//
//	res, err := rdb.Get(ctx, key).Result()
//
//	if errors.Is(err, redis.Nil) {
//		fmt.Println("key not found")
//		u = UserHash{
//			Id:          0,
//			TgUserId:    "123456",
//			UserName:    "test",
//			FirstName:   "test",
//			LastName:    "test",
//			UpdateId:    10,
//			LastRequest: &tgbotapi.Message{MessageID: 10000},
//		}
//
//		bytes, err := json.Marshal(&u)
//		if err != nil {
//			panic(err)
//		}
//		if err := rdb.Set(ctx, key, bytes, 0).Err(); err != nil {
//			panic(err)
//		}
//	} else if err != nil {
//		panic(err)
//	} else {
//		var cached UserHash
//		if err := json.Unmarshal([]byte(res), &cached); err != nil {
//			panic(err)
//		}
//		fmt.Printf("Unmarshaled: %+v\n", cached.LastRequest.MessageID)
//	}
//}

func main() {
	ctx := utils.Init()
	slog.InfoContext(ctx, "Logger for consumer is initialized")
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_ = rdb.FlushDB(ctx).Err()

	ping := rdb.Ping(context.Background())
	result, err := ping.Result()
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}
	slog.InfoContext(ctx, "after setting redis up result is ", result)

	consumer := kafka.NewConsumer(
		[]string{brokers},
		topic,
		group,
		worker.New(rdb),
	)
	if err := consumer.Consume(ctx); err != nil {
		panic(err)
	}
}

func ExampleClient(ctx context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
