package repository

import (
	"context"
	"errors"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/db"
	"google.golang.org/api/iterator"
)

type CategoryRepo struct {
	db db.FirebaseDatabaseInterface
}

func NewCategoryRepository(db db.FirebaseDatabaseInterface) (*CategoryRepo, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}
	return &CategoryRepo{db: db}, nil
}

func (c *CategoryRepo) Create(cat *entity.CategoryResponse) (*entity.CategoryResponse, error) {

	ctx := context.Background()

	_, err := c.db.Collection("categories").Doc(cat.ID).Set(ctx, cat)
	// ref, _, err := c.db.Collection("categories").Add(ctx, cat)
	if err != nil {
		return nil, err
	}

	doc, err := c.db.Collection("categories").Doc(cat.ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var category entity.CategoryResponse
	err = doc.DataTo(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (c *CategoryRepo) Get() ([]entity.CategoryResponse, error) {

	iter := c.db.Documents(context.Background(), "categories")

	defer iter.Stop()

	var documents []entity.CategoryResponse
	for {
		var user entity.CategoryResponse
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&user)
		if err != nil {
			return nil, err
		}
		user.ID = doc.Ref.ID

		documents = append(documents, user)
	}

	return documents, nil

}

func (c *CategoryRepo) GetById(id *string) (*entity.CategoryResponse, error) {

	doc, err := c.db.Collection("categories").Doc(*id).Get(context.Background())
	if err != nil {
		return nil, err
	}

	var category entity.CategoryResponse
	err = doc.DataTo(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (c *CategoryRepo) GetByFilterMany(key string, value *string) ([]entity.CategoryResponse, error) {
	iter := c.db.Documents(context.Background(), "categories")

	defer iter.Stop()

	var documents []entity.CategoryResponse
	for {
		var user entity.CategoryResponse
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&user)
		if err != nil {
			return nil, err
		}
		user.ID = doc.Ref.ID

		documents = append(documents, user)
	}

	return documents, nil
}

func (c *CategoryRepo) GetByFilterOne(key string, value *string) (*entity.CategoryResponse, error) {
	iter := c.db.Documents(context.Background(), "categories")

	defer iter.Stop()

	var categories entity.CategoryResponse
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&categories)
		if err != nil {
			return nil, err
		}

		if categories.Name == *value {
			categories.ID = doc.Ref.ID
			return &categories, nil
		}
	}

	return &categories, nil
}

func (c *CategoryRepo) Update(data *entity.CategoryResponse) (*entity.CategoryResponse, error) {
	_, err := c.db.Collection("categories").Doc(data.ID).Set(context.Background(), data)
	if err != nil {
		return nil, err
	}

	doc, err := c.db.Collection("categories").Doc(data.ID).Get(context.Background())
	if err != nil {
		return nil, err
	}

	var category entity.CategoryResponse
	err = doc.DataTo(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil

}

func (c *CategoryRepo) Delete(id *string) error {
	_, err := c.db.Collection("categories").Doc(*id).Delete(context.Background())
	if err != nil {
		return err
	}

	return nil
}
