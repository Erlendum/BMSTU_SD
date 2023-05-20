package mongo_repository

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/pkg/queries"
	"backend/internal/repository"
	"context"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ComparisonListMongo struct {
	ComparisonListId uint64 `bson:"comparisonList_id"`
	UserId           uint64 `bson:"user_id"`
	TotalPrice       uint64 `bson:"comparisonList_total_price"`
	Amount           uint64 `bson:"comparisonList_amount"`
}

type ComparisonsListsInstrumentsMongo struct {
	ComparisonListsInstrumentsId uint64 `bson:"comparisonLists_instruments_id"`
	ComparisonListId             uint64 `bson:"comparisonList_id"`
	InstrumentId                 uint64 `bson:"instrument_id"`
}

type ComparisonListMongoRepository struct {
	db *mongo.Database
}

func NewComparisonListMongoRepository(db *mongo.Database) repository.ComparisonListRepository {
	return &ComparisonListMongoRepository{db: db}
}

func (i *ComparisonListMongoRepository) Create(comparisonList *models.ComparisonList) error {
	opts := options.Find().SetSort(bson.D{{"comparisonList_id", -1}})
	opts.SetLimit(1)
	cur, err := i.db.Collection("comparisonLists").Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	comparisonListMaxId := new(ComparisonListMongo)
	err = cur.Decode(comparisonListMaxId)
	if err != nil && cur.Current != nil {
		return err
	}
	comparisonList.ComparisonListId = comparisonListMaxId.ComparisonListId + 1

	comparisonListMongo := &ComparisonListMongo{}
	err = copier.Copy(comparisonListMongo, comparisonList)
	if err != nil {
		return err
	}

	_, err = i.db.Collection("comparisonLists").InsertOne(context.Background(), comparisonListMongo)
	if err != nil {
		return err
	}
	return nil
}

func (i *ComparisonListMongoRepository) AddInstrument(id uint64, instrumentId uint64) error {
	opts := options.Find().SetSort(bson.D{{"comparisonLists_instruments_id", -1}})
	opts.SetLimit(1)
	cur, err := i.db.Collection("comparisonLists_instruments").Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())

	comparisonListsInstrumentsMaxId := new(ComparisonsListsInstrumentsMongo)
	err = cur.Decode(comparisonListsInstrumentsMaxId)
	if err != nil && cur.Current != nil {
		return err
	}

	comparisonListsInstrumentsMongo := &ComparisonsListsInstrumentsMongo{}

	comparisonListsInstrumentsMongo.ComparisonListsInstrumentsId = comparisonListsInstrumentsMaxId.ComparisonListsInstrumentsId + 1
	comparisonListsInstrumentsMongo.ComparisonListId = id
	comparisonListsInstrumentsMongo.InstrumentId = instrumentId

	_, err = i.db.Collection("comparisonLists_instruments").InsertOne(context.Background(), comparisonListsInstrumentsMongo)
	if err != nil {
		return err
	}
	return nil
}

func (i *ComparisonListMongoRepository) DeleteInstrument(id uint64, instrumentId uint64) error {
	_, err := i.db.Collection("comparisonLists_instruments").DeleteOne(context.Background(), bson.M{"comparisonList_id": id, "instrument_id": instrumentId})
	return err
}

func (i *ComparisonListMongoRepository) Get(userId uint64) (*models.ComparisonList, error) {
	cur, err := i.db.Collection("comparisonLists").Find(context.Background(), bson.M{"user_id": userId})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	if cur.Current == nil {
		return nil, repositoryErrors.ObjectDoesNotExists
	}
	comparisonListMongo := new(ComparisonListMongo)
	err = cur.Decode(comparisonListMongo)
	if err != nil {
		return nil, err
	}
	comparisonList := &models.ComparisonList{}
	err = copier.Copy(comparisonList, &comparisonListMongo)
	if err != nil {
		return nil, err
	}
	return comparisonList, nil
}

