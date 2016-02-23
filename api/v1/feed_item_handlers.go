package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"partisan/auth"
	"partisan/dao"
	"partisan/db"
	m "partisan/models"
	"strconv"
	"time"

	"partisan/Godeps/_workspace/src/github.com/gorilla/websocket"
	"partisan/Godeps/_workspace/src/github.com/jinzhu/gorm"

	"partisan/Godeps/_workspace/src/github.com/gin-gonic/gin"
)

// PostComments stores like data for ease
type PostComments struct {
	RecordID uint64
	Count    int
}

// FeedIndex shows all Feed Items for a particular user
func FeedIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	page := getPage(c)

	friendIDs, err := dao.ConfirmedFriendIDs(user, db)
	if err != nil {
		handleError(err, c)
		return
	}

	friendIDs = append(friendIDs, user.ID)

	feedItems, err := dao.GetFeedByUserIDs(user.ID, friendIDs, page, db)
	if err != nil {
		handleError(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"feed_items": feedItems})
}

func FeedShow(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	uID := c.Param("user_id")
	userID, err := strconv.Atoi(uID)
	if err != nil {
		handleError(&ErrParseID{err}, c)
		return
	}

	page := getPage(c)

	feedItems, err := dao.GetFeedByUserIDs(user.ID, []uint64{uint64(userID)}, page, db)
	if err != nil {
		handleError(err, c)
	}

	c.JSON(http.StatusOK, gin.H{"feed_items": feedItems})
}

func FeedSocket(c *gin.Context) {
	db := db.GetDB(c)
	user, err := auth.CurrentUser(c)
	if err != nil {
		handleError(err, c)
		return
	}

	friendIDs, err := dao.ConfirmedFriendIDs(user, db)
	if err != nil {
		handleError(err, c)
		return
	}

	friendIDs = append(friendIDs, user.ID)

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	msg := make(chan string)
	quit := make(chan bool)

	go feedReadLoop(conn, msg, quit)
	go feedWriteLoop(user.ID, friendIDs, db, conn, msg, quit)
}

func feedReadLoop(c *websocket.Conn, send chan string, quit chan bool) {
	for {
		msgType, r, err := c.NextReader()
		if err != nil {
			c.Close()
			quit <- true
			return
		}

		if msgType == websocket.BinaryMessage {
			fmt.Println("FeedSocket message shouldn't be a binary")
		}

		stamp, err := ioutil.ReadAll(r)
		if err != nil {
			fmt.Println("couldn't read timestamp", r)
		}

		send <- string(stamp)
	}
}

func feedWriteLoop(userID uint64, friendIDs []uint64, db *gorm.DB, c *websocket.Conn, received chan string, quit chan bool) {
	for {
		select {
		case stamp := <-received:
			sec, err := strconv.Atoi(stamp)
			if err != nil {
				fmt.Println("bad timestamp:", string(stamp))
			}

			after := time.Unix(int64(sec), int64(0))

			feedItems, err := dao.GetFeedByUserIDsAfter(userID, friendIDs, after, db)
			if err != nil {
				c.WriteJSON(gin.H{"feed_items": []m.FeedItem{}})
			}

			if err := c.WriteJSON(gin.H{"feed_items": feedItems}); err != nil {
				return
			}
		case <-quit:
			return // kill the loop
		}
	}
}
