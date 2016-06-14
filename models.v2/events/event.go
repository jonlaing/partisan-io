package events

import (
	"mime/multipart"
	"os"
	"time"

	"github.com/jasonmoo/geo"
	"github.com/nu7hatch/gouuid"

	"partisan/imager"
	"partisan/logger"
	"partisan/matcher"

	models "partisan/models.v2"
)

type Event struct {
	ID            string               `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
	Title         string               `json:"title"`
	StartDate     time.Time            `json:"start_date"`
	EndDate       time.Time            `json:"end_date"`
	Latitude      float64              `json:"latitude"`
	Longitude     float64              `json:"longitude"`
	Location      string               `json:"location"`
	PoliticalMap  matcher.PoliticalMap `json:"-" sql:"type:varchar(255)"`
	CenterX       int                  `json:"-"`
	CenterY       int                  `json:"-"`
	Summary       string               `json:"summary"`
	CoverPhotoURL string               `json:"cover_photo_url"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`

	Hosts      []Subscriber `json:"hosts" sql:"-"`
	GoingCount int          `json:"going_count" sql:"-"`
	MaybeCount int          `json:"maybe_count" sql:"-"`
	Match      float64      `json:"match" sql:"-"`
	RSVP       RSVPType     `json:"rsvp" sql:"-"`
}

type Events []Event

type CreatorBinding struct {
	Title     string `json:"title" form:"title" binding:"required"`
	StartDate string `json:"start_date" form:"start_date" binding:"required"`
	EndDate   string `json:"end_date" form:"end_date" binding:"required"`
	Location  string `json:"location" form:"location"`
	Summary   string `json:"summary" form:"summary"`
}

type UpdaterBinding struct {
	Title     string `json:"title" form:"title" binding:"required"`
	StartDate string `json:"start_date" form:"start_date"`
	EndDate   string `json:"end_date" form:"end_date"`
	Location  string `json:"location" form:"location"`
	Summary   string `json:"summary" form:"summary"`
}

func New(host Subscriber, b CreatorBinding) (e Event, s EventSubscription, errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	id, err := uuid.NewV4()
	if err != nil {
		errs["id"] = err
	}

	e.ID = id.String()
	e.Title = b.Title
	e.Summary = b.Summary
	e.CreatedAt = time.Now()
	e.UpdatedAt = e.CreatedAt

	if time, err := parseTime(b.StartDate); err != nil {
		logger.Error.Println(err)
		errs["start_date"] = err
	} else {
		e.StartDate = time
	}

	if time, err := parseTime(b.EndDate); err != nil {
		logger.Error.Println(err)
		errs["end_date"] = err
	} else {
		e.EndDate = time
	}

	var herrs models.ValidationErrors
	s, herrs = e.NewHost(host)
	if len(herrs) > 0 {
		for k, v := range herrs {
			errs[k] = v
		}
	}

	// Blank location is "Everywhere"
	if b.Location != "" {
		if err := e.GetLocation(b.Location); err != nil {
			errs["location"] = err
		}
	}

	x, y := e.PoliticalMap.Center()
	e.CenterX = x
	e.CenterY = y

	if verrs := e.Validate(); len(verrs) > 0 {
		for k, v := range verrs {
			errs[k] = v
		}
	}

	return
}

func (e *Event) Update(b UpdaterBinding) (errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	if b.Title != "" {
		e.Title = b.Title
	}

	if len(b.StartDate) > 0 {
		if time, err := parseTime(b.StartDate); err != nil {
			errs["start_date"] = err
		} else {
			e.StartDate = time
		}
	}

	if len(b.EndDate) > 0 {
		if time, err := parseTime(b.EndDate); err != nil {
			errs["end_date"] = err
		} else {
			e.EndDate = time
		}
	}

	if b.Summary != "" {
		e.Summary = b.Summary
	}

	if b.Location != "" {
		if err := e.GetLocation(b.Location); err != nil {
			errs["location"] = err
		}
	}

	if verrs := e.Validate(); len(verrs) > 0 {
		for k, v := range verrs {
			errs[k] = v
		}
	}

	return
}

func (e Event) Validate() (errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	if e.StartDate.Before(time.Now()) {
		errs["start_date"] = ErrStartDate
	}

	if e.EndDate.Before(time.Now()) {
		errs["end_date"] = ErrEndDate
	}

	if e.EndDate.Before(e.StartDate) {
		errs["end_date"] = ErrEndDate
	}

	return
}

func (e Event) HasHost(host Subscriber) bool {
	for _, h := range e.Hosts {
		if h != nil && h.GetID() == host.GetID() {
			return true
		}
	}

	return false
}

func (e Event) CanUpdate(host Subscriber) bool {
	return e.HasHost(host)
}

func (e Event) CanDelete(host Subscriber) bool {
	return e.HasHost(host)
}

// GetLocation finds the latitude/longitude by postal code
func (e *Event) GetLocation(address string) error {
	location, err := geo.Geocode(address)
	if err != nil {
		return err
	}

	e.Location = location.Address
	e.Latitude = location.Lat
	e.Longitude = location.Lng

	return nil
}

func (e *Event) AttachCoverPhoto(f multipart.File) error {
	var err error
	var fullPath string
	isS3 := false

	if len(os.Getenv("AWS_ACCESS_KEY_ID")) > 0 {
		isS3 = true
	}

	processor := imager.ImageProcessor{File: f}

	if err := processor.Resize(1500); err != nil {
		return err
	}

	if isS3 {
		fullPath, err = processor.Save("/img")
	} else {
		fullPath, err = processor.Save("/localfiles/img")
	}
	if err != nil {
		return err
	}

	e.CoverPhotoURL = fullPath

	return nil
}

func parseTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", s)
}
