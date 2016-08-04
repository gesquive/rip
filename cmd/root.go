package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path"

	cli "github.com/gesquive/cli-log"
	"github.com/gesquive/rip/format"
	"github.com/spf13/cobra"
)

var cfgFile string
var displayVersion string

var logDebug bool
var showVersion bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:       "rip [flags] <host>[:<port>] <tcp|udp> <file_path> [<file_path>...]",
	Short:     "Sends a text file line by line to a remote host/port",
	Long:      `Sends a text file line by line to a remote host/port.`,
	ValidArgs: []string{"host:port", "tcp|udp", "file_path"},
	Run:       run,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	displayVersion = version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&logDebug, "debug", "D", false,
		"Write debug messages to console")
	RootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false,
		"Show the version and exit")

	RootCmd.PersistentFlags().MarkHidden("debug")
}

func initConfig() {
	if logDebug {
		cli.SetLogLevel(cli.LevelDebug)
	}
	if showVersion {
		cli.Info(displayVersion)
		os.Exit(0)
	}
	cli.Debug("Running with debug turned on")
}

func run(cmd *cobra.Command, args []string) {
	//TODO: Add ability to receive piped input
	if len(args) < 3 {
		cmd.Usage()
		os.Exit(1)
	}

	address := args[0]
	protocol := args[1]
	files := args[2:]

	for _, f := range files {
		err := sendTextFile(f, protocol, address)
		if err != nil {
			cli.Errorf("Failed to send '%s' to %s://%s\n", f, protocol, address)
			cli.Error(err.Error())
		}
	}
}

func sendTextFile(filePath string, network string, address string) (err error) {
	textFile, err := os.Open(filePath)
	defer textFile.Close()
	if err != nil {
		return err
	}
	fileInfo, err := textFile.Stat()
	if err != nil {
		return err
	}
	bytesRead := uint64(0)
	totalSize := uint64(fileInfo.Size())
	fileName := fmt.Sprintf("%15s", path.Base(filePath))

	scanner := bufio.NewScanner(textFile)
	scanner.Split(bufio.ScanLines)

	destPort, err := net.Dial(network, address)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		bytesRead += uint64(len(line)) + 1
		fmt.Fprintf(destPort, line)
		cli.Infof("\r%s : %s %s", fileName,
			format.Percent(bytesRead, totalSize), format.Progress(bytesRead, totalSize))
	}
	cli.Infof("\r%s : %s\n", fileName, cli.Green("%16s", "complete"))
	return err
}
