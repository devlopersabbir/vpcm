package inventory

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mongoServerRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

type Counter struct {
	ID    string `bson:"_id"`
	Value uint   `bson:"value"`
}

func getNextSequence(ctx context.Context, db *mongo.Database, sequenceName string) (uint, error) {
	coll := db.Collection("counters")
	filter := bson.M{"_id": sequenceName}
	update := bson.M{"$inc": bson.M{"value": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var counter Counter
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&counter)
	if err != nil {
		return 0, err
	}
	return counter.Value, nil
}

func NewMongoRepository(db *mongo.Database) ServerRepository {
	return &mongoServerRepository{
		db:   db,
		coll: db.Collection("servers"),
	}
}

func (r *mongoServerRepository) Create(ctx context.Context, server *Server) error {
	id, err := getNextSequence(ctx, r.db, "server_id")
	if err != nil {
		return err
	}
	server.ID = id
	server.CreatedAt = time.Now()
	server.UpdatedAt = time.Now()

	_, err = r.coll.InsertOne(ctx, server)
	return err
}

func (r *mongoServerRepository) GetByID(ctx context.Context, id uint) (*Server, error) {
	var server Server
	err := r.coll.FindOne(ctx, bson.M{"id": id}).Decode(&server)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}
	return &server, nil
}

func (r *mongoServerRepository) GetByUUID(ctx context.Context, uuid string) (*Server, error) {
	var server Server
	err := r.coll.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&server)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}
	return &server, nil
}

func (r *mongoServerRepository) List(ctx context.Context) ([]Server, error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var servers []Server
	if err := cursor.All(ctx, &servers); err != nil {
		return nil, err
	}
	if servers == nil {
		servers = []Server{}
	}
	return servers, nil
}

func (r *mongoServerRepository) Update(ctx context.Context, server *Server) error {
	server.UpdatedAt = time.Now()
	_, err := r.coll.ReplaceOne(ctx, bson.M{"id": server.ID}, server)
	return err
}

func (r *mongoServerRepository) Delete(ctx context.Context, id uint) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrServerNotFound
	}
	return nil
}

func (r *mongoServerRepository) AddTag(ctx context.Context, serverID uint, tagName string) error {
	var server Server
	err := r.coll.FindOne(ctx, bson.M{"id": serverID}).Decode(&server)
	if err != nil {
		return err
	}

	// Add unique tag
	for _, t := range server.Tags {
		if t.Name == tagName {
			return nil
		}
	}

	newTag := Tag{Name: tagName}
	_, err = r.coll.UpdateOne(ctx, bson.M{"id": serverID}, bson.M{"$push": bson.M{"tags": newTag}})
	return err
}

func (r *mongoServerRepository) RemoveTag(ctx context.Context, serverID uint, tagName string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"id": serverID}, bson.M{"$pull": bson.M{"tags": bson.M{"name": tagName}}})
	return err
}
