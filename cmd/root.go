package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "rip <host>[:<port>] <tcp|udp> <file_path>",
	Short: "Sends a text file line by line to a remote host/port",
	Long: `Sends a text file line by line to a remote host/port.
	The data sent will not include line endings.`,
	Run: run,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	//TODO: add version flag
	//TODO: add debug logging
}

func initConfig() {
}

func run(cmd *cobra.Command, args []string) {
	if len(args) != 3 {
		cmd.Usage()
		os.Exit(1)
	}

	address := args[0]
	protocol := args[1]
	file := args[2]

	sendTextFile(file, protocol, address)
}

func sendTextFile(path string, network string, address string) (err error) {
	textFile, err := os.Open(path)
	defer textFile.Close()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(textFile)
	scanner.Split(bufio.ScanLines)

	destPort, err := net.Dial(network, address)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(line)
		// TODO: Add progress bar!
		fmt.Fprintf(destPort, line)
	}
	return err
}
