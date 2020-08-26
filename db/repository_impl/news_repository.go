package repository_impl

import (
	"tamastudy_news_crawler/db"
	"tamastudy_news_crawler/domain/model"
	"tamastudy_news_crawler/domain/repository_interface"
)

type newsRepository struct {
	dbHandler *db.Handler
}

func NewNewsRepository(dbHandler *db.Handler) repository_interface.INewsRepository {
	newsRepository := newsRepository{dbHandler}
	return newsRepository
}

func (n newsRepository) DeleteAllByPortal(portal string) error{
	if err := n.dbHandler.Conn.Where("portal = ?", portal).Delete(model.News{}).Error; err != nil{
		return err
	}
	return nil
}

func (n newsRepository) Create(news *model.News)  error{
	if err := n.dbHandler.Conn.Create(news).Error; err != nil{
		return err
	}
	return nil
}

func (n newsRepository) DeleteAllByPortalAndAllCreate(portal string, news []*model.News)  error{
	tx := n.dbHandler.Conn.Begin()

	if err := n.DeleteAllByPortal(portal); err != nil{
		tx.Rollback()
		return err
	}

	for _, record := range news{
		if err := n.Create(record); err != nil{
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}