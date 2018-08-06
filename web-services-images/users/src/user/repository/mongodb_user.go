package repository

import (
	"context"
	"web-service-users/src/model"
	"web-service-users/src/user"

	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	log "github.com/sirupsen/logrus"
)

type mongodbUserRepository struct {
	DataBase   *mongo.Database
	Collection *mongo.Collection
}

func NewMongodbUserRepository(db *mongo.Database) user.Repository {
	return &mongodbUserRepository{db, db.Collection("users")}
}

func (m *mongodbUserRepository) userToDoc(user *model.User) *bson.Document {
	return bson.NewDocument(
		bson.EC.String("user_name", user.UserName),
		bson.EC.String("name", user.Name),
		bson.EC.Int32("quantity_kudos", user.QuantityKudos),
		bson.EC.DateTime("created_at", user.CreateAt.Unix()),
		bson.EC.DateTime("updated_at", user.UpdateAt.Unix()))
}

func (m *mongodbUserRepository) docToUser(doc *bson.Document) *model.User {
	return &model.User{
		UserName:      doc.Lookup("user_name").StringValue(),
		Name:          doc.Lookup("name").StringValue(),
		QuantityKudos: doc.Lookup("quantity_kudos").Int32(),
		CreateAt:      doc.Lookup("created_at").DateTime(),
		UpdateAt:      doc.Lookup("updated_at").DateTime(),
	}
}

func (m *mongodbUserRepository) Store(user *model.User) error {
	doc := m.userToDoc(user)
	_, err := m.Collection.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongodbUserRepository) GetByUserName(userName string) (*model.User, error) {

	result := m.Collection.FindOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("user_name", userName),
		))
	doc := bson.NewDocument()
	err := result.Decode(doc)
	if err != nil {
		return nil, err
	}
	userFound := m.docToUser(doc)
	return userFound, nil
}

func (m *mongodbUserRepository) DeleteByUserName(userName string) error {

	doc := bson.NewDocument()
	err := m.Collection.FindOne(
		context.Background(),
		bson.NewDocument(
			bson.EC.String("user_name", userName),
		)).Decode(doc)

	if err != nil {
		return err
	}
	_, err = m.Collection.DeleteOne(context.Background(), doc)
	if err != nil {
		return err
	}
	return nil

}

func (m *mongodbUserRepository) FetchAllUsers(pageSize int64, numberPage int64) ([]*model.User, error) {
	skip := pageSize * (numberPage - 1)
	cursor, err := m.Collection.Find(context.Background(), nil, findopt.Skip(skip), findopt.Limit(pageSize))
	if err != nil {
		return nil, err
	}
	var users []*model.User
	for cursor.Next(context.Background()) {
		doc := bson.NewDocument()
		err := cursor.Decode(doc)
		if err != nil {
			log.Fatal(err)
		}
		newUser := m.docToUser(doc)
		users = append(users, newUser)
	}
	return users, nil

}

func (m *mongodbUserRepository) UpdateQuantityKudos(userName string, quantity int32, updateDate time.Time) error {

	_, err := m.Collection.UpdateOne(context.Background(),
		bson.NewDocument(
			bson.EC.String("user_name", userName),
		),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$set",
				bson.EC.Int32("quantity_kudos", quantity),
				bson.EC.DateTime("updated_at", updateDate.Unix()),
			),
		),
	)
	if err != nil {
		return err
	}
	return nil

}
