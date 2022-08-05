package service

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/config"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/kafka"
	"github.com/stellarisJAY/goim/pkg/proto/pb"
	"github.com/stellarisJAY/goim/pkg/snowflake"
	"github.com/stellarisJAY/goim/pkg/wordfilter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageServiceImpl struct {
	transferProducer *kafka.Producer
	idGenerator      *snowflake.Snowflake
	wordFilter       wordfilter.Filter
}

func NewMessageServiceImpl() *MessageServiceImpl {
	transProducer, err := kafka.NewProducer(config.Config.Kafka.Addrs, pb.MessageTransferTopic)
	if err != nil {
		panic(err)
	}
	// 从配置文件读取敏感词
	filter := wordfilter.NewTrieTreeFilter()
	filter.Build(config.Config.SensitiveWords)
	return &MessageServiceImpl{
		transferProducer: transProducer,
		idGenerator:      snowflake.NewSnowflake(config.Config.MachineID),
		wordFilter:       filter,
	}
}

// SyncOfflineMessages 同步离线消息
// 1. 从MongoDB查询用户提供的序列号开始的消息
// 2. 将消息按照序列号排序
// 3. 返回消息列表
func (m *MessageServiceImpl) SyncOfflineMessages(ctx context.Context, request *pb.SyncMsgRequest) (*pb.SyncMsgResponse, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	// 按照seq排序
	opts := options.Find().SetSort(bson.D{{"seq", 1}})
	// 查询 to == userID AND seq > lastSeq
	cursor, err := database.
		Collection(db.CollectionOfflineMessage).
		Find(context.TODO(), bson.D{{"to", request.UserID}, {"seq", bson.D{{"$gt", request.LastSeq}}}}, opts)

	if err != nil {
		return &pb.SyncMsgResponse{Code: pb.Error, Message: err.Error()}, nil
	}
	// Unmarshal
	msgs := make([]*pb.BaseMsg, 0, cursor.RemainingBatchLength())
	for cursor.Next(context.TODO()) {
		raw := cursor.Current
		message := new(pb.BaseMsg)
		err := bson.Unmarshal(raw, message)
		if err != nil {
			continue
		}
		msgs = append(msgs, message)
	}
	response := &pb.SyncMsgResponse{
		Code:     pb.Success,
		Message:  "",
		Messages: msgs,
	}
	// 更新客户端的 initSeq 和 lastSeq
	if msgs != nil && len(msgs) > 0 {
		response.InitSeq = msgs[0].Seq
		response.LastSeq = msgs[len(msgs)-1].Seq
	}
	return response, nil
}
