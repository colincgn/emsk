package cmd

import (
	"github.com/colincgn/emsk/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"strings"
)

const (
	// The environment variable prefix of all environment variables bound to our command line flags.
	// For example, --number is bound to EMSK_NUMBER.
	envPrefix = "EMSK"
	defaultBootstrapServers = "localhost:29092"
	tlsEnabledEnvVariable = "TLS_ENABLED"
	bootstrapServersEnvVariable = "BOOTSTRAP_SERVERS"
)

func init() {
	viper.SetEnvPrefix(envPrefix)
	err := viper.BindEnv(bootstrapServersEnvVariable)

	if err != nil {
		log.Fatal("unable to bind environment bootstrap variables")
	}
	err = viper.BindEnv(tlsEnabledEnvVariable)
	if err != nil {
		log.Fatal("unable to bind tls enabled variables")
	}

	initBootstrapServersFlag()
	initTslEnabledFlag()

	cobra.OnInitialize(createKafkaClient)
}

var kafka *pkg.Kafka

func createKafkaClient() {

	tls := false
	if tlsEnabled == "" {
		tls = viper.GetBool(tlsEnabledEnvVariable)
	} else {
		tls = strings.EqualFold(tlsEnabled, "true")
	}

	sEnv := viper.GetString(bootstrapServersEnvVariable)
	servers := []string{""}
	if sEnv != "" {
		servers = strings.Split(sEnv, ",")
	} else {
		servers = strings.Split(bootstrapServers, ",")
	}

	k, err := pkg.NewKafka("us-west-2", servers, tls)

	if err != nil {
		log.Fatal("Unable to create a Kafka client", err)
	}
	kafka = k
}

var tlsEnabled string
func initTslEnabledFlag() {
	rootCmd.PersistentFlags().StringVarP(&tlsEnabled, "tls-enabled","t", "", "boolean flag if brokers are TLS enabled or not")
}

var bootstrapServers string
func initBootstrapServersFlag() {
	rootCmd.PersistentFlags().StringVarP(&bootstrapServers, "bootstrap-servers","s", defaultBootstrapServers, "a comma separated list of bootstrap-servers")
}

var rootCmd = &cobra.Command{
	Use:   "emsk",
	Short: "emsk",
	Long: `
Command line utility application to help with simple Kafka commands when working with MSK and AWS Lambda.
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}