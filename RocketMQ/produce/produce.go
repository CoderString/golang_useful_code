package produce

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"os"
	"strconv"
	"sync"
	"time"
)

func SendSyncMessage(topic string) {
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(2),
	)

	if err != nil {
		fmt.Println("adminClient new error: " + err.Error())
		os.Exit(1)
	}

	if err = p.Start(); err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	defer p.Shutdown()

	for i := 0; i < 10; i++ {
		res, err := p.SendSync(context.Background(), primitive.NewMessage(topic, []byte(fmt.Sprintf("Hello RocketMQ %d", i))))
		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
		time.Sleep(1 * time.Second)
	}

}

func SendAsyncMessage(topic string) {
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(2))

	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		err := p.SendAsync(context.Background(),
			func(ctx context.Context, result *primitive.SendResult, e error) {
				if e != nil {
					fmt.Printf("receive message error: %s\n", err)
				} else {
					fmt.Printf("send message success: result=%s\n", result.String())
				}
				wg.Done()
			}, primitive.NewMessage(topic, []byte("Hello RocketMQ Go Client!")))

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		}
	}
	wg.Wait()
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

func SendBatchMessage(topic string) {
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	var msgs []*primitive.Message
	for i := 0; i < 10; i++ {
		msgs = append(msgs, primitive.NewMessage(topic,
			[]byte("Hello RocketMQ Go Client! num: "+strconv.Itoa(i))))
	}

	res, err := p.SendSync(context.Background(), msgs...)

	if err != nil {
		fmt.Printf("send message error: %s\n", err)
	} else {
		fmt.Printf("send message success: result=%s\n", res.String())
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

func SendDelayMessage(topic string) {
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	for i := 0; i < 10; i++ {
		msg := primitive.NewMessage(topic, []byte("Hello RocketMQ Go Client!"))
		msg.WithDelayTimeLevel(3)
		res, err := p.SendSync(context.Background(), msg)

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

func SendTransactionMessage(topic string) {
	p, _ := rocketmq.NewTransactionProducer(
		NewDemoListener(),
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(1),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s\n", err.Error())
		os.Exit(1)
	}

	for i := 0; i < 10; i++ {
		res, err := p.SendMessageInTransaction(context.Background(),
			primitive.NewMessage(topic, []byte("Hello RocketMQ again "+strconv.Itoa(i))))

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	time.Sleep(5 * time.Minute)
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

func SendMessageWithTag(topic string) {
	p, _ := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	tags := []string{"TagA", "TagB", "TagC"}
	for i := 0; i < 3; i++ {
		tag := tags[i%3]
		msg := primitive.NewMessage(topic,
			[]byte("Hello RocketMQ Go Client!"))
		msg.WithTag(tag)

		res, err := p.SendSync(context.Background(), msg)
		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}

func SendMessageWithNameSpace(topic string) {
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"127.0.0.1:9876"})),
		producer.WithRetry(2),
		//producer.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
		producer.WithNamespace("namespace"),
	)

	if err != nil {
		fmt.Println("init producer error: " + err.Error())
		os.Exit(0)
	}

	err = p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	for i := 0; i < 100000; i++ {
		res, err := p.SendSync(context.Background(), primitive.NewMessage(topic,
			[]byte("Hello RocketMQ Go Client!")))

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}
