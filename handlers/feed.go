package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bsaii/funky-hub-go/database"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type FunkyFeed struct {
	Nickname  string   `json:"nickname" binding:"required"`
	Interests []string `json:"interests" binding:"required"`
	About     string   `json:"about_you" binding:"required"`
}

type FunkyFeeds struct {
	Id        string   `json:"id"`
	About     string   `json:"about_you"`
	Interests []string `json:"interests"`
	Nickname  string   `json:"nickname"`
	CreatedAt string   `json:"created_at"`
}

func CreateFeed(c *gin.Context) {
	var feed FunkyFeed

	if err := c.ShouldBindJSON(&feed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if feed.About == "" || len(feed.Interests) == 0 || feed.Nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fields cannot be empty"})
		return
	}

	stmt, err := database.DB.Prepare("INSERT INTO funkys (about_you, interests, nickname) VALUES($1, $2, $3) RETURNING id, about_you, interests, nickname, created_at")
	if err != nil {
		log.Printf("Failed to add a funky person. Error: %v", err)
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": "Failed to add a funky person", "err": err})
		return
	}
	defer stmt.Close()

	var newFeed FunkyFeeds
	err = stmt.QueryRow(feed.About, pq.Array(feed.Interests), feed.Nickname).Scan(&newFeed.Id, &newFeed.Nickname, pq.Array(&newFeed.Interests), &newFeed.About, &newFeed.CreatedAt)
	if err != nil {
		log.Printf("Failed to add a funky person. Error: %v", err)
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": "Failed to add a funky person", "err": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "Succefully added a new funky person.", "feed": newFeed})
}

func GetFeeds(c *gin.Context) {
	var feeds []FunkyFeeds

	rows, err := database.DB.Query("SELECT * FROM funkys ORDER BY created_at DESC")
	if err != nil {
		log.Printf("Failed to fetch all funky feeds. Error: %v", err)
		c.JSON(http.StatusExpectationFailed, gin.H{"msg": "Failed to fetch all funky feeds.", "err": err})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var feed FunkyFeeds

		if err := rows.Scan(&feed.Id, &feed.Nickname, pq.Array(&feed.Interests), &feed.About, &feed.CreatedAt); err != nil {
			log.Print(err.Error())
			c.JSON(http.StatusExpectationFailed, gin.H{"msg": "Failed to fetch all funky feeds."})
			return
		}
		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		log.Printf("The query encounted an error. Error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to fetch all funky feeds."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feeds": feeds})
}

func GetFeedById(c *gin.Context) {
	feedId := c.Param("id")
	var feed FunkyFeeds

	row := database.DB.QueryRow("SELECT * FROM funkys WHERE id = $1", feedId)
	if err := row.Scan(&feed.Id, &feed.Nickname, pq.Array(&feed.Interests), &feed.About, &feed.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No such feed with id: %v", feedId)
			c.JSON(http.StatusBadRequest, gin.H{"msg": "No such feed exists."})
			return
		}
		log.Printf("Failed to get feed for id: %v. Error: %v", feedId, err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Failed to get feed for id provided", "err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feed": feed})
}
