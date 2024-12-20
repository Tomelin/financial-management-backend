package entity

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

type ICategory interface {
	Create(*CategoryResponse) (*CategoryResponse, error)
	Get() ([]CategoryResponse, error)
	GetById(id *string) (*CategoryResponse, error)
	Update(data *CategoryResponse) (*CategoryResponse, error)
	Delete(id *string) error
	GetByFilterMany(key string, value *string) ([]CategoryResponse, error)
	GetByFilterOne(key string, value *string) (*CategoryResponse, error)
}

type Category struct {
	Name        string `json:"name" binding:"required,min=3,max=50" firestore:"name"`
	Description string `json:"description,omitempty" firestore:"description"`
}

type CategoryResponse struct {
	ID string `json:"id" firestore:"id"`
	Category
}

func NewCategory(name, description *string) (*CategoryResponse, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	if name == nil || *name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if len(*name) > 50 {
		return nil, fmt.Errorf("name must be less than 50 characters")
	}

	return &CategoryResponse{
		ID: id.String(),
		Category: Category{
			Name:        *name,
			Description: *description,
		},
	}, nil
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (c *Category) IsEmpty(data *Category) bool {
	return data == nil || reflect.DeepEqual(*data, Category{})
}

func (c *CategoryResponse) Validate() error {
	_, err := uuid.Parse(c.ID)
	if err != nil {
		return fmt.Errorf("id is required")
	}

	return c.Category.Validate()
}

func (c *CategoryResponse) IsEmpty(data *CategoryResponse) bool {
	return data == nil || reflect.DeepEqual(*data, CategoryResponse{})
}
