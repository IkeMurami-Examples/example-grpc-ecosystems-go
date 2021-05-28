package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string

const (
	appName 	   string = "example"
	defaultCfgPath        = "/etc/" + appName
	defaultCfgName        = appName + ".yaml"
	usageMessage          = "config file (default is " + defaultCfgPath + "/" + defaultCfgName + ")"
)

var rootCmd = &cobra.Command{
	Use: appName,
	Short: "Root command",
	Long: "Root command",
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(gitCommit string) {
	viper.Set("GIT_COMMIT", gitCommit)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", usageMessage)
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Toggle debug mode")
	_ = viper.BindPFlag("DEBUG", rootCmd.PersistentFlags().Lookup("debug"))
	rootCmd.PersistentFlags().BoolP("pretty_log", "p", false, "Pretty log format")
	_ = viper.BindPFlag("PRETTY_LOG", rootCmd.PersistentFlags().Lookup("pretty_log"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name `appName.yaml`.
		viper.SetConfigType("yaml")
		viper.AddConfigPath(defaultCfgPath)
		viper.SetConfigName(appName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
