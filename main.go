package main

import (
	"github.com/seadiaz/kafka-to-websocket/kafka"
	server "github.com/seadiaz/kafka-to-websocket/websocket"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "kafka-to-websocket",
		Short: "consume messages from a topic and publish to a websocket",
		Long:  "consume messages from a topic and publish to a websocket",
		Run:   executeRootCommand,
	}
)

func init() {
	rootCmd.Flags().StringP("brokers", "b", "localhost:9092", "list of brokers separated by comma")
	rootCmd.Flags().StringP("group-id", "g", "k2w", "name of the kafka group id")
	rootCmd.Flags().StringP("topic", "t", "", "topic from to consume messages")
	rootCmd.Flags().StringP("addr", "a", "localhost:8000", "address to listen websocket server")
	rootCmd.Flags().StringP("base-path", "p", "", "base path for websocket server")

	rootCmd.MarkFlagRequired("topic")

	viper.SetEnvPrefix("k2w")
	viper.AutomaticEnv()
}

func main() {
	rootCmd.Execute()
}

func executeRootCommand(cmd *cobra.Command, _ []string) {
	paramsKafka := &kafka.Params{
		Brokers: cmd.Flag("brokers").Value.String(),
		GroupID: cmd.Flag("group-id").Value.String(),
		Topic:   cmd.Flag("topic").Value.String(),
	}
	c := make(chan []byte)
	go kafka.Run(paramsKafka, c)

	paramsWebsocket := &server.Params{
		Addr:     cmd.Flag("addr").Value.String(),
		BasePath: cmd.Flag("base-path").Value.String(),
		HTTPUser: viper.GetString("HTTP_USER"),
		HTTPPass: viper.GetString("HTTP_PASS"),
	}
	server.Run(paramsWebsocket, c)
}
