package group

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"time"
)

func GetAllSubscriptionGroup() {
	nameSrvAddr := []string{"127.0.0.1:9876"}
	brokerAddr := "127.0.0.1:10911"
	adminClient, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(nameSrvAddr)),
		//admin.WithCredentials(primitive.Credentials{
		//	AccessKey: "RocketMQ",
		//	SecretKey: "12345678",
		//}),
	)
	result, err := adminClient.GetAllSubscriptionGroup(context.Background(), brokerAddr, 3*time.Second)
	if err != nil {
		fmt.Println("GetAllSubscriptionGroup error:", err.Error())
	}

	for k, v := range result.SubscriptionGroupTable {
		fmt.Println(k, "=", v)
	}
	if err = adminClient.Close(); err != nil {
		fmt.Printf("Shutdown admin error: %s", err.Error())
	}
}
