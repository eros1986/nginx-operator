package app

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/eros1986/nginx-operator/pkg/operator"
)

var (
	kubeconfig     string
	watchNamespace string
	resyncSeconds  uint32
)

var serverCmd = &cobra.Command{
	Use:           "server",
	Short:         "Lanch server",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, agrs []string) error {
		config := &operator.OperatorConfig{
			KubeConfigPath: kubeconfig,
			WatchNamespace: watchNamespace,
			ResyncPeriod:   time.Duration(resyncSeconds) * time.Second,
		}
		operator, err := operator.NewOperator(config)
		if err != nil {
			return err
		}

		ctx := context.TODO()
		stopc := make(chan struct{})

		return operator.Run(ctx, stopc)
	},
}

func init() {
	serverCmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "c", "", "pacht to kube config")
	serverCmd.Flags().StringVar(&watchNamespace, "watchNamespace", "",
		"the namespace which operator watches")
	serverCmd.Flags().Uint32Var(&resyncSeconds, "resyncSeconds", 30,
		"resync seconds")
	rootCmd.AddCommand(serverCmd)
}