func (i *ComparisonListMongoRepository) GetUser(id uint64) (*models.User, error) {
	cur, err := i.db.Collection("comparisonLists").Find(context.Background(), bson.M{"comparisonList_id": id})
	if err != nil {
		return nil, err
	}
	cur.Next(context.Background())
	if cur.Current == nil {
		return nil, repositoryErrors.ObjectDoesNotExists
	}
	comparisonListMongo := new(ComparisonListMongo)
	err = cur.Decode(comparisonListMongo)
	if err != nil {
		return nil, err
	}
	userId := comparisonListMongo.UserId
	cur.Close(context.Background())

	cur, err = i.db.Collection("users").Find(context.Background(), bson.M{"user_id": userId})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	if cur.Current == nil {
		return nil, repositoryErrors.ObjectDoesNotExists
	}
	userMongo := new(UserMongo)
	err = cur.Decode(userMongo)
	if err != nil {
		return nil, err
	}
	user := &models.User{}
	err = copier.Copy(user, &userMongo)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *ComparisonListMongoRepository) GetInstruments(userId uint64) ([]models.Instrument, error) {
	cur, err := i.db.Collection("comparisonLists").Find(context.Background(), bson.M{"user_id": userId})
	if err != nil {
		return nil, err
	}

	cur.Next(context.Background())
	if cur.Current == nil {
		return nil, repositoryErrors.ObjectDoesNotExists
	}
	comparisonListMongo := new(ComparisonListMongo)
	err = cur.Decode(comparisonListMongo)
	if err != nil {
		return nil, err
	}
	comparisonListId := comparisonListMongo.ComparisonListId
	cur.Close(context.Background())

	opts := options.Find().SetSort(bson.D{{"instrument_id", 1}})
	cur, err = i.db.Collection("comparisonLists_instruments").Find(context.Background(), bson.M{"comparisonList_id": comparisonListId}, opts)
	if err != nil {
		return nil, err
	}

	instrumentsIds := make([]uint64, 0)
	for cur.Next(context.Background()) {
		comparisonListsInstrumentsMongo := new(ComparisonsListsInstrumentsMongo)
		err := cur.Decode(comparisonListsInstrumentsMongo)
		if err != nil {
			return nil, err
		}

		instrumentsIds = append(instrumentsIds, comparisonListsInstrumentsMongo.InstrumentId)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.Background())

	var instrumentsMongo []InstrumentMongo
	for k := range instrumentsIds {
		cur, err := i.db.Collection("instruments").Find(context.Background(), bson.M{"instrument_id": instrumentsIds[k]})
		if err != nil {
			return nil, err
		}
		cur.Next(context.Background())
		if cur.Current == nil {
			return nil, repositoryErrors.ObjectDoesNotExists
		}
		instrumentMongo := new(InstrumentMongo)
		err = cur.Decode(instrumentMongo)
		if err != nil {
			return nil, err
		}
		instrumentsMongo = append(instrumentsMongo, *instrumentMongo)
		cur.Close(context.Background())
	}

	var instruments []models.Instrument
	for i := range instrumentsMongo {
		instrument := &models.Instrument{}
		err = copier.Copy(instrument, &instrumentsMongo[i])
		if err != nil {
			return nil, err
		}
		instruments = append(instruments, *instrument)
	}

	return instruments, nil
}

func (i *ComparisonListMongoRepository) comparisonListFieldToDBField(field models.ComparisonListField) (string, error) {
	switch field {
	case models.ComparisonListFieldUserId:
		return "comparisonList_user_id", nil
	case models.ComparisonListFieldTotalPrice:
		return "comparisonList_total_price", nil
	case models.ComparisonListFieldAmount:
		return "comparisonList_amount", nil
	}
	return "", repositoryErrors.InvalidField
}

func (i *ComparisonListMongoRepository) Update(id uint64, fieldsToUpdate models.ComparisonListFieldsToUpdate) error {
	filter := bson.D{{"comparisonList_id", id}}
	updateFields := make(map[string]any, len(fieldsToUpdate))
	for key, value := range fieldsToUpdate {
		field, err := i.comparisonListFieldToDBField(key)
		if err != nil {
			return err
		}
		updateFields[field] = value
	}
	update := queries.CreateMongoUpdateQuery(updateFields)

	_, err := i.db.Collection("comparisonLists").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *ComparisonListMongoRepository) Clear(id uint64) error {
	_, err := i.db.Collection("comparisonLists_instruments").DeleteMany(context.Background(), bson.M{"comparisonList_id": id})
	return err
}
