package v2

import (
	"net/http"
	"partisan/auth"
	"partisan/db"

	"partisan/models.v2/events"
	"partisan/models.v2/posts"
	"partisan/models.v2/users"

	"github.com/gin-gonic/gin"
)

func EventIndex(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	page := getPage(c)

	events, _ := events.SearchForUser(user, page*25, db)
	c.JSON(http.StatusOK, gin.H{"events": events})
}

func EventShow(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	eventID := c.Param("event_id")

	event, err := events.GetByID(eventID, user, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func EventCreate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var binding events.CreatorBinding
	if err := c.Bind(&binding); err != nil {
		if err := c.BindJSON(&binding); err != nil {
			c.AbortWithError(http.StatusNotAcceptable, ErrBinding)
			return
		}
	}

	event, sub, errs := events.New(user, binding)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	tmpFile, _, err := c.Request.FormFile("cover_photo")
	if err == nil {
		defer tmpFile.Close()

		if err := event.AttachCoverPhoto(tmpFile); err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			return
		}
	}

	if err := db.Create(&event).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Create(&sub).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"event": event})
}

func EventUpdate(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	eventID := c.Param("event_id")

	event, err := events.GetByID(eventID, user, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if !event.CanUpdate(user) {
		c.AbortWithError(http.StatusUnauthorized, ErrCannotUpdate)
		return
	}

	var binding events.UpdaterBinding
	if err := c.Bind(&binding); err != nil {
		if err := c.BindJSON(&binding); err != nil {
			c.AbortWithError(http.StatusNotAcceptable, ErrBinding)
			return
		}
	}

	if errs := event.Update(binding); len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Save(&event).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func EventAddHost(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	eventID := c.Param("event_id")

	event, err := events.GetByID(eventID, user, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if !event.CanUpdate(user) {
		c.AbortWithError(http.StatusUnauthorized, ErrCannotUpdate)
		return
	}

	hostID := c.Param("user_id")
	host, err := users.GetByID(hostID, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if event.HasHost(host) {
		c.AbortWithError(http.StatusNotAcceptable, ErrAlreadyExists)
		return
	}

	if sub, err := event.GetSubscription(host, db); err == nil {
		if sub.RSVP != events.RTGoing {
			event.GoingCount++
			event.MaybeCount--
		}

		sub.RSVP = events.RTHost
		if err := db.Save(&sub).Error; err != nil {
			c.AbortWithError(http.StatusNotAcceptable, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"event": event})
		return
	}

	sub, errs := event.NewHost(host)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	if err := db.Create(&sub).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, errs)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"event": event})
}

func EventRemoveHost(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	eventID := c.Param("event_id")
	hostID := c.Param("user_id")

	event, err := events.GetByID(eventID, user, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if !event.CanUpdate(user) {
		c.AbortWithError(http.StatusUnauthorized, ErrCannotUpdate)
		return
	}

	subs, err := event.GetHostSubscriptions(db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if len(subs) < 2 {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	for _, sub := range subs {
		if sub.Subscriber.GetID() == hostID {
			if err := db.Delete(&sub).Error; err != nil {
				c.AbortWithError(http.StatusNotAcceptable, err)
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "removed"})
			return
		}
	}

	c.AbortWithStatus(http.StatusNotFound)
}

func EventGoing(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	eventID := c.Param("event_id")

	event, err := events.GetByID(eventID, user, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if sub, err := event.GetSubscription(user, db); err == nil {
		if sub.RSVP != events.RTHost {
			sub.RSVP = events.RTGoing
			if err := db.Save(&sub).Error; err != nil {
				c.AbortWithError(http.StatusNotAcceptable, err)
				return
			}

			event.GoingCount++
			event.MaybeCount--
		}

		c.JSON(http.StatusOK, gin.H{"event": event})
		return
	}

	sub, errs := event.NewGuest(user, events.RTGoing)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Create(&sub).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func EventMaybe(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	eventID := c.Param("event_id")

	event, err := events.GetByID(eventID, user, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if sub, err := event.GetSubscription(user, db); err == nil {
		if sub.RSVP != events.RTHost {
			sub.RSVP = events.RTMaybe
			if err := db.Save(&sub).Error; err != nil {
				c.AbortWithError(http.StatusNotAcceptable, err)
				return
			}

			event.GoingCount--
			event.MaybeCount++
		}

		c.JSON(http.StatusOK, gin.H{"event": event})
		return
	}

	sub, errs := event.NewGuest(user, events.RTMaybe)
	if len(errs) > 0 {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	if err := db.Create(&sub).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func EventUnsubscribe(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	eventID := c.Param("event_id")

	event, err := events.GetByID(eventID, user, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	sub, err := event.GetSubscription(user, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if sub.RSVP == events.RTGoing {
		event.GoingCount--
	}

	if sub.RSVP == events.RTMaybe {
		event.MaybeCount--
	}

	if err := db.Delete(&sub).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"event": event})
}

func EventPosts(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	eventID := c.Param("event_id")
	page := getPage(c)

	ps, err := posts.ListByParent(user.ID, posts.PTEvent, eventID, page*25, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": ps})
}

func EventDestroy(c *gin.Context) {
	db := db.GetDB(c)

	user, err := auth.CurrentUser(c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	eventID := c.Param("event_id")

	event, err := events.GetByID(eventID, user, db)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if !event.CanDelete(user) {
		c.AbortWithError(http.StatusUnauthorized, ErrCannotDelete)
		return
	}

	if err := db.Delete(&event).Error; err != nil {
		c.AbortWithError(http.StatusNotAcceptable, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})

	db.Where("event_id = ?", eventID).Delete(events.EventSubscription{})
}
