package models

import "testing"

func TestCanDelete(t *testing.T) {
	basic := User{ID: 1}
	advertiser := User{ID: 2, Type: UTAdvertiser}
	admin := User{ID: 3, Type: UTAdmin}

	post1 := Post{UserID: 1}
	post2 := Post{UserID: 2}

	if !basic.CanDelete(&post1) {
		t.Error("Basic user should be able to delete their own post")
	}

	if basic.CanDelete(&post2) {
		t.Error("Basic user should not be able to delete this post")
	}

	if !advertiser.CanDelete(&post2) {
		t.Error("Advertiser user should be able to delete their own post")
	}

	if advertiser.CanDelete(&post1) {
		t.Error("Advertsier user should not be able to delete this post")
	}

	if !admin.CanDelete(&post1) || !admin.CanDelete(&post2) {
		t.Error("Admin user should be able to delete anything")
	}
}
