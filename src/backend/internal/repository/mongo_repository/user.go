package mongo_repository

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/repository"
	"context"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserMongo struct {
	UserId    uint64            `bson:"user_id"`
	Login     string            `bson:"user_login"`
	Password  string            `bson:"user_password"`
	Fio       string            `bson:"user_fio"`
	DateBirth string            `bson:"user_date_birth"`
	Gender    models.UserGender `bson:"user_gender"`
	IsAdmin   string            `bson:"user_is_admin"`
}

type UserMongoRepository struct {
	db *mongo.Collection
}

func NewUserMongoRepository(db *mongo.Collection) repository.UserRepository {
	return &UserMongoRepository{db: db}
}

func (i *UserMongoRepository) toMongo(user *models.User) (*UserMongo, error) {
	userMongo := &UserMongo{}
	err := copier.Copy(userMongo, user)
	if err != nil {
		return nil, err
	}
	userMongo.IsAdmin = "false"
	if user.IsAdmin {
		userMongo.IsAdmin = "true"
	}
	userMongo.DateBirth = user.DateBirth.Format(LayoutDate)
	return userMongo, nil
}

func (i *UserMongoRepository) fromMongo(userMongo *UserMongo) (*models.User, error) {
	user := &models.User{}
	err := copier.Copy(user, &userMongo)
	if err != nil {
		return nil, err
	}
	user.IsAdmin = false
	if userMongo.IsAdmin == "true" {
		user.IsAdmin = true
	}
	user.DateBirth, err = time.Parse(LayoutDate, userMongo.DateBirth)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserMongoRepository) Create(user *models.User) error {
	opts := options.Find().SetSort(bson.D{{"user_id", -1}})
	opts.SetLimit(1)
	cur, err := i.db.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	userMaxId := new(UserMongo)
	err = cur.Decode(userMaxId)
	if err != nil && cur.Current != nil {
		return err
	}
	user.UserId = userMaxId.UserId + 1

	userMongo, err := i.toMongo(user)
	if err != nil {
		return err
	}

	_, err = i.db.InsertOne(context.Background(), userMongo)
	if err != nil {
		return err
	}
	return nil
}

func (i *UserMongoRepository) Get(login string) (*models.User, error) {
	cur, err := i.db.Find(context.Background(), bson.M{"user_login": login})
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
	user, err := i.fromMongo(userMongo)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserMongoRepository) GetById(id uint64) (*models.User, error) {
	cur, err := i.db.Find(context.Background(), bson.M{"user_id": id})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	cur.Next(context.Background())
	userMongo := new(UserMongo)
	err = cur.Decode(userMongo)
	if err != nil {
		return nil, err
	}
	user, err := i.fromMongo(userMongo)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserMongoRepository) GetList() ([]models.User, error) {
	opts := options.Find().SetSort(bson.D{{"user_id", 1}})
	cur, err := i.db.Find(context.Background(), bson.M{}, opts)
	defer cur.Close(context.Background())

	if err != nil {
		return nil, err
	}

	var usersMongo []UserMongo
	var users []models.User

	for cur.Next(context.Background()) {
		userMongo := new(UserMongo)
		err := cur.Decode(userMongo)
		if err != nil {
			return nil, err
		}

		usersMongo = append(usersMongo, *userMongo)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	for k := range usersMongo {
		user, err := i.fromMongo(&usersMongo[k])
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}
