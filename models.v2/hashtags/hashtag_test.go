package hashtags

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nu7hatch/gouuid"

	"github.com/jinzhu/gorm"
)

var testdb *gorm.DB

type testTagger struct {
	Type    string
	ID      string
	Content string
}

func (t testTagger) GetType() string {
	return t.Type
}

func (t testTagger) GetID() string {
	return t.ID
}

func (t testTagger) GetContent() string {
	return t.Content
}

func hashtagSetup() {
	var tags []Hashtag
	testdb.Find(&tags)
	testdb.Delete(&tags)

	var taxs []Taxonomy
	testdb.Find(&taxs)
	testdb.Delete(&taxs)
}

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

	testdb.AutoMigrate(Hashtag{}, Taxonomy{})
}

func TestMain(m *testing.M) {
	var tags []Hashtag
	testdb.Find(&tags)
	testdb.Delete(&tags)

	var taxs []Taxonomy
	testdb.Find(&taxs)
	testdb.Delete(&taxs)

	m.Run()
}

func TestNewCreate(t *testing.T) {
	htagger := testTagger{
		ID: genUUID(),
	}

	var preHCount, preTCount int
	testdb.Model(&Hashtag{}).Count(&preHCount)
	testdb.Model(&Taxonomy{}).Count(&preTCount)

	if err := Create(&htagger, "blah", testdb); err != nil {
		t.Error(err)
	} else {
		var nowHCount, nowTCount int
		testdb.Model(&Hashtag{}).Count(&nowHCount)
		testdb.Model(&Taxonomy{}).Count(&nowTCount)

		if preHCount >= nowHCount {
			t.Error("Counts should be different:", preHCount, nowHCount)
		}

		if preTCount >= nowTCount {
			t.Error("Counts should be different:", preTCount, nowTCount)
		}
	}
}

func TestExistingCreate(t *testing.T) {
	htagger1 := testTagger{ID: genUUID()}
	htagger2 := testTagger{ID: genUUID()}

	Create(&htagger1, "blah", testdb) // create it ahead of time

	var preHCount, preTCount int
	testdb.Model(&Hashtag{}).Count(&preHCount)
	testdb.Model(&Taxonomy{}).Count(&preTCount)

	if err := Create(&htagger2, "blah", testdb); err != nil {
		t.Error(err)
	} else {
		var nowHCount, nowTCount int
		testdb.Model(&Hashtag{}).Count(&nowHCount)
		testdb.Model(&Taxonomy{}).Count(&nowTCount)

		if preHCount != nowHCount {
			t.Error("Counts should be the same:", preHCount, nowHCount)
		}

		if preTCount >= nowTCount {
			t.Error("Counts should be different:", preTCount, nowTCount)
		}
	}
}

func TestExtractTags(t *testing.T) {
	if tags := ExtractTags("#blah"); len(tags) != 1 {
		t.Error("Expected 1 tag, got:", len(tags))
	} else {
		if tags[0] != "blah" {
			t.Error("Expected tag to be \"blah\", got:", tags[0])
		}
	}

	if tags := ExtractTags("#foo #bar"); len(tags) != 2 {
		t.Error("Expected 2 tag, got:", len(tags))
	} else {
		if tags[0] != "foo" || tags[1] != "bar" {
			t.Error("Expected tags to be \"foo\" and \"bar\", got:", tags)
		}
	}

	if tags := ExtractTags("#foo #bar\n# header"); len(tags) != 2 {
		t.Error("Expected 2 tag, got:", len(tags))
	} else {
		if tags[0] != "foo" || tags[1] != "bar" {
			t.Error("Expected tags to be \"foo\" and \"bar\", got:", tags)
		}
	}

}

func genUUID() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return id.String()
}
