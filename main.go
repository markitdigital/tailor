package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	// TimeoutSeconds is the number of seconds to wait between service status queries
	TimeoutSeconds = 20

	// DelayOnExit is the number of seconds to wait after a service has exited before returning an exit code to allow logs
	// to finigh being collected
	DelayOnExit = 5
)

// LogLine contains a log line and a source file
type LogLine struct {
	SourceFile string
	Line       string
}

func init() {
	// set logrus to use JSON & stdout
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {

	// exit with usage if we aren't passed any arguments
	flag.Parse()
	if flag.NArg() == 0 {
		cwd, _ := os.Getwd()
		fmt.Printf("USAGE: %s\\tailor.exe [windows service name] {log file locations or globs (optional) ...}\n", cwd)
		os.Exit(0)
	}

	// treat the first positional argument as the service name, additional args as logfile paths
	serviceName := flag.Arg(0)
	filePaths := []string{}
	if flag.NArg() > 1 {
		filePaths = flag.Args()[1:]
	}

	// log a startup message
	logrus.WithField("source", "tailor").Debugf("Monitoring service: %s", serviceName)

	// create a goroutine per logfile which tails and writes lines to stdout
	for _, path := range filePaths {
		matches, _ := filepath.Glob(path)
		for _, file := range matches {
			go func(file string) {
				logrus.WithField("source", "tailor").Debugf("Tailing logfile: %s", file)
				t, _ := tail.TailFile(file, tail.Config{Follow: true, Poll: true, ReOpen: true, Logger: logrus.WithFields(logrus.Fields{"source": file})})
				for line := range t.Lines {
					logrus.WithField("source", file).Debug(line.Text)
				}
			}(file)
		}
	}

	// continually ping the service, exit when the service is not in a running state
	for {

		// connect the service manager
		m, err := mgr.Connect()
		if err != nil {
			logrus.WithField("source", "tailor").Debugf("unable to connect Windows service manager: %v", err)
		}
		defer m.Disconnect()

		// open the specified service
		s, err := m.OpenService(serviceName)
		if err != nil {
			logrus.WithField("source", "tailor").Debugf("could not access service: %v", err)
		}
		defer s.Close()

		// query for the service's status
		status, err := s.Query()
		if err != nil {
			logrus.WithField("source", "tailor").Debugf("could not query service: %v", err)
		}

		// if the service is not running, delay for a bit to collect any additional logs, then exit
		if status.State != svc.Running {
			time.Sleep(time.Second * DelayOnExit)
			os.Exit(1)
		}

		// delay between service status checks
		time.Sleep(time.Second * TimeoutSeconds)

	}

}
