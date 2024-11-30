package service

import (
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/utils"
)

type CategorySvc struct {
	repo entity.ICategory
}

// func NewCategorySvcsitory(db entity.ICategory) (*CategorySvc, error) {
func NewCategorySvcsitory(db entity.ICategory) (entity.ICategory, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}
	return &CategorySvc{repo: db}, nil
}

func (c *CategorySvc) Create(cat *entity.CategoryResponse) (*entity.CategoryResponse, error) {

	res, err := c.GetByFilterOne("name", &cat.Name)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}

	if res.IsEmpty(res) {
		response, err := c.repo.Create(cat)
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	return nil, errors.New("category already exists")
}

func (c *CategorySvc) Get() ([]entity.CategoryResponse, error) {
	return c.repo.Get()
}

func (c *CategorySvc) GetById(id *string) (*entity.CategoryResponse, error) {

	if id == nil || *id == "" {
		return nil, errors.New("id cannot be empty")
	}

	err := utils.ValidateUUID(id)
	if err != nil {
		return nil, err
	}

	data, err := c.repo.GetById(id)
	if err != nil {
		if strings.Contains(err.Error(), "rpc error: code = NotFound") {
			return nil, errors.New("not found")
		}
	}
	return data, err
}

func (c *CategorySvc) GetByFilterMany(key string, value *string) ([]entity.CategoryResponse, error) {
	if key == "" {
		return nil, errors.New("key cannot be empty")
	}

	if value == nil || *value == "" {
		return nil, errors.New("value cannot be empty")
	}

	data, err := c.repo.GetByFilterMany(key, value)
	if err != nil {
		return nil, err
	}

	var items []entity.CategoryResponse
	for _, item := range data {
		if item.Category.Name == *value {
			items = append(items, item)
		}
	}

	return items, nil
}

func (c *CategorySvc) GetByFilterOne(key string, value *string) (*entity.CategoryResponse, error) {

	if key == "" {
		return nil, errors.New("key cannot be empty")
	}

	if value == nil || *value == "" {
		return nil, errors.New("value cannot be empty")
	}

	return c.repo.GetByFilterOne(key, value)
}

func (c *CategorySvc) Update(data *entity.CategoryResponse) (*entity.CategoryResponse, error) {

	if data == nil {
		return nil, errors.New("category cannot be empty")
	}

	err := utils.ValidateUUID(&data.ID)
	if err != nil {
		return nil, err
	}

	if err := data.Validate(); err != nil {
		return nil, err
	}

	return c.repo.Update(data)
}

func (c *CategorySvc) Delete(id *string) error {

	if id == nil || *id == "" {
		return errors.New("id cannot be empty")
	}

	err := utils.ValidateUUID(id)
	if err != nil {
		return err
	}

	data, err := c.GetById(id)
	if err != nil {
		return err
	}

	if data != nil {
		return c.repo.Delete(id)
	}

	return errors.New("not found")
}

func (c *CategorySvc) getFieldValueByName(v interface{}, fieldName string) (interface{}, error) {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(fieldName)
	log.Println(f.Interface())
	if !f.IsValid() {
		return nil, errors.New("field not found")
	}
	return f.Interface(), nil
}
