package service

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/vediagames/vediagames.com/bucket/bunny"
	"github.com/vediagames/vediagames.com/image/domain"
	"github.com/vediagames/vediagames.com/image/processor/imagor"
)

func TestService_Get(t *testing.T) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	bunnyBucket := bunny.New(bunny.Config{
		URL:       "https://storage.bunnycdn.com",
		AccessKey: "",
		Zone:      "vedia-games",
		Client:    httpClient,
	})

	processer := imagor.New(imagor.Config{
		Client:       httpClient,
		URL:          "http://localhost:8000",
		Secret:       "vediagames",
		BucketClient: bunnyBucket,
	})

	svc := New(Config{
		URL:       "https://images.vediagames.com/file/vg-images",
		Processor: processer,
		Client:    httpClient,
	})

	get, err := svc.Get(context.TODO(), domain.GetRequest{
		Slug: "kirka-io",
		Image: domain.Image{
			Format: domain.FormatWebp,
			Width:  384,
			Height: 215,
		},
		Original: domain.OriginalThumbnail512x384,
		Resource: domain.ResourceGame,
	})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(get.URL)

}

func Test_imagePath(t *testing.T) {
	want := "games/kirka-io/thumb512x384.webp"
	got := imagePath(domain.ResourceGame, "kirka-io", domain.Image{
		Format: domain.FormatWebp,
		Width:  512,
		Height: 384,
	})

	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}

	t.Logf("want %q, got %q", want, got)
}

func Test_imageURL(t *testing.T) {
	want := "https://images.vediagames.com/file/vg-images/games/kirka-io/thumb512x384.webp"
	imgPath := imagePath(domain.ResourceGame, "kirka-io", domain.Image{
		Format: domain.FormatWebp,
		Width:  512,
		Height: 384,
	})

	got := imageURL("https://images.vediagames.com/file/vg-images", imgPath)

	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}

	t.Logf("want %q, got %q", want, got)
}
