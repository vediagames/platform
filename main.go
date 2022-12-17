package main

import "github.com/vediagames/vediagames.com/auth"

func main() {
	//var (
	//	configFilePathFlag string
	//)
	//
	//rootCmd := &cobra.Command{}
	//
	//rootCmd.PersistentFlags().
	//	StringVarP(&configFilePathFlag, "config", "c", "config.yml", "Path to the config file")
	//
	//rootCmd.AddCommand(cmd.ServerCmd())
	//rootCmd.AddCommand(cmd.MigrateCmd())
	//rootCmd.AddCommand(cmd.StubCmd())
	//rootCmd.AddCommand(cmd.RefreshCmd())
	//
	//logger := environment.InitLogger()
	//
	//ctx := logger.WithContext(context.Background())
	//
	//cfg, err := config.New(configFilePathFlag)
	//if err != nil {
	//	logger.Fatal().Err(fmt.Errorf("failed to load config: %w", err))
	//}
	//
	//ctx = context.WithValue(ctx, config.ContextKey, cfg)
	//rootCmd.SetContext(ctx)
	//
	//if err := rootCmd.Execute(); err != nil {
	//	logger.Fatal().Err(fmt.Errorf("failed to execute: %w", err))
	//}

	auth.Auth()
}
