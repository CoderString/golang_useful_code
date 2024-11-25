package consume

import (
	"context"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"log"
	"time"
)

func pull() {
	resp, err := pullConsumer.Pull(context.TODO(), 1)
	if err != nil {
		log.Printf("[pull error] err=%v", err)
		time.Sleep(sleepTime)
		return
	}
	switch resp.Status {
	case primitive.PullFound:
		log.Printf("[pull message successfully] MinOffset:%d, MaxOffset:%d, nextOffset: %d, len:%d\n", resp.MinOffset, resp.MaxOffset, resp.NextBeginOffset, len(resp.GetMessages()))
		var queue *primitive.MessageQueue
		if len(resp.GetMessages()) <= 0 {
			return
		}
		for _, msg := range resp.GetMessageExts() {
			// todo LOGIC CODE HERE
			queue = msg.Queue

			//log.Println(msg.Queue, msg.QueueOffset, msg.GetKeys(), msg.MsgId, string(msg.Body))
			log.Println(msg)
		}
		// update offset
		err = pullConsumer.UpdateOffset(queue, resp.NextBeginOffset)
		if err != nil {
			log.Printf("[pullConsumer.UpdateOffset] err=%v", err)
		}

	case primitive.PullNoNewMsg, primitive.PullNoMsgMatched:
		log.Printf("[no pull message]   next = %d\n", resp.NextBeginOffset)
		time.Sleep(sleepTime)
		return
	case primitive.PullBrokerTimeout:
		log.Printf("[pull broker timeout]  next = %d\n", resp.NextBeginOffset)

		time.Sleep(sleepTime)
		return
	case primitive.PullOffsetIllegal:
		log.Printf("[pull offset illegal] next = %d\n", resp.NextBeginOffset)
		return
	default:
		log.Printf("[pull error]  next = %d\n", resp.NextBeginOffset)
	}
}
