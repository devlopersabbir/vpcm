package notes

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var ErrNoteNotFound = errors.New("note not found")

type mongoNoteRepository struct {
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

func NewMongoRepository(db *mongo.Database) NoteRepository {
	return &mongoNoteRepository{
		db:   db,
		coll: db.Collection("notes"),
	}
}

func (r *mongoNoteRepository) Create(ctx context.Context, note *Note) error {
	id, err := getNextSequence(ctx, r.db, "note_id")
	if err != nil {
		return err
	}
	note.ID = id
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	_, err = r.coll.InsertOne(ctx, note)
	return err
}

func (r *mongoNoteRepository) GetByID(ctx context.Context, id uint) (*Note, error) {
	var note Note
	err := r.coll.FindOne(ctx, bson.M{"id": id}).Decode(&note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNoteNotFound
		}
		return nil, err
	}
	return &note, nil
}

func (r *mongoNoteRepository) ListByServer(ctx context.Context, serverID uint) ([]Note, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"server_id": serverID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notesList []Note
	if err := cursor.All(ctx, &notesList); err != nil {
		return nil, err
	}
	if notesList == nil {
		notesList = []Note{}
	}
	return notesList, nil
}

func (r *mongoNoteRepository) Update(ctx context.Context, note *Note) error {
	note.UpdatedAt = time.Now()
	_, err := r.coll.ReplaceOne(ctx, bson.M{"id": note.ID}, note)
	return err
}

func (r *mongoNoteRepository) Delete(ctx context.Context, id uint) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNoteNotFound
	}
	return nil
}
