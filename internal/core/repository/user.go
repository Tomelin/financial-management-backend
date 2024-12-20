package repository

import (
	"context"
	"errors"

	"github.com/synera-br/financial-management/src/backend/internal/core/entity"
	"github.com/synera-br/financial-management/src/backend/pkg/db"
	"github.com/synera-br/financial-management/src/backend/pkg/observability"
	"google.golang.org/api/iterator"
)

type UserRepo struct {
	db    db.FirebaseDatabaseInterface
	trace *observability.Tracer
}

func NewUserRepository(db db.FirebaseDatabaseInterface) (entity.IUser, error) {
	if db == nil {
		return nil, errors.New("db is required")
	}
	return &UserRepo{db: db}, nil
}

func (u *UserRepo) Create(ctx context.Context, user *entity.AccountUser) (*entity.AccountUser, error) {

	_, err := u.db.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return nil, err
	}

	doc, err := u.db.Collection("users").Doc(user.ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var userResponse entity.AccountUser
	err = doc.DataTo(&userResponse)
	if err != nil {
		return nil, err
	}

	return &userResponse, nil
}

func (u *UserRepo) Get(ctx context.Context) ([]entity.AccountUser, error) {
	iter := u.db.Documents(ctx, "users")

	defer iter.Stop()

	var documents []entity.AccountUser
	for {
		var user entity.AccountUser
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

func (u *UserRepo) GetById(ctx context.Context, id *string) (*entity.AccountUser, error) {
	doc, err := u.db.Collection("users").Doc(*id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var user entity.AccountUser
	err = doc.DataTo(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) GetByEmail(ctx context.Context, email *string) (*entity.AccountUser, error) {

	doc, err := u.db.Collection("users").Where("email", "==", *email).Documents(ctx).Next()
	if err != nil {
		if err == iterator.Done {
			return nil, errors.New("not found")
		}

		return nil, err
	}

	var user entity.AccountUser
	err = doc.DataTo(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) Update(ctx context.Context, data *entity.AccountUser) (*entity.AccountUser, error) {
	_, err := u.db.Collection("users").Doc(data.ID).Set(ctx, data)
	if err != nil {
		return nil, err
	}

	doc, err := u.db.Collection("users").Doc(data.ID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var user entity.AccountUser
	err = doc.DataTo(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) Delete(ctx context.Context, id *string) error {
	_, err := u.db.Collection("users").Doc(*id).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepo) GetByFilterMany(ctx context.Context, filter []entity.QueryDBClause) ([]entity.AccountUser, error) {
	iter := u.db.Documents(ctx, "users")

	defer iter.Stop()

	var documents []entity.AccountUser
	for {
		var user entity.AccountUser
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

func (u *UserRepo) GetByFilterOne(ctx context.Context, filter []entity.QueryDBClause) (*entity.AccountUser, error) {
	// ctx, span := c.trace.Trace.Start(ctx, "TransactionCategoryRepo.GetByFilterOne")
	// defer span.End()

	hasQuery := false
	hasQueryOr := false
	query := u.db.Collection("users").Query
	queryOr := u.db.Collection("users").Query
	for _, f := range filter {
		if f.Clause == "" || f.Clause == entity.QueryClauseAnd {
			hasQuery = true
			for _, q := range f.Queries {
				condition := checkFirebaseCondition(&q.Condition)
				if q.Key != "" && q.Value != "" && condition != "" {
					query = query.Where(q.Key, condition, q.Value)
				}
			}
		}
		if f.Clause == entity.QueryClauseOr {
			hasQueryOr = true
			for _, q := range f.Queries {
				condition := checkFirebaseCondition(&q.Condition)
				if q.Key != "" && q.Value != "" && condition != "" {
					queryOr = queryOr.Where(q.Key, condition, q.Value)
				}
			}
		}

	}

	var users []entity.AccountUser
	if hasQuery {
		iter := query.Documents(ctx)
		defer iter.Stop()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}

			var user entity.AccountUser
			err = doc.DataTo(&user)
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}
	}

	if hasQueryOr {
		iter := queryOr.Documents(ctx)
		defer iter.Stop()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}

			var user entity.AccountUser
			err = doc.DataTo(&user)
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}
	}

	if len(users) == 0 {
		return nil, errors.New("not found")
	}

	return &users[0], nil
}
