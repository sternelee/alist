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
	// "github.com/anacrolix/torrent/bencode"
	// "github.com/anacrolix/torrent/metainfo"
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

	var articles []model.FeedArticle
	for _, item := range rss.Channel.Items {
		log.Debugf("%+v", item)
		article := model.FeedArticle{
			Description: item.Description,
			Date:        getTime(item.PubDate),
			Link:        item.Link,
			Title:       item.Title,
			TorrentURL:  item.Link + ".torrent",
		}
		articles = append(articles, article)
	}

	common.SuccessResp(c, common.PageResp{
		Content: articles,
	})
}

func getTime(str string) time.Time {
	stamp, _ := time.Parse("2006-01-02 15:04:05", str)
	return stamp
}

// func torrentHash(url string) string {
// 	// 下载种子文件
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
//
// 	// 读取种子文件内容
// 	data, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// 解析种子文件
// 	metaInfo, err := metainfo.Load(data)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// 获取 info 字典
// 	info := metaInfo.Info
//
// 	// 创建磁力链接
// 	magnetLink, err := bencode.EncodeString(map[string]interface{}{
// 		"xt": "urn:btih:" + info.HashInfoBytes().HexString(),
// 		"dn": info.Name,
// 		"tr": metaInfo.AnnounceList,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// 打印磁力链接
// 	fmt.Println("磁力链接:", magnetLink)
// }
