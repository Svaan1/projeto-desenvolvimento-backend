package quiz

import (
	"backendProject/internal/db"
	"context"
	"log"
)

type Repository struct {
	DB db.Database
}

func NewRepository(db db.Database) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r *Repository) GetQuiz(ctx context.Context, qry string) (Quiz, error) {
	quiz := Quiz{}
	log.Printf("Getting quiz from key: %s", qry)
	err := r.DB.GetObject(ctx, qry, &quiz)
	return quiz, err

}

func (r *Repository) SetQuiz(ctx context.Context, key string, quiz Quiz) error {
	log.Printf("Setting quiz with key: %s", key)
	return r.DB.SetObject(ctx, key, quiz)
}
