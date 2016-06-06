package events

import (
	"time"

	"github.com/jasonmoo/geo"
	"github.com/nu7hatch/gouuid"

	"partisan/matcher"

	models "partisan/models.v2"
)

type Event struct {
	ID            string               `json:"id" gorm:"primary_key" sql:"type:uuid;default:uuid_generate_v4()"`
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
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Location  string    `json:"location"`
	Summary   string    `json:"summary"`
}

type UpdaterBinding struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Location  string    `json:"location"`
	Summary   string    `json:"summary"`
}

func New(host Subscriber, b CreatorBinding) (e Event, s EventSubscription, errs models.ValidationErrors) {
	errs = make(models.ValidationErrors)

	id, err := uuid.NewV4()
	if err != nil {
		errs["id"] = err
	}

	e.ID = id.String()
	e.StartDate = b.StartDate
	e.EndDate = b.EndDate
	e.Summary = b.Summary
	e.CreatedAt = time.Now()
	e.UpdatedAt = e.CreatedAt

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

	if !b.StartDate.IsZero() {
		e.StartDate = b.StartDate
	}

	if !b.EndDate.IsZero() {
		e.EndDate = b.EndDate
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
