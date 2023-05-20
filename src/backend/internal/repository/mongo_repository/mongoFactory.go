package mongo_repository

import (
	"backend/config"
	"backend/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	LayoutDate = "2006-01-02 15:04:05 -0700 MST"
)

type MongoRepositoryFields struct {
	Db     *mongo.Database
	Config config.Config
}

func CreateMongoRepositoryFields(fileName, filePath string) *MongoRepositoryFields {
	fields := new(MongoRepositoryFields)
	err := fields.Config.ParseConfig(fileName, filePath)
	if err != nil {
		return nil
	}
	fields.Db, err = fields.Config.Mongo.InitDB()
	if err != nil {
		return nil
	}
	return fields
}

func CreateInstrumentMongoRepository(fields *MongoRepositoryFields) repository.InstrumentRepository {
	return NewInstrumentMongoRepository(fields.Db.Collection("instruments"))
}

func CreateUserMongoRepository(fields *MongoRepositoryFields) repository.UserRepository {
	return NewUserMongoRepository(fields.Db.Collection("users"))
}

func CreateComparisonListMongoRepository(fields *MongoRepositoryFields) repository.ComparisonListRepository {
	return NewComparisonListMongoRepository(fields.Db)
}

func CreateOrderMongoRepository(fields *MongoRepositoryFields) repository.OrderRepository {
	return NewOrderMongoRepository(fields.Db)
}

func CreateDiscountMongoRepository(fields *MongoRepositoryFields) repository.DiscountRepository {
	return NewDiscountMongoRepository(fields.Db.Collection("discounts"))
}
