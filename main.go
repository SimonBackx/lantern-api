package main

import (
	"flag"
	"github.com/kardianos/service"
	"log"
)

var logger service.Logger

// Program structures.
// Define Start and Stop methods.
type program struct {
	exit     chan bool
	finished chan bool
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan bool)
	p.finished = make(chan bool)

	// Start should not block. Do the actual work async.
	go run(p.exit, p.finished)
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	//logger.Info("I'm Stopping!")
	p.exit <- true

	// Wachten tot stoppen voltooid is
	<-p.finished
	return nil
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "Lantern-Api",
		DisplayName: "Lantern-Api",
		Description: "Cyber threat collection in the darkweb",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
