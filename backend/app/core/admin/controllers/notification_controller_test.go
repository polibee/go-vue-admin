package controllers

import "testing"

func TestNormalizeNotificationQuery(t *testing.T) {
	page, perPage := normalizeNotificationQuery(0, 1000)
	if page != 1 || perPage != 100 {
		t.Fatalf("unexpected normalized query: page=%d perPage=%d", page, perPage)
	}
	page, perPage = normalizeNotificationQuery(3, 25)
	if page != 3 || perPage != 25 {
		t.Fatalf("valid query should be preserved: page=%d perPage=%d", page, perPage)
	}
}
