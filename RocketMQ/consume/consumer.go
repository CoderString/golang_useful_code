package consume

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"log"
	"net/http"
	"os"
	"time"
)

func ConsumerSimpleMessage(topic string) {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		//consumer.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
	)
	if err != nil {
		fmt.Println("init consumer error: " + err.Error())
		os.Exit(0)
	}

	err = c.Subscribe(topic, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("subscribe callback: %v \n", msgs)
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("Shutdown Consumer error: %s", err.Error())
	}
}

func ConsumerBroadCastMessage() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerModel(consumer.BroadCasting),
	)
	err := c.Subscribe("min", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("subscribe callback: %v \n", msgs)
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("Shutdown Consumer error: %s", err.Error())
	}
}

func ConsumerOrderlyMessage() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		consumer.WithConsumerModel(consumer.Clustering),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset),
		consumer.WithConsumerOrder(true),
	)
	err := c.Subscribe("TopicTest", consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		orderlyCtx, _ := primitive.GetOrderlyCtx(ctx)
		fmt.Printf("orderly context: %v\n", orderlyCtx)
		fmt.Printf("subscribe orderly callback: %v \n", msgs)
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("Shutdown Consumer error: %s", err.Error())
	}
}

func ConsumerMessageWithTag() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithGroupName("testGroup"),
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
	)
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: "TagA || TagC",
	}
	err := c.Subscribe("TopicTest", selector, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		fmt.Printf("subscribe callback: %v \n", msgs)
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	time.Sleep(time.Hour)
	err = c.Shutdown()
	if err != nil {
		fmt.Printf("shutdown Consumer error: %s", err.Error())
	}
}

func ConsumerMessageByPoll() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	rlog.SetLogLevel("info")
	var nameSrv, err = primitive.NewNamesrvAddr(nameSrvAddr)
	if err != nil {
		log.Fatalf("NewNamesrvAddr err: %v", err)
	}
	pullConsumer, err = rocketmq.NewPullConsumer(
		consumer.WithGroupName(consumerGroupName),
		consumer.WithNameServer(nameSrv),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: accessKey,
			SecretKey: secretKey,
		}),
		consumer.WithNamespace(namespace),
		consumer.WithMaxReconsumeTimes(2),
	)
	if err != nil {
		log.Fatalf("fail to new pullConsumer: %v", err)
	}
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: tag,
	}
	err = pullConsumer.Subscribe(topic, selector)
	if err != nil {
		log.Fatalf("fail to Subscribe: %v", err)
	}
	err = pullConsumer.Start()
	if err != nil {
		log.Fatalf("fail to Start: %v", err)
	}

	for {
		poll()
	}
}

func ConsumerMessageByPull() {
	rlog.SetLogLevel("info")
	var nameSrv, err = primitive.NewNamesrvAddr(nameSrvAddr)
	if err != nil {
		log.Fatalf("NewNamesrvAddr err: %v", err)
	}
	pullConsumer, err = rocketmq.NewPullConsumer(
		consumer.WithGroupName(consumerGroupName),
		consumer.WithNameServer(nameSrv),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: accessKey,
			SecretKey: secretKey,
		}),
		consumer.WithNamespace(namespace),
		consumer.WithMaxReconsumeTimes(2),
	)
	if err != nil {
		log.Fatalf("fail to new pullConsumer: %v", err)
	}
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: tag,
	}
	err = pullConsumer.Subscribe(topic, selector)
	if err != nil {
		log.Fatalf("fail to Subscribe: %v", err)
	}
	err = pullConsumer.Start()
	if err != nil {
		log.Fatalf("fail to Start: %v", err)

	}

	refreshPersistOffsetDuration := time.Second * 5

	timer := time.NewTimer(refreshPersistOffsetDuration)
	go func() {
		for ; true; <-timer.C {
			err = pullConsumer.PersistOffset(context.TODO(), topic)
			if err != nil {
				log.Printf("[pullConsumer.PersistOffset] err=%v", err)
			}
			timer.Reset(refreshPersistOffsetDuration)
		}
	}()

	for i := 0; i <= 4; i++ {
		go func() {
			for {
				pull()
			}
		}()
	}
	// make current thread hold to see pull result. TODO should update here
	time.Sleep(10000 * time.Second)
}
