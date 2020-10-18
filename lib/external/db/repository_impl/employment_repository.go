package repository_impl

import (
	"tamastudy_news_crawler/lib/entity/repository_inter"
	"tamastudy_news_crawler/lib/external/db"
)

type employmentRepository struct {
	dbHandler *db.Handler
}

func NewEmploymentRepository(dbHandler *db.Handler) repository_inter.IEmploymentRepository {
	employmentRepository := employmentRepository{dbHandler}
	return employmentRepository
}