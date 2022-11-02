package dao

import (
	"context"
	"github.com/stellarisJAY/goim/pkg/db"
	"github.com/stellarisJAY/goim/pkg/db/model"
	"go.mongodb.org/mongo-driver/bson"
)

func AddNotification(notification *model.Notification) error {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	collection := database.Collection(db.CollectionNotification)
	_, err := collection.InsertOne(context.TODO(), notification)
	return err
}

func MarkNotificationRead(id int64) error {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	collection := database.Collection(db.CollectionNotification)
	filter := bson.D{{"id", id}}
	update := bson.D{{"$set", bson.D{{"read", true}}}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func RemoveNotification(id int64) error {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	collection := database.Collection(db.CollectionNotification)
	filter := bson.D{{"id", id}}
	_, err := collection.DeleteMany(context.TODO(), filter)
	return err
}

func ListNotReadNotifications(receiver int64) ([]*model.Notification, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	collection := database.Collection(db.CollectionNotification)
	filter := bson.D{{"receiver", receiver}, {"read", false}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	notifications := make([]*model.Notification, 0)
	err = cursor.All(context.TODO(), &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func ListAllNotifications(receiver int64) ([]*model.Notification, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	collection := database.Collection(db.CollectionNotification)
	filter := bson.D{{"receiver", receiver}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	notifications := make([]*model.Notification, 0)
	err = cursor.All(context.TODO(), &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func ListNotificationOfType(receiver int64, nType byte) ([]*model.Notification, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	collection := database.Collection(db.CollectionNotification)
	filter := bson.D{{"receiver", receiver}, {"type", nType}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	notifications := make([]*model.Notification, 0)
	err = cursor.All(context.TODO(), &notifications)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func GetNotification(id int64) (*model.Notification, error) {
	database := db.DB.MongoDB.Database(db.MongoDBName)
	collection := database.Collection(db.CollectionNotification)
	filter := bson.D{{"id", id}}
	result := collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	notification := &model.Notification{}
	if err := result.Decode(notification); err != nil {
		return nil, err
	}
	return notification, nil
}
