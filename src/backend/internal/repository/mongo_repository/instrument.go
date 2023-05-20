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

type InstrumentMongo struct {
	InstrumentId uint64 `bson:"instrument_id"`
	Name         string `bson:"instrument_name"`
	Price        uint64 `bson:"instrument_price"`
	Material     string `bson:"instrument_material"`
	Type         string `bson:"instrument_type"`
	Brand        string `bson:"instrument_brand"`
	Img          string `bson:"instrument_img"`
}
type InstrumentMongoRepository struct {
	db *mongo.Collection
}

func NewInstrumentMongoRepository(db *mongo.Collection) repository.InstrumentRepository {
	return &InstrumentMongoRepository{db: db}
}

func (i *InstrumentMongoRepository) Create(instrument *models.Instrument) error {
	opts := options.Find().SetSort(bson.D{{"instrument_id", -1}})
	opts.SetLimit(1)
	cur, err := i.db.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	instrumentMaxId := new(InstrumentMongo)
	err = cur.Decode(instrumentMaxId)
	if err != nil && cur.Current != nil {
		return err
	}
	instrument.InstrumentId = instrumentMaxId.InstrumentId + 1

	instrumentMongo := &InstrumentMongo{}
	err = copier.Copy(instrumentMongo, instrument)
	if err != nil {
		return err
	}

	_, err = i.db.InsertOne(context.Background(), instrumentMongo)
	if err != nil {
		return err
	}
	return nil
}

func (i *InstrumentMongoRepository) instrumentFieldToDBField(field models.InstrumentField) (string, error) {
	switch field {
	case models.InstrumentFieldName:
		return "instrument_name", nil
	case models.InstrumentFieldPrice:
		return "instrument_price", nil
	case models.InstrumentFieldMaterial:
		return "instrument_material", nil
	case models.InstrumentFieldType:
		return "instrument_type", nil
	case models.InstrumentFieldBrand:
		return "instrument_brand", nil
	case models.InstrumentFieldImg:
		return "instrument_img", nil
	}
	return "", repositoryErrors.InvalidField
}

func (i *InstrumentMongoRepository) Update(id uint64, fieldsToUpdate models.InstrumentFieldsToUpdate) error {
	filter := bson.D{{"instrument_id", id}}
	updateFields := make(map[string]any, len(fieldsToUpdate))
	for key, value := range fieldsToUpdate {
		field, err := i.instrumentFieldToDBField(key)
		if err != nil {
			return err
		}
		updateFields[field] = value
	}
	update := queries.CreateMongoUpdateQuery(updateFields)

	_, err := i.db.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *InstrumentMongoRepository) Delete(id uint64) error {
	_, err := i.db.DeleteOne(context.Background(), bson.M{"instrument_id": id})
	return err
}

func (i *InstrumentMongoRepository) Get(id uint64) (*models.Instrument, error) {
	cur, err := i.db.Find(context.Background(), bson.M{"instrument_id": id})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	if cur.Current == nil {
		return nil, repositoryErrors.ObjectDoesNotExists
	}
	instrumentMongo := new(InstrumentMongo)
	err = cur.Decode(instrumentMongo)
	if err != nil {
		return nil, err
	}
	instrument := &models.Instrument{}
	err = copier.Copy(instrument, &instrumentMongo)
	if err != nil {
		return nil, err
	}
	return instrument, nil
}

func (i *InstrumentMongoRepository) GetList() ([]models.Instrument, error) {
	opts := options.Find().SetSort(bson.D{{"instrument_id", 1}})
	cur, err := i.db.Find(context.Background(), bson.M{}, opts)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	var instrumentsMongo []InstrumentMongo
	var instruments []models.Instrument

	for cur.Next(context.Background()) {
		instrumentMongo := new(InstrumentMongo)
		err := cur.Decode(instrumentMongo)
		if err != nil {
			return nil, err
		}

		instrumentsMongo = append(instrumentsMongo, *instrumentMongo)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	for i := range instrumentsMongo {
		instrument := &models.Instrument{}
		err = copier.Copy(instrument, &instrumentsMongo[i])
		if err != nil {
			return nil, err
		}
		instruments = append(instruments, *instrument)
	}

	if instruments == nil {
		return nil, repositoryErrors.ObjectDoesNotExists
	}
	return instruments, nil
}
