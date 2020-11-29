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
)

func init() {
	viper.SetEnvPrefix(envPrefix)
	err := viper.BindEnv("BOOTSTRAP_SERVERS")

	if err != nil {
		log.Fatal("unable to bind environment bootstrap variables")
	}
	err = viper.BindEnv("TLS_ENABLED")
	if err != nil {
		log.Fatal("unable to bind tls enabled variables")
	}

	initBootstrapServersFlag()
	initTslEnabledFlag()

	cobra.OnInitialize(createClient)
}

var klient *pkg.Klient

func createClient() {

	tls := false
	if tlsEnabled == "" {
		tls = viper.GetBool("TLS_ENABLED")
	} else {
		tls = strings.EqualFold(tlsEnabled, "true")
	}

	sEnv := viper.GetString("BOOTSTRAP_SERVERS")
	servers := []string{""}
	if sEnv != "" {
		servers = strings.Split(sEnv, ",")
	} else {
		servers = strings.Split(bootstrapServers, ",")
	}

	k, err := pkg.NewKlient("us-west-2", servers, tls)

	if err != nil {
		log.Fatal("Unable to create a Kafka client", err)
	}
	klient = k
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