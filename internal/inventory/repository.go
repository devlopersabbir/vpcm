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

// ─── Core CRUD ────────────────────────────────────────────────────────────────

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

// ─── Tags ─────────────────────────────────────────────────────────────────────

func (r *mongoServerRepository) AddTag(ctx context.Context, serverID uint, tagName string) error {
	var server Server
	if err := r.coll.FindOne(ctx, bson.M{"id": serverID}).Decode(&server); err != nil {
		return err
	}
	for _, t := range server.Tags {
		if t.Name == tagName {
			return nil
		}
	}
	_, err := r.coll.UpdateOne(ctx, bson.M{"id": serverID}, bson.M{"$push": bson.M{"tags": Tag{Name: tagName}}})
	return err
}

func (r *mongoServerRepository) RemoveTag(ctx context.Context, serverID uint, tagName string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"id": serverID}, bson.M{"$pull": bson.M{"tags": bson.M{"name": tagName}}})
	return err
}

// ─── Metadata Upserts (MongoDB stubs) ────────────────────────────────────────
// Full MongoDB embedded-document implementation can replace these stubs when
// the Mongo backend is promoted to production.

func (r *mongoServerRepository) UpsertNetwork(ctx context.Context, n *ServerNetwork) error {
	_, err := r.coll.UpdateOne(ctx,
		bson.M{"id": n.ServerID},
		bson.M{"$set": bson.M{"network": n}},
		options.UpdateOne().SetUpsert(false))
	return err
}

func (r *mongoServerRepository) UpsertHardware(ctx context.Context, h *ServerHardware) error {
	_, err := r.coll.UpdateOne(ctx,
		bson.M{"id": h.ServerID},
		bson.M{"$set": bson.M{"hardware": h}},
		options.UpdateOne().SetUpsert(false))
	return err
}

func (r *mongoServerRepository) UpsertOS(ctx context.Context, o *ServerOS) error {
	_, err := r.coll.UpdateOne(ctx,
		bson.M{"id": o.ServerID},
		bson.M{"$set": bson.M{"os": o}},
		options.UpdateOne().SetUpsert(false))
	return err
}

// ─── Software ─────────────────────────────────────────────────────────────────

func (r *mongoServerRepository) ReplaceSoftware(ctx context.Context, serverID uint, software []Software) error {
	_, err := r.coll.UpdateOne(ctx,
		bson.M{"id": serverID},
		bson.M{"$set": bson.M{"software": software}},
		options.UpdateOne().SetUpsert(false))
	return err
}

func (r *mongoServerRepository) GetSoftware(ctx context.Context, serverID uint) ([]Software, error) {
	var server struct {
		Software []Software `bson:"software"`
	}
	if err := r.coll.FindOne(ctx, bson.M{"id": serverID}).Decode(&server); err != nil {
		return []Software{}, nil
	}
	if server.Software == nil {
		return []Software{}, nil
	}
	return server.Software, nil
}

// ─── Joined Views (MongoDB) ───────────────────────────────────────────────────

func (r *mongoServerRepository) GetServerView(ctx context.Context, id uint) (*ServerView, error) {
	var doc struct {
		Server   `bson:",inline"`
		Network  *ServerNetwork  `bson:"network"`
		Hardware *ServerHardware `bson:"hardware"`
		OS       *ServerOS       `bson:"os"`
		Software []Software      `bson:"software"`
	}
	if err := r.coll.FindOne(ctx, bson.M{"id": id}).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}
	return mongoDocToView(doc.Server, doc.Network, doc.Hardware, doc.OS, doc.Software), nil
}

func (r *mongoServerRepository) GetServerViewByUUID(ctx context.Context, uuid string) (*ServerView, error) {
	var doc struct {
		Server   `bson:",inline"`
		Network  *ServerNetwork  `bson:"network"`
		Hardware *ServerHardware `bson:"hardware"`
		OS       *ServerOS       `bson:"os"`
		Software []Software      `bson:"software"`
	}
	if err := r.coll.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}
	return mongoDocToView(doc.Server, doc.Network, doc.Hardware, doc.OS, doc.Software), nil
}

func (r *mongoServerRepository) ListServerViews(ctx context.Context) ([]ServerView, error) {
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []struct {
		Server   `bson:",inline"`
		Network  *ServerNetwork  `bson:"network"`
		Hardware *ServerHardware `bson:"hardware"`
		OS       *ServerOS       `bson:"os"`
		Software []Software      `bson:"software"`
	}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	views := make([]ServerView, 0, len(docs))
	for _, d := range docs {
		views = append(views, *mongoDocToView(d.Server, d.Network, d.Hardware, d.OS, d.Software))
	}
	return views, nil
}

func mongoDocToView(s Server, n *ServerNetwork, h *ServerHardware, o *ServerOS, sw []Software) *ServerView {
	if sw == nil {
		sw = []Software{}
	}
	return &ServerView{
		ID:        s.ID,
		UUID:      s.UUID,
		Name:      s.Name,
		Host:      s.Host,
		Port:      s.Port,
		Username:  s.Username,
		AuthType:  s.AuthType,
		Provider:  s.Provider,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		LastSeen:  s.LastSeen,
		Tags:      s.Tags,
		Network:   n,
		Hardware:  h,
		OS:        o,
		Software:  sw,
	}
}

// ─── Misc ─────────────────────────────────────────────────────────────────────

func (r *mongoServerRepository) Flush(ctx context.Context) error {
	if _, err := r.coll.DeleteMany(ctx, bson.M{}); err != nil {
		return err
	}
	_, _ = r.db.Collection("connection_logs").DeleteMany(ctx, bson.M{})
	_, _ = r.db.Collection("counters").DeleteMany(ctx, bson.M{})
	return nil
}

// ─── Connection Logs ──────────────────────────────────────────────────────────

func (r *mongoServerRepository) CreateConnectionLog(ctx context.Context, s *ConnectionLog) error {
	id, err := getNextSequence(ctx, r.db, "connection_log_id")
	if err != nil {
		return err
	}
	s.ID = id
	_, err = r.db.Collection("connection_logs").InsertOne(ctx, s)
	return err
}

func (r *mongoServerRepository) UpdateConnectionLog(ctx context.Context, s *ConnectionLog) error {
	_, err := r.db.Collection("connection_logs").ReplaceOne(ctx, bson.M{"id": s.ID}, s)
	return err
}

func (r *mongoServerRepository) GetConnectionLogs(ctx context.Context, serverID uint) ([]ConnectionLog, error) {
	opts := options.Find().SetSort(bson.D{{Key: "logged_in_at", Value: -1}})
	cursor, err := r.db.Collection("connection_logs").Find(ctx, bson.M{"server_id": serverID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []ConnectionLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	if logs == nil {
		logs = []ConnectionLog{}
	}
	return logs, nil
}
