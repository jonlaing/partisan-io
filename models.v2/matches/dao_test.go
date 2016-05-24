package matches

import (
	"fmt"
	"os"
	"testing"
	"time"

	"partisan/matcher"

	"partisan/models.v2/users"

	_ "github.com/lib/pq"
	"github.com/nu7hatch/gouuid"

	"github.com/jinzhu/gorm"
)

var testdb *gorm.DB
var me, tooYoung, tooOld, outOfArea, opponent, wrongGender, find users.User

func init() {
	var err error
	connString := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_TEST_USER"), os.Getenv("DB_TEST_NAME"), os.Getenv("DB_TEST_PW"))
	testdb, err = gorm.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	if err := testdb.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		panic(err)
	}

	testdb.AutoMigrate(users.User{})
}

func TestMain(m *testing.M) {
	key, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	apiKey := key.String()

	exp := time.Now().Add(168 * time.Hour)

	politicalMap := matcher.PoliticalMap{
		0, 0, 0, 0,
		0, 0, 0, 0,
		1, 1, 0, 0,
		1, 1, 0, 0,
	}
	px, py := politicalMap.Center()

	me = users.User{
		Username:     "me",
		Email:        "me@email.com",
		Gender:       "dragon",
		Birthdate:    time.Now().AddDate(-25, 0, 0),
		PoliticalMap: politicalMap,
		CenterX:      px,
		CenterY:      py,
		LookingFor:   3,
		APIKey:       apiKey,
		APIKeyExp:    exp,
	}
	if err := testdb.Save(&me).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&me)

	tooYoung = users.User{
		Username:     "tooyoung",
		Email:        "tooyoung@email.com",
		Gender:       "unicorn",
		Birthdate:    time.Now(),
		PoliticalMap: politicalMap,
		CenterX:      px,
		CenterY:      py,
		LookingFor:   3,
		APIKey:       apiKey,
		APIKeyExp:    exp,
	}
	if err := testdb.Save(&tooYoung).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&tooYoung)

	tooOld = users.User{
		Username:     "tooold",
		Email:        "tooold@email.com",
		Gender:       "unicorn",
		Birthdate:    time.Now().AddDate(-200, 0, 0),
		PoliticalMap: politicalMap,
		CenterX:      px,
		CenterY:      py,
		LookingFor:   3,
		APIKey:       apiKey,
		APIKeyExp:    exp,
	}
	if err := testdb.Save(&tooOld).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&tooOld)

	outOfArea = users.User{
		Username:     "outOfArea",
		Email:        "outOfArea@email.com",
		Gender:       "unicorn",
		Birthdate:    time.Now().AddDate(-25, 0, 0),
		PoliticalMap: politicalMap,
		CenterX:      px,
		CenterY:      py,
		Latitude:     1.0,
		Longitude:    1.0,
		LookingFor:   3,
		APIKey:       apiKey,
		APIKeyExp:    exp,
	}
	if err := testdb.Save(&outOfArea).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&outOfArea)

	opponentMap := matcher.PoliticalMap{
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 0, 0, 0,
		0, 0, 0, 0,
	}
	opx, opy := opponentMap.Center()
	opponent = users.User{
		Username:     "opponent",
		Email:        "opponent@email.com",
		Gender:       "unicorn",
		Birthdate:    time.Now().AddDate(-25, 0, 0),
		PoliticalMap: opponentMap,
		CenterX:      opx,
		CenterY:      opy,
		LookingFor:   3,
		APIKey:       apiKey,
		APIKeyExp:    exp,
	}
	if err := testdb.Save(&opponent).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&opponent)

	wrongGender = users.User{
		Username:     "wronggender",
		Email:        "wronggender@email.com",
		Gender:       "dragon",
		Birthdate:    time.Now().AddDate(-25, 0, 0),
		PoliticalMap: politicalMap,
		CenterX:      px,
		CenterY:      py,
		LookingFor:   3,
		APIKey:       apiKey,
		APIKeyExp:    exp,
	}
	if err := testdb.Save(&wrongGender).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&wrongGender)

	find = users.User{
		Username:     "find",
		Email:        "find@email.com",
		Gender:       "unicorn",
		Birthdate:    time.Now().AddDate(-25, 0, 0),
		PoliticalMap: politicalMap,
		CenterX:      px,
		CenterY:      py,
		LookingFor:   3,
		APIKey:       apiKey,
		APIKeyExp:    exp,
	}
	if err := testdb.Save(&find).Error; err != nil {
		panic(err)
	}
	defer testdb.Delete(&find)

	m.Run()
}

func TestList(t *testing.T) {
	search := SearchBinding{
		Gender:     "unicorn",
		MinAge:     21,
		MaxAge:     30,
		Radius:     25,
		LookingFor: 7,
		Page:       0,
	}

	matches, err := List(me, search, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	if len(matches) != 1 {
		t.Error("Expected only one result, found:", len(matches))
		return
	}

	if matches[0].User.ID != find.ID {
		t.Error("Expected to find user with ID:", find.ID, "got:", matches[0].User.ID)
	}

	if matches[0].Match == 0 {
		t.Error("Expected non-zero match")
	}

	search.Gender = ""

	matches, err = List(me, search, testdb)
	if err != nil {
		t.Error("Unexpected error:", err)
	}

	if len(matches) != 2 {
		t.Error("Expected 2 results, found:", len(matches))
		return
	}
}
