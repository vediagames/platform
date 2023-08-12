package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"

	"github.com/vediagames/platform/config"
	"github.com/vediagames/platform/quote"
)

type Quote struct {
	Quote    string `json:"quote"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

func QuotesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "quotes",
		Short: "Insert quotes",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg := cmd.Context().Value(config.ContextKey).(config.Config)

			db, err := sqlx.Open("postgres", cfg.PostgreSQL.VediaGamesConnectionString)
			if err != nil {
				return fmt.Errorf("failed to open db connection: %w", err)
			}

			quoteService := quote.New(db)

			category := "mom"
			apiURL := fmt.Sprintf("https://api.api-ninjas.com/v1/quotes?category=%s&limit=10", category)
			req, err := http.NewRequest("GET", apiURL, nil)
			if err != nil {
				return err
			}

			req.Header.Set("X-Api-Key", "dqOXrQOm9eP0QBnpLO0g6w==3ISkz5v1BrnkgS4y")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				var quotes []Quote
				err = json.Unmarshal(body, &quotes)
				if err != nil {
					return err
				}

				for _, qu := range quotes {
					fmt.Printf("%q,%q\n", qu.Author, qu.Quote)
					err := quoteService.Insert([]quote.Quote{
						{
							Author: qu.Author,
							Quote:  qu.Quote,
						},
					})
					if err != nil {
						fmt.Println("failed to insert: ", err)
					}
				}
			} else {
				fmt.Println("Error:", resp.StatusCode)
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				fmt.Println(body)
			}

			// 			file, err := os.Open(cfg.QuotesCSV)
			// 			if err != nil {
			// 				return fmt.Errorf("failed to open csv file: %w", err)
			// 			}
			// 			defer file.Close()

			// 			lines, err := csv.NewReader(file).ReadAll()
			// 			if err != nil {
			// 				return fmt.Errorf("failed to read csv file: %w", err)
			// 			}

			// 			for i, line := range lines {
			// 				q := quote.Quote{
			// 					Author: line[0],
			// 					Quote:  line[1],
			// 				}

			// 				if err := quoteService.Insert([]quote.Quote{q}); err != nil {
			// 					return fmt.Errorf("failed to insert quote at index %d: %w", i, err)
			// 				}
			// 			}

			return nil
		},
	}
}
