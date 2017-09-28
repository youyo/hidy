package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//var cfgFile string
var (
	name    string
	file    string
	str     string
	value   string
	profile string
)

var RootCmd = &cobra.Command{
	Use:   "hidy",
	Short: "Manage credentials using aws parameter store.",
	Long:  `Manage credentials using aws parameter store.`,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	/*
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			viper.AddConfigPath(home)
			viper.SetConfigName(".hidy")
		}
	*/

	viper.SetEnvPrefix("hidy")
	viper.AutomaticEnv()

	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	//}
}
