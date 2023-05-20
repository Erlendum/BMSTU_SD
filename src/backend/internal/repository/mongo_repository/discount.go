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
	"time"
)

type DiscountMongo struct {
	DiscountId   uint64 `bson:"discount_id"`
	InstrumentId uint64 `bson:"instrument_id"`
	UserId       uint64 `bson:"user_id"`
	Amount       uint64 `bson:"discount_amount"`
	Type         string `bson:"discount_type"`
	DateBegin    string `bson:"discount_date_begin"`
	DateEnd      string `bson:"discount_date_end"`
}

type DiscountMongoRepository struct {
	db *mongo.Collection
}

func NewDiscountMongoRepository(db *mongo.Collection) repository.DiscountRepository {
	return &DiscountMongoRepository{db: db}
}

func (i *DiscountMongoRepository) toMongo(discount *models.Discount) (*DiscountMongo, error) {
	discountMongo := &DiscountMongo{}
	err := copier.Copy(discountMongo, discount)
	if err != nil {
		return nil, err
	}

	discountMongo.DateBegin = discount.DateBegin.Format(LayoutDate)
	discountMongo.DateEnd = discount.DateEnd.Format(LayoutDate)
	return discountMongo, nil
}

func (i *DiscountMongoRepository) fromMongo(discountMongo *DiscountMongo) (*models.Discount, error) {
	discount := &models.Discount{}
	err := copier.Copy(discount, &discountMongo)
	if err != nil {
		return nil, err
	}

	discount.DateBegin, err = time.Parse(LayoutDate, discountMongo.DateBegin)
	if err != nil {
		return nil, err
	}
	discount.DateEnd, err = time.Parse(LayoutDate, discountMongo.DateEnd)
	if err != nil {
		return nil, err
	}
	return discount, nil
}

func (i *DiscountMongoRepository) Create(discount *models.Discount) error {
	opts := options.Find().SetSort(bson.D{{"discount_id", -1}})
	opts.SetLimit(1)
	cur, err := i.db.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	discountMaxId := new(DiscountMongo)
	err = cur.Decode(discountMaxId)
	if err != nil && cur.Current != nil {
		return err
	}
	discount.DiscountId = discountMaxId.DiscountId + 1

	discountMongo, err := i.toMongo(discount)
	if err != nil {
		return err
	}

	_, err = i.db.InsertOne(context.Background(), discountMongo)
	if err != nil {
		return err
	}
	return nil
}

func (i *DiscountMongoRepository) discountFieldToDBField(field models.DiscountField) (string, error) {
	switch field {
	case models.DiscountFieldInstrumentId:
		return "instrument_id", nil
	case models.DiscountFieldUserId:
		return "user_id", nil
	case models.DiscountFieldAmount:
		return "discount_amount", nil
	case models.DiscountFieldType:
		return "discount_type", nil
	case models.DiscountFieldDateBegin:
		return "discount_date_begin", nil
	case models.DiscountFieldDateEnd:
		return "discount_date_end", nil
	}
	return "", repositoryErrors.InvalidField
}

func (i *DiscountMongoRepository) Update(id uint64, fieldsToUpdate models.DiscountFieldsToUpdate) error {
	filter := bson.D{{"discount_id", id}}
	if fieldsToUpdate[models.DiscountFieldDateBegin] != nil {
		fieldsToUpdate[models.DiscountFieldDateBegin] = fieldsToUpdate[models.DiscountFieldDateBegin].(time.Time).Format(LayoutDate)
	}
	if fieldsToUpdate[models.DiscountFieldDateEnd] != nil {
		fieldsToUpdate[models.DiscountFieldDateEnd] = fieldsToUpdate[models.DiscountFieldDateEnd].(time.Time).Format(LayoutDate)
	}
	updateFields := make(map[string]any, len(fieldsToUpdate))
	for key, value := range fieldsToUpdate {
		field, err := i.discountFieldToDBField(key)
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

func (i *DiscountMongoRepository) Delete(id uint64) error {
	_, err := i.db.DeleteOne(context.Background(), bson.M{"discount_id": id})
	return err
}

func (i *DiscountMongoRepository) Get(id uint64) (*models.Discount, error) {
	cur, err := i.db.Find(context.Background(), bson.M{"discount_id": id})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	if cur.Current == nil {
		return nil, repositoryErrors.ObjectDoesNotExists
	}
	discountMongo := new(DiscountMongo)
	err = cur.Decode(discountMongo)
	if err != nil {
		return nil, err
	}
	user, err := i.fromMongo(discountMongo)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *DiscountMongoRepository) GetList() ([]models.Discount, error) {
	opts := options.Find().SetSort(bson.D{{"discount_id", 1}})
	cur, err := i.db.Find(context.Background(), bson.M{}, opts)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	var discountsMongo []DiscountMongo
	var discounts []models.Discount

	for cur.Next(context.Background()) {
		discountMongo := new(DiscountMongo)
		err := cur.Decode(discountMongo)
		if err != nil {
			return nil, err
		}

		discountsMongo = append(discountsMongo, *discountMongo)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	for k := range discountsMongo {
		discount, err := i.fromMongo(&discountsMongo[k])
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, *discount)
	}

	return discounts, nil
}

func (i *DiscountMongoRepository) GetSpecificList(instrumentId uint64, userId uint64) ([]models.Discount, error) {
	opts := options.Find().SetSort(bson.D{{"discount_id", 1}})
	cur, err := i.db.Find(context.Background(), bson.M{"instrument_id": instrumentId, "user_id": userId}, opts)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	var discountsMongo []DiscountMongo
	var discounts []models.Discount

	for cur.Next(context.Background()) {
		discountMongo := new(DiscountMongo)
		err := cur.Decode(discountMongo)
		if err != nil {
			return nil, err
		}

		discountsMongo = append(discountsMongo, *discountMongo)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	for k := range discountsMongo {
		discount, err := i.fromMongo(&discountsMongo[k])
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, *discount)
	}

	return discounts, nil
}
