package word

import (
	"context"
	"errors"
	"github.com/Go-roro/wordrop/internal/common"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "words"

type Repository struct {
	collection *mongo.Collection
}

func NewWordRepo(db *mongo.Database) *Repository {
	return &Repository{
		collection: db.Collection(collectionName),
	}
}

func (r *Repository) SaveWord(word *Word) (*Word, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	word.ID = primitive.NewObjectID()
	now := time.Now()
	word.CreatedAt = now
	word.UpdatedAt = now

	result, err := r.collection.InsertOne(ctx, word)
	if err != nil {
		return nil, err
	}

	word.ID = result.InsertedID.(primitive.ObjectID)
	log.Printf("Word %s saved with ID: %s\n", word.Text, word.ID)
	return word, nil
}

func (r *Repository) FindById(id string) (*Word, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectIDFromHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Fail to convert to ObjectID from id: %s", id)
		return nil, err
	}

	result := r.collection.FindOne(ctx, bson.M{"_id": objectIDFromHex})
	findWord := &Word{}
	if err := result.Decode(findWord); err != nil {
		log.Printf("Word with ID: %s not found.", id)
		return nil, err
	}

	return findWord, nil
}

func (r *Repository) UpdateWord(word *Word) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	target := bson.M{"_id": word.ID}
	update := bson.M{
		"$set": bson.M{
			"text":            word.Text,
			"english_meaning": word.EnglishMeaning,
			"korean_meaning":  word.KoreanMeanings,
			"description":     word.Description,
			"synonyms":        word.Synonyms,
			"examples":        word.Examples,
			"is_delivered":    word.IsDelivered,
			"created_at":      word.CreatedAt,
			"updated_at":      time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, target, update)
	if err != nil {
		log.Printf("Word with ID: %s failed to update", word.ID)
		return err
	}

	return nil
}

type SearchParams struct {
	IsDelivered *bool  `json:"is_delivered"`
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
	SortBy      string `json:"sort_by"`
	SortOrder   string `json:"sort_order"`
}

const defaultSortBy = "created_at"
const defaultSortOrder = "desc"
const defaultPageSize = 10
const maxPageSize = 30

var allowedSortFields = map[string]bool{
	"created_at": true,
}

var allowedSortOrders = map[string]int{
	"asc":  1,
	"desc": -1,
}

func (r *Repository) FindWords(params *SearchParams) (*common.PageResult[*Word], error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := setupFilter(params)
	findOptions, err := setupOptions(params)
	if err != nil {
		return nil, err
	}

	pageSize := defaultPageSize
	if params.PageSize > 0 {
		pageSize = params.PageSize
	}

	if pageSize > maxPageSize {
		log.Printf("Page size %d exceeds maximum allowed %d, setting to max", pageSize, maxPageSize)
		pageSize = maxPageSize
	}

	page := 1
	if params.Page > 0 {
		page = params.Page
	}
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))

	results, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("Failed to find words with filter: %v, error: %v", filter, err)
		return nil, err
	}

	var words []*Word
	for results.Next(ctx) {
		var word Word
		if err := results.Decode(&word); err != nil {
			log.Printf("Failed to decode word: %v", err)
			return nil, err
		}
		words = append(words, &word)
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	return common.NewPageResult(words, page, int64(pageSize), total), nil
}

func setupFilter(params *SearchParams) bson.M {
	filter := bson.M{}
	if deliveredFilter := params.IsDelivered; deliveredFilter != nil {
		filter["is_delivered"] = *deliveredFilter
	}
	return filter
}

func setupOptions(params *SearchParams) (*options.FindOptions, error) {
	findOptions := options.Find()

	sortField, err := parseSortBy(params.SortBy)
	if err != nil {
		return nil, err
	}

	sortOder, err := parseSortOrder(params.SortOrder)
	if err != nil {
		return nil, err
	}

	findOptions.SetSort(bson.D{{sortField, sortOder}})
	return findOptions, nil
}

func parseSortBy(sortedBy string) (string, error) {
	if sortedBy == "" {
		return defaultSortBy, nil
	}

	if _, ok := allowedSortFields[sortedBy]; !ok {
		log.Printf("Invalid sort field: %s", sortedBy)
		return "", errors.New("invalid sort field")
	}
	return sortedBy, nil
}

func parseSortOrder(sortOrder string) (int, error) {
	if sortOrder == "" {
		return allowedSortOrders[defaultSortOrder], nil
	}

	if _, ok := allowedSortOrders[sortOrder]; !ok {
		log.Printf("Invalid sort order: %s", sortOrder)
		return 0, errors.New("invalid sort order")
	}

	return allowedSortOrders[sortOrder], nil
}
