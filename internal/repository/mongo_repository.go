package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/M2rk13/Otus-327619/internal/config"
	"github.com/M2rk13/Otus-327619/internal/model/api"
	logmodel "github.com/M2rk13/Otus-327619/internal/model/log"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	mongoClient *mongo.Client
	redisClient *redis.Client
	dbName      string
	redisTTL    time.Duration
}

func NewMongoStore(ctx context.Context, mongoCfg config.MongoConfig, redisCfg config.RedisConfig) (*MongoStore, error) {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoCfg.URI))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       0,
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &MongoStore{
		mongoClient: mongoClient,
		redisClient: redisClient,
		dbName:      mongoCfg.Database,
		redisTTL:    redisCfg.TTL,
	}, nil
}

func (s *MongoStore) Close(ctx context.Context) {
	if s.mongoClient != nil {
		s.mongoClient.Disconnect(ctx)
	}

	if s.redisClient != nil {
		s.redisClient.Close()
	}
}

func (s *MongoStore) logChangeToRedis(action, entity, id string) {
	key := fmt.Sprintf("log:%s:%s:%s", entity, id, time.Now().Format(time.RFC3339Nano))
	value := fmt.Sprintf("Action '%s' performed on entity '%s' with ID '%s'", action, entity, id)
	s.redisClient.Set(key, value, s.redisTTL)
}

func (s *MongoStore) collection(name string) *mongo.Collection {
	return s.mongoClient.Database(s.dbName).Collection(name)
}

func (s *MongoStore) CreateRequest(req *api.Request) {
	req.Id = uuid.New().String()
	_, _ = s.collection("requests").InsertOne(context.Background(), req)
	s.logChangeToRedis("CREATE", "request", req.Id)
}

func (s *MongoStore) GetRequestByID(id string) *api.Request {
	var req api.Request
	err := s.collection("requests").FindOne(context.Background(), bson.M{"id": id}).Decode(&req)

	if err != nil {
		return nil
	}

	return &req
}

func (s *MongoStore) GetAllRequests() []*api.Request {
	var results []*api.Request
	cursor, err := s.collection("requests").Find(context.Background(), bson.M{})

	if err != nil {
		return results
	}

	defer cursor.Close(context.Background())
	_ = cursor.All(context.Background(), &results)

	return results
}

func (s *MongoStore) UpdateRequest(req *api.Request) bool {
	res, _ := s.collection("requests").UpdateOne(context.Background(), bson.M{"id": req.Id}, bson.M{"$set": req})

	if res.ModifiedCount > 0 {
		s.logChangeToRedis("UPDATE", "request", req.Id)

		return true
	}

	return false
}

func (s *MongoStore) DeleteRequest(id string) bool {
	res, _ := s.collection("requests").DeleteOne(context.Background(), bson.M{"id": id})

	if res.DeletedCount > 0 {
		s.logChangeToRedis("DELETE", "request", id)

		return true
	}

	return false
}

func (s *MongoStore) CreateResponse(resp *api.Response) {
	resp.Id = uuid.New().String()
	_, _ = s.collection("responses").InsertOne(context.Background(), resp)
	s.logChangeToRedis("CREATE", "response", resp.Id)
}

func (s *MongoStore) GetResponseByID(id string) *api.Response {
	var resp api.Response
	err := s.collection("responses").FindOne(context.Background(), bson.M{"id": id}).Decode(&resp)

	if err != nil {
		return nil
	}

	return &resp
}

func (s *MongoStore) GetAllResponses() []*api.Response {
	var results []*api.Response
	cursor, err := s.collection("responses").Find(context.Background(), bson.M{})

	if err != nil {
		return results
	}

	defer cursor.Close(context.Background())
	_ = cursor.All(context.Background(), &results)

	return results
}

func (s *MongoStore) UpdateResponse(resp *api.Response) bool {
	res, _ := s.collection("responses").UpdateOne(context.Background(), bson.M{"id": resp.Id}, bson.M{"$set": resp})

	if res.ModifiedCount > 0 {
		s.logChangeToRedis("UPDATE", "response", resp.Id)

		return true
	}

	return false
}

func (s *MongoStore) DeleteResponse(id string) bool {
	res, _ := s.collection("responses").DeleteOne(context.Background(), bson.M{"id": id})

	if res.DeletedCount > 0 {
		s.logChangeToRedis("DELETE", "response", id)

		return true
	}

	return false
}

func (s *MongoStore) CreateConversionLog(log *logmodel.ConversionLog) {
	log.Id = uuid.New().String()
	_, _ = s.collection("conversion_logs").InsertOne(context.Background(), log)
	s.logChangeToRedis("CREATE", "conversion_log", log.Id)
}

func (s *MongoStore) GetConversionLogByID(id string) *logmodel.ConversionLog {
	var log logmodel.ConversionLog
	err := s.collection("conversion_logs").FindOne(context.Background(), bson.M{"id": id}).Decode(&log)

	if err != nil {
		return nil
	}

	return &log
}

func (s *MongoStore) GetAllConversionLogs() []*logmodel.ConversionLog {
	var results []*logmodel.ConversionLog
	cursor, err := s.collection("conversion_logs").Find(context.Background(), bson.M{})

	if err != nil {
		return results
	}

	defer cursor.Close(context.Background())
	_ = cursor.All(context.Background(), &results)

	return results
}

func (s *MongoStore) UpdateConversionLog(log *logmodel.ConversionLog) bool {
	res, _ := s.collection("conversion_logs").UpdateOne(context.Background(), bson.M{"id": log.Id}, bson.M{"$set": log})

	if res.ModifiedCount > 0 {
		s.logChangeToRedis("UPDATE", "conversion_log", log.Id)

		return true
	}

	return false
}

func (s *MongoStore) DeleteConversionLog(id string) bool {
	res, _ := s.collection("conversion_logs").DeleteOne(context.Background(), bson.M{"id": id})

	if res.DeletedCount > 0 {
		s.logChangeToRedis("DELETE", "conversion_log", id)

		return true
	}

	return false
}

func (s *MongoStore) GetNewConversionRequests() []*api.Request        { return nil }
func (s *MongoStore) GetNewConversionResponses() []*api.Response      { return nil }
func (s *MongoStore) GetNewConversionLogs() []*logmodel.ConversionLog { return nil }
