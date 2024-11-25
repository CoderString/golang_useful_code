package topic

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func GetTopicList() {
	nameSrvAddr := []string{"127.0.0.1:9876"}
	adminClient, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(nameSrvAddr)),
		//admin.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
	)
	if err != nil {
		fmt.Println("adminClient new error:", err.Error())
	}
	defer adminClient.Close()
	result, err := adminClient.FetchAllTopicList(context.Background())
	if err != nil {
		fmt.Println("FetchAllTopicList error:", err.Error())
	}

	for _, topicName := range result.TopicList {
		fmt.Println("topic:", topicName)
	}
}

func CreateTopic(topicName string) {
	nameSrvAddr := []string{"127.0.0.1:9876"}
	brokerAddr := "127.0.0.1:10911"
	adminClient, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(nameSrvAddr)),
		//admin.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
	)
	if err = adminClient.CreateTopic(context.Background(), admin.WithTopicCreate(topicName), admin.WithBrokerAddrCreate(brokerAddr)); err != nil {
		fmt.Println("Create topic error:", err.Error())
	}
	if err = adminClient.Close(); err != nil {
		fmt.Printf("Shutdown admin error: %s", err.Error())
	}
}

func DeleteTopic(topicName string) {
	nameSrvAddr := []string{"127.0.0.1:9876"}
	brokerAddr := "127.0.0.1:10911"
	adminClient, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(nameSrvAddr)),
		//admin.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
	)
	if err = adminClient.DeleteTopic(context.Background(), admin.WithTopicDelete(topicName), admin.WithBrokerAddrDelete(brokerAddr), admin.WithNameSrvAddr(nameSrvAddr)); err != nil {
		fmt.Println("Delete topic error:", err.Error())
	}

	if err = adminClient.Close(); err != nil {
		fmt.Printf("Shutdown admin error: %s", err.Error())
	}
}
