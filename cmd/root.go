package cmd

import (
	"github.com/msoft-dev/filetify/pkg/shared"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "filetify",
	Short: "Filetify is a file syncrhonization client-server application",
	Long: `With Filetify you can quickly synchronize files within
you local network. Just put Filetify on your server and execute it as server,
then run your Filetify clients on every computer you want. Filetify will automatically
synchronize your files and store them on server. Fast, no database, no additional resources, just storage.'`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//Configuration
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.filetify.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".filetify" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".filetify")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("An error occured during loading '%v' configuration file: %v \n", viper.ConfigFileUsed(), err)
	}

	log.Printf("Configuration file loaded: %+v\n", *shared.GetConfiguration())
}
