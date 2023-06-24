package datafuser

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"

	"github.com/vediagames/zeroerror"

	gamedomain "github.com/vediagames/platform/game/domain"
	tagdomain "github.com/vediagames/platform/tag/domain"
)

type Config struct {
	DB          *sqlx.DB
	TagService  tagdomain.Service
	GameService gamedomain.Service
}

func (c Config) Validate() error {
	var err zeroerror.Error

	err.AddIf(c.DB == nil, fmt.Errorf("empty DB"))
	err.AddIf(c.TagService == nil, fmt.Errorf("empty tag service"))
	err.AddIf(c.GameService == nil, fmt.Errorf("empty game service"))

	return err.Err()
}

func New(c Config) Service {
	if err := c.Validate(); err != nil {
		panic(fmt.Errorf("invalid config: %w", err))
	}

	return Service{
		db:          c.DB,
		tagService:  c.TagService,
		gameService: c.GameService,
	}
}

type fuser func(context.Context) error

type file string

const (
	FileTrending       = file("trending")
	FilePromotedTags   = file("promoted_tags")
	FileTopTags        = file("top_tags")
	FileWhatOthersPlay = file("what_others_play")
	FileQuotes         = file("quotes")
)

type FilePath string

type Service struct {
	db          *sqlx.DB
	tagService  tagdomain.Service
	gameService gamedomain.Service
	fusers      map[string]fuser
	files       map[file]FilePath
}

func (s Service) Fuse(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)

	for name, f := range s.fusers {
		f := f
		name := name
		errGroup.Go(func() error {
			if err := f(ctx); err != nil {
				return fmt.Errorf("failed to fuse %q: %w", name, err)
			}

			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("failed to wait: %w", err)
	}

	return nil
}

type Game struct {
	Slug      string
	Placement string
}

func (s Service) fuseTrending(ctx context.Context) error {

	// 	var records []Record
	// 	for _, line := range lines {
	// 		data := Record{
	// 			Name:  line[0],
	// 			Slug:  line[1],
	// 			Likes: line[2],
	// 		}
	// 		records = append(records, data)
	// 	}

	return nil
}

func (s Service) fuseWhatOthersPlay(ctx context.Context) error {
	return nil
}

type PromotedTag struct {
	Slug      string
	Placement string
	Icon      string
}

func (s Service) fusePromotedTags(ctx context.Context) error {
	return nil
}

type TopTag struct {
	Slug      string
	Placement string
}

func (s Service) fuseTopTags(ctx context.Context) error {
	return nil
}

type Quote struct {
	ID     string
	Author string
	Quote  string
}

func (s Service) fuseQuotes(ctx context.Context) error {
	return nil
}

func (s Service) readCsv(f file) ([][]string, error) {
	filePath, ok := s.files[FileTrending]
	if !ok {
		return nil, fmt.Errorf("file not found in map %q", f)
	}

	file, err := os.Open(string(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filePath, err)
	}

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read all: %w", err)
	}

	return lines, nil
}
