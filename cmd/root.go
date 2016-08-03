package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path"

	log "github.com/gesquive/cli-log"
	"github.com/gesquive/rip/format"
	"github.com/spf13/cobra"
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
		log.Error(err)
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
			log.Errorf("Failed to send '%s' to %s://%s\n", f, protocol, address)
			log.Error(err.Error())
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
		log.Infof("\r%s : %s %s", fileName,
			format.Percent(bytesRead, totalSize), format.Progress(bytesRead, totalSize))
	}
	log.Infof("\r%s : %s\n", fileName, log.Green("%16s", "complete"))
	return err
}
