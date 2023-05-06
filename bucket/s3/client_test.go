package s3

import (
	"context"
	"strings"
	"testing"
)

func TestClient_Upload(t *testing.T) {
	c := New(context.TODO(), cfg)

	reader := strings.NewReader("HELLO WORLD")

	err := c.Upload(context.TODO(), "here.txt", reader)
	if err != nil {
		t.Fatal(err)
	}
}
