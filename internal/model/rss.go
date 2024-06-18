package model

import "time"

type FeedArticle struct {
  ID uint `json:"id" gorm:"primaryKey"`
  FID uint
  Date time.Time
  Title string
  Author string
  Description string
  TorrentURL string
  MagnetURL string
  Link string `json:"link" gorm:"unique"` // 去重
  IsRead bool
}

type Feed struct {
  UID uint `json:"uid" gorm:"primaryKey"`
  Url string
  Title string
  LastBuildDate time.Time
  IsLoading bool
  HasError bool
  Enabled bool
  Priority int
  UseRegex bool
  MustContain string
  MustNotContain string
  EpisodeFilter string
  AffectedFeeds string
  LastMatch time.Time
  IgnoreDays int
  SmartFilter bool
  PreviouslyMatchedEpisodes string
}
