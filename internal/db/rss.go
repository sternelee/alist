package db

import (
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/pkg/errors"
)

func CreateFeed(feed *model.Feed) error {
	// TODO: 需要爬取article 才保存
	return errors.WithStack(db.Create(feed).Error)
}

func UpdateFeed(feed *model.Feed) error {
	return errors.WithStack(db.Save(feed).Error)
}

func DeleteFeedById(id int) error {
	return errors.WithStack(db.Delete(&model.Storage{}, id).Error)
}

func GetFeeds(pageIndex, pageSize int) ([]model.Feed, int64, error) {
	feedDB := db.Model(&model.Feed{})
	var count int64
	if err := feedDB.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get feed count")
	}
	var feeds []model.Feed
	if err := feedDB.Order(columnName("id")).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&feeds).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return feeds, count, nil
}

func CreateFeedArticles(articles *[]model.FeedArticle) error {
  return errors.WithStack(db.Create(articles).Error)
}

func GetFeedArticles(pageIndex, pageSize int) ([]model.FeedArticle, int64, error) {
	articleDB := db.Model(&model.FeedArticle{})
	var count int64
	if err := articleDB.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get feed count")
	}
	var articles []model.FeedArticle
	if err := articleDB.Order(columnName("id")).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return articles, count, nil
}

func GetFeedArticlesByFid(fid, pageIndex, pageSize int) ([]model.FeedArticle, int64, error) {
	articleDB := db.Model(&model.FeedArticle{})
	var count int64
	if err := articleDB.Where("fid = ?", fid).Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get feed count")
	}
	var articles []model.FeedArticle
	if err := articleDB.Where("fid = ?", fid).Order(columnName("id")).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return articles, count, nil
}

func UpdateFeedArticle(article *model.FeedArticle) error {
  return errors.WithStack(db.Save(article).Error)
}

func DeleteFeedArticleById(id int) error {
  return errors.WithStack(db.Delete(&model.FeedArticle{}, id).Error)
}

func DeleteFeedArticleByFid(fid int) error {
  return errors.WithStack(db.Where("fid = ?", fid).Delete(&model.FeedArticle{}).Error)
}
