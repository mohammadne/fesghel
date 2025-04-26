package functional

import (
	"context"
	"testing"
)

func TestURLsShorten(t *testing.T) {
	urlsService.Shorten(context.TODO(), "https://example.com/a-very-long-url")
}
