package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path"
	"runtime"
	"time"

	"go.uber.org/ratelimit"

	"github.com/gesquive/cli"
	"github.com/gesquive/rip/format"
	"github.com/spf13/cobra"
)

var (
	buildVersion = "v0.2.1-dev"
	buildCommit  = ""
	buildDate    = ""
)

var debug bool
var showVersion bool
var msgRate int

func main() {
	Execute()
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "rip [flags] <host>[:<port>] <tcp|udp> <file_path> [<file_path>...]",
	Short:            "Sends a text file line by line to a remote host/port",
	Long:             `Sends a text file line by line to a remote host/port.`,
	ValidArgs:        []string{"host:port", "tcp|udp", "file_path"},
	PersistentPreRun: preRun,
	Run:              run,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmd.SetHelpTemplate(fmt.Sprintf("%s\nVersion:\n  github.com/gesquive/rip %s\n",
		RootCmd.HelpTemplate(), buildVersion))
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false,
		"Write debug messages to console")
	RootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false,
		"Show the version and exit")

	RootCmd.PersistentFlags().IntVarP(&msgRate, "rate-limit", "r", -1,
		"Message rate allowed per second, use -1 for no limit")

	RootCmd.PersistentFlags().MarkHidden("debug")
}

func preRun(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Printf("github.com/gesquive/rip\n")
		fmt.Printf(" Version:    %s\n", buildVersion)
		if len(buildCommit) > 6 {
			fmt.Printf(" Git Commit: %s\n", buildCommit[:7])
		}
		if buildDate != "" {
			fmt.Printf(" Build Date: %s\n", buildDate)
		}
		fmt.Printf(" Go Version: %s\n", runtime.Version())
		fmt.Printf(" OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
	if debug {
		cli.SetPrintLevel(cli.LevelDebug)
	}
	cli.Debug("Running with debug turned on")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		os.Exit(1)
	}

	address := args[0]
	protocol := args[1]
	files := args[2:]
	cli.Info("Sending data to %s://%s", protocol, address)

	// Detect if user is piping in text
	fileInput, err := os.Stdin.Stat()
	if err != nil {
		cli.Error(err.Error())
		os.Exit(2)
	}

	pipeFound := fileInput.Mode()&os.ModeNamedPipe != 0
	if pipeFound {
		cli.Debug("Pipe input detected, sending")
		err := sendTextFile(os.Stdin, protocol, address, msgRate)
		if err != nil {
			cli.Error("Failed to send piped data")
			cli.Error(err.Error())
		}
	}

	errCount := 0
	for _, f := range files {
		textFile, err := os.Open(f)
		defer textFile.Close()
		if err != nil {
			cli.Errorf("Failed to send '%s'\n\r", f)
			cli.Error(err.Error())
			errCount++
			continue
		}
		err = sendTextFile(textFile, protocol, address, msgRate)
		if err != nil {
			cli.Errorf("Failed to send '%s' to %s://%s\n", f, protocol, address)
			cli.Error(err.Error())
			errCount++
		}
	}

	if !pipeFound && len(files) == 0 {
		cli.Warn("No data was piped into or specified on the command line.\n")
		cmd.Usage()
		os.Exit(1)
	} else if errCount == 0 {
		cli.Info("All files successfully sent.")
	} else {
		cli.Warn("There were some errors while sending files.")
	}
}

func sendTextFile(textFile *os.File, network string, address string, rate int) (err error) {
	fileInfo, err := textFile.Stat()
	if err != nil {
		return err
	}
	bytesRead := uint64(0)
	linesRead := uint64(0)
	totalSize := uint64(fileInfo.Size())
	fileName := path.Base(textFile.Name())

	scanner := bufio.NewScanner(textFile)
	scanner.Split(bufio.ScanLines)

	destPort, err := net.Dial(network, address)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	stats := make(chan struct{})
	var limiter ratelimit.Limiter
	if rate < 0 {
		limiter = ratelimit.NewUnlimited()
	} else {
		limiter = ratelimit.New(rate)
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				cli.Infof("\rtransfer: %s %6dpkts %s", fileName, linesRead, format.Progress(bytesRead, totalSize))
			case <-stats:
				ticker.Stop()
				return
			}
		}
	}()

	start := time.Now()
	for scanner.Scan() {
		line := scanner.Text()
		linesRead++
		bytesRead += uint64(len(line)) + 1
		limiter.Take()
		_, err = fmt.Fprint(destPort, line)
		if err != nil {
			cli.Errorf("\r")
			return err
		}
	}
	end := time.Now()

	close(stats)
	cli.Info(cli.Green("\rtransfer: Successfully sent '%s' (%d packets in %s secs) %.2f pkt/sec",
		fileName, linesRead, end.Sub(start).String(), float64(linesRead)/end.Sub(start).Seconds()))
	return err
}
