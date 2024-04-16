package handles

import (
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/alist-org/alist/v3/internal/db"
	"github.com/alist-org/alist/v3/internal/model"
	"github.com/alist-org/alist/v3/server/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func GetFeeds(c *gin.Context) {
	var req model.PageReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	req.Validate()
	log.Debugf("%+v", req)
	feeds, total, err := db.GetFeeds(req.Page, req.PerPage)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, common.PageResp{
		Content: feeds,
		Total:   total,
	})
}

func CreateFeed(c *gin.Context) {
	var req model.Feed
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := db.CreateFeed(&req); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
  // TODO: 同时搜索并添加 FeedArticle 列表
	common.SuccessResp(c)
}

func UpdateFeed(c *gin.Context) {
	var req model.Feed
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := db.UpdateFeed(&req); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c)
}

func DeleteFeed(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := db.DeleteFeedById(id); err != nil {
		common.ErrorResp(c, err, 500, true)
		return
	}
	if err := db.DeleteFeedArticleByFid(id); err != nil {
		common.ErrorResp(c, err, 500, true)
		return
	}
	common.SuccessResp(c)
}

func DeleteFeedArticle(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := db.DeleteFeedArticleById(id); err != nil {
		common.ErrorResp(c, err, 500, true)
		return
	}
	common.SuccessResp(c)
}

func GetFeedArticles(c *gin.Context) {
	var req model.PageReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	req.Validate()
	log.Debugf("%+v", req)
	articles, total, err := db.GetFeedArticles(req.Page, req.PerPage)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, common.PageResp{
		Content: articles,
		Total:   total,
	})
}

// 定义 RSS 结构体
type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func SearchFeed(c *gin.Context) {
	link := c.Query("link")
	// 发送 HTTP 请求
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 解析 XML 数据
	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		panic(err)
	}

	// https://acg.rip/.xml?term=%E5%85%B3%E4%BA%8E%E6%88%91%E8%BD%AC%E7%94%9F%E5%8F%98%E6%88%90%E5%8F%B2%E8%8E%B1%E5%A7%86%E8%BF%99%E6%A1%A3%E4%BA%8B+%E7%AC%AC%E4%B8%89%E5%AD%A3

  var articles []model.FeedArticle
	for _, item := range rss.Channel.Items {
    article := model.FeedArticle {
      Description: item.Description,
      Date: getTime(item.PubDate),
      Link: item.Link,
      Title: item.Title,
      // TorrentURL: item.Link + ".torrent",
    }
    articles = append(articles, article)
	}

	common.SuccessResp(c, common.PageResp{
		Content: articles,
	})
}

func getTime(str string) (time.Time) {
	stamp, _ := time.Parse("2006-01-02 15:04:05", str)
  return stamp
}
