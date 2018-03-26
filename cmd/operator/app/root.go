package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const defaultDebugLevel uint32 = 4

var debugLevel uint32

var rootCmd = &cobra.Command{
	Use:           "demo-operator",
	Short:         "demo-operator is an operator demo.",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Please use -h to see usage")
	},
}

func init() {
	rootCmd.PersistentFlags().Uint32VarP(&debugLevel, "debuglevel", "l",
		defaultDebugLevel,
		"log debug level: 0[panic] 1[fatal] 2[error] 3[warn] 4[info] 5[debug]")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(log.Level(debugLevel))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
