package repo

import (
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

func (r *MongoRepo) SaveBulkMatches(ctx context.Context, matches []Match) error {
	var mongoModels []mongo.WriteModel

	for _, match := range matches {
		cursor := r.collection.FindOne(ctx, bson.M{"external_id": match.ExternalID})
		var matchSaved Match
		err := cursor.Decode(&matchSaved)
		if err == nil && matchSaved.LoadState != "without_games" {
			continue
		}

		filter := bson.M{"external_id": match.ExternalID}
		update := bson.M{"$set": match}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		mongoModels = append(mongoModels, model)
	}

	if len(mongoModels) == 0 {
		return nil
	}

	_, err := r.collection.BulkWrite(ctx, mongoModels)
	if err != nil {
		return err
	}
	return nil
}
func (r *MongoRepo) SaveMatch(ctx context.Context, match Match) error {
	filter := bson.M{"external_id": match.ExternalID}
	update := bson.M{"$set": match}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *MongoRepo) GetMatches(ctx context.Context, filter bson.M) ([]Match, error) {
	var matches []Match
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

func (r *MongoRepo) GetTeamByName(ctx context.Context, name string) (Team, error) {
	var team Team
	filter := primitive.M{"name": name}
	err := r.collection.FindOne(ctx, filter).Decode(&team)
	return team, err
}

func (r *MongoRepo) SaveTeamByName(ctx context.Context, team Team) error {
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

func (r *MongoRepo) SaveBulkGames(ctx context.Context, games []Game) error {
	var wModel []mongo.WriteModel

	for _, game := range games {
		filter := bson.M{"external_id": game.ExternalID}
		update := bson.M{"$set": game}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		wModel = append(wModel, model)
	}
	_, err := r.collection.BulkWrite(ctx, wModel)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepo) GetGames(ctx context.Context, filter bson.M) ([]Game, error) {
	var games []Game
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &games); err != nil {
		return nil, err
	}
	return games, nil
}

func (r *MongoRepo) GetPlayerByExternalID(ctx context.Context, externalID string) (Player, error) {
	var player Player
	filter := primitive.M{"external_id": externalID}
	err := r.collection.FindOne(ctx, filter).Decode(&player)
	return player, err
}

func (r *MongoRepo) SaveBulkPlayers(ctx context.Context, players []Player) error {
	var wModel []mongo.WriteModel

	for _, player := range players {
		filter := bson.M{"external_id": player.ExternalID}
		update := bson.M{"$set": player}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)
		wModel = append(wModel, model)
	}
	_, err := r.collection.BulkWrite(ctx, wModel)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepo) SaveGame(ctx context.Context, game Game) error {
	filter := bson.M{"external_id": game.ExternalID}
	update := bson.M{"$set": game}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *MongoRepo) UpdateGameByExternalID(ctx context.Context, externalID string, updateData bson.M) error {
	filter := bson.M{"external_id": externalID}
	update := bson.M{"$set": updateData}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepo) SaveFrame(ctx context.Context, frame Frame) error {
	filter := bson.M{"game_id": frame.GameID, "timestamp": frame.TimeStamp}
	update := bson.M{"$set": frame}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}
