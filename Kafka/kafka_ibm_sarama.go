package ibm_sarama

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestSendSyncMessage(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{"192.168.0.102:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	defer client.Close()

	msg := &sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.StringEncoder("Hello Kafka from Sarama!"),
	}

	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Message sent successfully! PID: %v, Offset: %v\n", pid, offset)
}

func TestSendSyncBatchMessage(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{"192.168.0.102:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	defer client.Close()

	msgs := make([]*sarama.ProducerMessage, 0)
	for _, str := range []string{"msg1", "msg2", "msg3"} {
		msg := &sarama.ProducerMessage{
			Topic: "topic",
			Value: sarama.ByteEncoder(str),
		}
		msgs = append(msgs, msg)
	}
	if err := client.SendMessages(msgs); err != nil {
		fmt.Printf("send error: %#v\n", err)
	}
}

func TestSendAsyncMessage(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal        // 仅等待本地副本确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机分区
	config.Producer.Return.Successes = true                   // 返回成功消息
	config.Producer.Return.Errors = true                      // 返回错误消息

	producer, err := sarama.NewAsyncProducer([]string{"192.168.0.102:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to start Kafka producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Failed to close Kafka producer: %v", err)
		}
	}()

	go handleResponses(producer)
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Async message %d", i+1)
		producer.Input() <- &sarama.ProducerMessage{
			Topic: "test",
			Value: sarama.StringEncoder(message),
		}
		fmt.Printf("Message queued: %s\n", message)
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("All messages sent asynchronously!")
}

func handleResponses(producer sarama.AsyncProducer) {
	for {
		select {
		case success := <-producer.Successes():
			fmt.Printf("Message sent successfully: topic=%s, partition=%d, offset=%d\n",
				success.Topic, success.Partition, success.Offset)
		case err := <-producer.Errors():
			fmt.Printf("Failed to send message: %v\n", err)
		}
	}

}

func TestConsumerMessage(t *testing.T) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_8_0_0

	consumer, err := sarama.NewConsumer([]string{"192.168.0.102:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	fmt.Println("Waiting for messages...")

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Received message: key=%s value=%s topic=%s partition=%d offset=%d\n",
				string(msg.Key), string(msg.Value), msg.Topic, msg.Partition, msg.Offset)

		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error occurred: %v\n", err)
		}
	}
}

// 消费者结构
type ConsumerGroupHandler struct{}

// 处理消息
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer group is set up.")
	return nil
}

func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	fmt.Println("Consumer group is cleaned up.")
	return nil
}

func (ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message received: key=%s, value=%s, topic=%s, partition=%d, offset=%d\n",
			string(msg.Key), string(msg.Value), msg.Topic, msg.Partition, msg.Offset)
		// 手动确认消息
		session.MarkMessage(msg, "")
	}
	return nil
}

func TestConsumerMessageByGroup(t *testing.T) {
	config := sarama.NewConfig()
	groupID := "my-consumer-group"
	brokers := []string{"192.168.0.102:9092"}
	topics := []string{"test"}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM)
		<-c
		cancel()
	}()

	fmt.Println("Starting Kafka consumer group...")
	handler := ConsumerGroupHandler{}
	for {
		if err := consumerGroup.Consume(ctx, topics, handler); err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}

		if ctx.Err() != nil {
			fmt.Println("Consumer group terminated.")
			return
		}
	}
}

