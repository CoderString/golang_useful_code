package consume

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"log"
	"time"
)

const (
	nameSrvAddr       = "http://127.0.0.1:9876"
	accessKey         = "rocketmq"
	secretKey         = "12345678"
	topic             = "test-topic"
	consumerGroupName = "testPullGroup"
	tag               = "testPull"
	namespace         = "ns"
)

var pullConsumer rocketmq.PullConsumer
var sleepTime = 1 * time.Second

func poll() {
	cr, err := pullConsumer.Poll(context.TODO(), time.Second*5)
	if consumer.IsNoNewMsgError(err) {
		return
	}
	if err != nil {
		log.Printf("[poll error] err=%v", err)
		time.Sleep(sleepTime)
		return
	}
	// todo LOGIC CODE HERE
	log.Println("msgList: ", cr.GetMsgList())
	log.Println("messageQueue: ", cr.GetMQ())
	log.Println("processQueue: ", cr.GetPQ())
	// pullConsumer.ACK(context.TODO(), cr, consumer.ConsumeRetryLater)
	pullConsumer.ACK(context.TODO(), cr, consumer.ConsumeSuccess)
}
