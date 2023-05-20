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

type OrderMongo struct {
	OrderId uint64 `bson:"order_id"`
	UserId  uint64 `bson:"user_id"`
	Price   uint64 `bson:"order_price"`
	Time    string `bson:"order_time"`
	Status  string `bson:"order_status"`
}

type OrderElementMongo struct {
	OrderElementId uint64 `bson:"order_element_id"`
	InstrumentId   uint64 `bson:"instrument_id"`
	OrderId        uint64 `bson:"order_id"`
	Amount         uint64 `bson:"order_element_amount"`
	Price          uint64 `bson:"order_element_price"`
}

type OrderMongoRepository struct {
	db *mongo.Database
}

func NewOrderMongoRepository(db *mongo.Database) repository.OrderRepository {
	return &OrderMongoRepository{db: db}
}

func (i *OrderMongoRepository) toMongo(order *models.Order) (*OrderMongo, error) {
	orderMongo := &OrderMongo{}
	err := copier.Copy(orderMongo, order)
	if err != nil {
		return nil, err
	}

	orderMongo.Time = order.Time.Format(LayoutDate)
	return orderMongo, nil
}

func (i *OrderMongoRepository) fromMongo(orderMongo *OrderMongo) (*models.Order, error) {
	order := &models.Order{}
	err := copier.Copy(order, &orderMongo)
	if err != nil {
		return nil, err
	}
	order.Time, err = time.Parse(LayoutDate, orderMongo.Time)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (i *OrderMongoRepository) Create(order *models.Order) (uint64, error) {
	opts := options.Find().SetSort(bson.D{{"order_id", -1}})
	opts.SetLimit(1)
	cur, err := i.db.Collection("orders").Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return 0, err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	orderMaxId := new(OrderMongo)
	err = cur.Decode(orderMaxId)
	if err != nil && cur.Current != nil {
		return 0, err
	}
	order.OrderId = orderMaxId.OrderId + 1
	order.Status = "Created"

	orderMongo, err := i.toMongo(order)
	if err != nil {
		return 0, err
	}

	_, err = i.db.Collection("orders").InsertOne(context.Background(), orderMongo)
	if err != nil {
		return 0, err
	}

	return orderMongo.OrderId, nil
}

func (i *OrderMongoRepository) CreateOrderElement(element *models.OrderElement) error {
	opts := options.Find().SetSort(bson.D{{"order_element_id", -1}})
	opts.SetLimit(1)
	cur, err := i.db.Collection("orderElements").Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	orderElementMaxId := new(OrderElementMongo)
	err = cur.Decode(orderElementMaxId)
	if err != nil && cur.Current != nil {
		return err
	}
	element.OrderElementId = orderElementMaxId.OrderElementId + 1

	elementMongo := &OrderElementMongo{}
	err = copier.Copy(elementMongo, element)
	if err != nil {
		return err
	}

	_, err = i.db.Collection("orderElements").InsertOne(context.Background(), elementMongo)
	if err != nil {
		return err
	}
	return nil
}

func (i *OrderMongoRepository) GetList(userId uint64) ([]models.Order, error) {
	opts := options.Find().SetSort(bson.D{{"order_id", 1}})
	cur, err := i.db.Collection("orders").Find(context.Background(), bson.M{"user_id": userId}, opts)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	var ordersMongo []OrderMongo
	var orders []models.Order

	for cur.Next(context.Background()) {
		orderMongo := new(OrderMongo)
		err := cur.Decode(orderMongo)
		if err != nil {
			return nil, err
		}

		ordersMongo = append(ordersMongo, *orderMongo)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	for k := range ordersMongo {
		order, err := i.fromMongo(&ordersMongo[k])
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}

	return orders, nil
}

func (i *OrderMongoRepository) GetListForAll() ([]models.Order, error) {
	opts := options.Find().SetSort(bson.D{{"order_id", 1}})
	cur, err := i.db.Collection("orders").Find(context.Background(), bson.M{}, opts)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	var ordersMongo []OrderMongo
	var orders []models.Order

	for cur.Next(context.Background()) {
		orderMongo := new(OrderMongo)
		err := cur.Decode(orderMongo)
		if err != nil {
			return nil, err
		}

		ordersMongo = append(ordersMongo, *orderMongo)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	for k := range ordersMongo {
		order, err := i.fromMongo(&ordersMongo[k])
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}

	return orders, nil
}

func (i *OrderMongoRepository) orderFieldToDBField(field models.OrderField) (string, error) {
	switch field {
	case models.OrderFieldUserId:
		return "user_id", nil
	case models.OrderFieldPrice:
		return "order_price", nil
	case models.OrderFieldTime:
		return "order_time", nil
	case models.OrderFieldStatus:
		return "order_status", nil
	}
	return "", repositoryErrors.InvalidField
}

func (i *OrderMongoRepository) Update(id uint64, fieldsToUpdate models.OrderFieldsToUpdate) error {
	filter := bson.D{{"order_id", id}}
	updateFields := make(map[string]any, len(fieldsToUpdate))
	for key, value := range fieldsToUpdate {
		field, err := i.orderFieldToDBField(key)
		if err != nil {
			return err
		}
		updateFields[field] = value
	}
	update := queries.CreateMongoUpdateQuery(updateFields)

	_, err := i.db.Collection("orders").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (i *OrderMongoRepository) GetOrderElements(id uint64) ([]models.OrderElement, error) {
	opts := options.Find().SetSort(bson.D{{"instrument_id", 1}})
	cur, err := i.db.Collection("orderElements").Find(context.Background(), bson.M{"order_id": id}, opts)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	var orderElementsMongo []OrderElementMongo
	var orderElements []models.OrderElement

	for cur.Next(context.Background()) {
		orderElementMongo := new(OrderElementMongo)
		err := cur.Decode(orderElementMongo)
		if err != nil {
			return nil, err
		}

		orderElementsMongo = append(orderElementsMongo, *orderElementMongo)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	for i := range orderElementsMongo {
		orderElement := &models.OrderElement{}
		err = copier.Copy(orderElement, &orderElementsMongo[i])
		if err != nil {
			return nil, err
		}
		orderElements = append(orderElements, *orderElement)
	}

	return orderElements, nil
}
