package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-logrusutil/logrusutil/logctx"
	"github.com/govinda-attal/winning11/feeds"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	cfgFile string
	cfg     Config

	articleFile string

	log = logctx.Default
)

var rootCmd = &cobra.Command{
	Use: "winning11",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/app-cfg.yaml", "application config file")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	validateCmd.PersistentFlags().StringVar(&articleFile, "article", "", "path to an article to validate")

	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(versionCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigFile(cfgFile)
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("unable to unmarshal configuration locally")
	}
	lvl, _ := logrus.ParseLevel(cfg.Log.Level)
	log.Logger.SetLevel(lvl)
	if cfg.Log.Format == "JSON" {
		log.Logger.SetFormatter(&logrus.JSONFormatter{})
	}
}

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		db := cfg.Validator.DB
		err := setupMigrations(ctx, db.URI, db.Name)
		if err != nil {
			panic(err)
		}
		fmt.Println("migrations applied")
	},
}

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:", version)
	},
}
var validateCmd = &cobra.Command{
	Use: "validate",
	Run: func(cmd *cobra.Command, args []string) {
		articleFile = path.Clean(articleFile)
		bb, err := os.ReadFile(articleFile)
		if err != nil {
			panic(err)
		}

		var article feeds.Article
		if err := json.Unmarshal(bb, &article); err != nil {
			panic(err)
		}

		ctx := context.Background()
		rp, err := ruleProvider(ctx)
		if err != nil {
			panic(err)
		}

		defer func() {
			if err := rp.Shutdown(ctx); err != nil {
				panic(err)
			}
		}()

		if err := feeds.NewValidator(rp).Validate(ctx, article); err != nil {
			fmt.Printf("given article with topic (%s) is invalid\n", article.Topic)
			fmt.Println("error(s):", err)
			os.Exit(1)
		}
		fmt.Printf("given article with topic (%s) is valid\n", article.Topic)
	},
}
