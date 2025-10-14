package repo

import (
	"betty/science/app/league_of_legends/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	collection *mongo.Collection
	db         *mongo.Database
}

func NewMongoRepo(collectionName string, db *mongo.Database) *MongoRepo {
	collection := db.Collection(collectionName)
	return &MongoRepo{collection: collection, db: db}
}

func (r *MongoRepo) SaveBulkMatches(ctx context.Context, matches []models.Match) error {
	var models []mongo.WriteModel

	for _, match := range matches {
		filter := bson.M{"external_id": match.ExternalID}
		update := bson.M{"$set": match}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		models = append(models, model)
	}
	_, err := r.collection.BulkWrite(ctx, models)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepo) GetMatches(ctx context.Context, filter bson.M) ([]models.Match, error) {
	var matches []models.Match
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &matches); err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *MongoRepo) UpdateTeamTournaments(ctx context.Context, id primitive.ObjectID, tournament string) error {
	filter := primitive.M{"_id": id}
	update := primitive.M{"$push": primitive.M{"tournaments": tournament}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepo) GetTeamByName(ctx context.Context, name string) (models.Team, error) {
	var team models.Team
	filter := primitive.M{"name": name}
	err := r.collection.FindOne(ctx, filter).Decode(&team)
	return team, err
}

func (r *MongoRepo) SaveTeamByName(ctx context.Context, team models.Team) error {
	filter := bson.M{"name": team.Name}
	update := bson.M{"$set": team}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *MongoRepo) UpdateTeamExternalID(ctx context.Context, id primitive.ObjectID, externalID string) error {
	filter := primitive.M{"_id": id}
	update := primitive.M{"$set": primitive.M{"external_id": externalID}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepo) SaveBulkGames(ctx context.Context, games []models.Game) error {
	var models []mongo.WriteModel

	for _, game := range games {
		filter := bson.M{"external_id": game.ExternalID}
		update := bson.M{"$set": game}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		models = append(models, model)
	}
	_, err := r.collection.BulkWrite(ctx, models)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepo) GetGames(ctx context.Context, filter bson.M) ([]models.Game, error) {
	var games []models.Game
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}
	return games, nil
}
