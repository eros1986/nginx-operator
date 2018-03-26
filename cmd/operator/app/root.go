package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const defaultDebugLevel uint32 = 4

var debugLevel uint32

var rootCmd = &cobra.Command{
	Use:           "nginx-operator",
	Short:         "nginx-operator is an nginx operator demo.",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Please use -h to see uaage")
	},
}

func init() {
	rootCmd.PersistentFlags().Uint32Varp(&debugLevel, "debuglevel", "l",
		debugLevel,
		"log debug level:0[panic] 1[fatal] 2[error] 3[warn] 4[info] 5[debug]")

	cobra.OnInitialize(initConfig)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
