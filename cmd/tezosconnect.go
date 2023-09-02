package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/OneOf-Inc/firefly-tezosconnect/internal/tezos"
	"github.com/hyperledger/firefly-common/pkg/config"
	"github.com/hyperledger/firefly-common/pkg/i18n"
	"github.com/hyperledger/firefly-common/pkg/log"
	fftmcmd "github.com/hyperledger/firefly-transaction-manager/cmd"
	"github.com/hyperledger/firefly-transaction-manager/pkg/fftm"
	txhandlerfactory "github.com/hyperledger/firefly-transaction-manager/pkg/txhandler/registry"
	"github.com/hyperledger/firefly-transaction-manager/pkg/txhandler/simple"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var sigs = make(chan os.Signal, 1)

var rootCmd = &cobra.Command{
	Use:   "tezosconnect",
	Short: "Hyperledger FireFly Connector for Tezos blockchain",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

var cfgFile string

var connectorConfig config.Section

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file")
	rootCmd.AddCommand(fftmcmd.ClientCommand())
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	fftm.InitConfig()
	connectorConfig = config.RootSection("connector")
	tezos.InitConfig(connectorConfig)
	txhandlerfactory.RegisterHandler(&simple.TransactionHandlerFactory{})
}

func run() error {
	initConfig()
	err := config.ReadConfig("tezosconnect", cfgFile)

	// Setup logging after reading config (even if failed), to output header correctly
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	ctx = log.WithLogger(ctx, logrus.WithField("pid", fmt.Sprintf("%d", os.Getpid())))
	ctx = log.WithLogger(ctx, logrus.WithField("prefix", "evmconnect"))

	config.SetupLogging(ctx)

	// Deferred error return from reading config
	if err != nil {
		cancelCtx()
		return i18n.WrapError(ctx, err, i18n.MsgConfigFailed)
	}

	// Init connector
	c, err := tezos.NewTezosConnector(ctx, connectorConfig)
	if err != nil {
		return err
	}
	m, err := fftm.NewManager(ctx, c)
	if err != nil {
		return err
	}

	// Setup signal handling to cancel the context, which shuts down the API Server
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.L(ctx).Infof("Shutting down due to %s", sig.String())
		cancelCtx()
	}()

	return runManager(ctx, m)
}

func runManager(ctx context.Context, m fftm.Manager) error {
	err := m.Start()
	if err != nil {
		return err
	}
	<-ctx.Done()
	m.Close()
	return nil
}
