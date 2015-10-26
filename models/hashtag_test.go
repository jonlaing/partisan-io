package models

import "testing"

func hashtagSetup() {
	var tags []Hashtag
	testDB.Find(&tags)
	testDB.Delete(&tags)

	var taxs []Taxonomy
	testDB.Find(&taxs)
	testDB.Delete(&taxs)
}

func TestNewCreateHashtag(t *testing.T) {
	hashtagSetup()
	defer hashtagSetup() // delete em all again

	htagger := Post{
		ID: 5,
	}

	var preHCount, preTCount int
	testDB.Model(&Hashtag{}).Count(&preHCount)
	testDB.Model(&Taxonomy{}).Count(&preTCount)

	if err := CreateHashtag(&htagger, "blah", &testDB); err != nil {
		t.Error(err)
	} else {
		var nowHCount, nowTCount int
		testDB.Model(&Hashtag{}).Count(&nowHCount)
		testDB.Model(&Taxonomy{}).Count(&nowTCount)

		if preHCount >= nowHCount {
			t.Error("Counts should be different:", preHCount, nowHCount)
		}

		if preTCount >= nowTCount {
			t.Error("Counts should be different:", preTCount, nowTCount)
		}
	}
}

func TestExistingCreateHashtag(t *testing.T) {
	hashtagSetup()
	defer hashtagSetup() // delete em all again

	htagger1 := Post{ID: 5}
	htagger2 := Post{ID: 6}

	CreateHashtag(&htagger1, "blah", &testDB) // create it ahead of time

	var preHCount, preTCount int
	testDB.Model(&Hashtag{}).Count(&preHCount)
	testDB.Model(&Taxonomy{}).Count(&preTCount)

	if err := CreateHashtag(&htagger2, "blah", &testDB); err != nil {
		t.Error(err)
	} else {
		var nowHCount, nowTCount int
		testDB.Model(&Hashtag{}).Count(&nowHCount)
		testDB.Model(&Taxonomy{}).Count(&nowTCount)

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
