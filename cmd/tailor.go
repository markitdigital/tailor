package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hpcloud/tail"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/markitdigital/tailor"
)

// create a channel to hold all the log messages
var messages chan *tailor.Message

func main() {

	// exit with usage if we aren't passed any arguments
	flag.Parse()
	if flag.NArg() < 1 {
		cwd, _ := os.Getwd()
		fmt.Printf("USAGE: %s\\tailor.exe [windows service name] {log file locations or globs (optional) ...}\n", cwd)
		os.Exit(0)
	}

	// treat the first positional argument as the service name, additional args as paths
	service, paths := flag.Arg(0), []string{}
	if flag.NArg() > 1 {
		paths = flag.Args()[1:]
	}

	// initialize message
	messages = make(chan *tailor.Message)

	// start goroutine that tails files and writes messages to channel
	go tailFiles(paths)

	// start goroutine that writes messages to stdout as structured JSON
	go writeMessages()

	// connect the service manager
	m, err := mgr.Connect()
	if err != nil {
		messages <- tailor.NewMessage("tailor", fmt.Sprintf("unable to connect Windows service manager: %v", err))
	}
	defer m.Disconnect()

	// open the specified service
	s, err := m.OpenService(service)
	if err != nil {
		messages <- tailor.NewMessage("tailor", fmt.Sprintf("could not access service: %v", err))
	}
	defer s.Close()

	// monitor the service
	for {

		// query for the service's status
		status, err := s.Query()
		if err != nil {
			messages <- tailor.NewMessage("tailor", fmt.Sprintf("could not query service: %v", err))
		}

		// if the service is not running, delay for a bit to collect any additional logs, then exit
		if status.State != svc.Running {
			time.Sleep(time.Second * 30)
			os.Exit(1)
		}

		// delay between service status checks
		time.Sleep(time.Second * 30)

	}

}

func tailFiles(paths []string) {
	foundLogFiles := map[string]int{}
	for {
		for _, path := range paths {
			files, _ := filepath.Glob(path)
			for _, file := range files {
				if _, ok := foundLogFiles[file]; !ok {
					foundLogFiles[file] = 1
					go func() {
						messages <- tailor.NewMessage("tailor", fmt.Sprintf("Tailing file: %s", file))
						t, _ := tail.TailFile(file, tail.Config{Follow: true, Poll: true, ReOpen: true})
						for line := range t.Lines {
							messages <- tailor.NewMessage(file, line.Text)
						}
					}()
				}
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func writeMessages() {
	for {
		message := <-messages
		bytes, _ := json.Marshal(message)
		fmt.Println(string(bytes))
	}
}
