package main

import (
	"fmt"
	"log"

	"github.com/alaa/consul-lock/consul"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	lock      = kingpin.Command("lock", "Accquire lock from consul")
	status    = kingpin.Command("status", "Check lock status")
	release   = kingpin.Command("release", "Release the lock, i.e: Unlock")
	releaseID = release.Arg("id", "session ID to release").Required().String()
)

func main() {
	lockPath, revisionPath := "terraform/lock", "terraform/revision"

	consul, err := consul.New(lockPath, revisionPath)
	if err != nil {
		log.Fatal(err)
	}

	switch kingpin.Parse() {
	case "lock":
		sid, err := consul.AcquireLock()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Lock Session ID: %s", sid)

	case "status":
		err := consul.Status()
		if err != nil {
			log.Fatal("Lock is currently accquired. Please wait. %s", err)
		}
		log.Printf("Lock status is available")

	case "release":
		fmt.Println(*releaseID)
		if err := consul.ReleaseLock(*releaseID); err != nil {
			log.Fatal(err)
		}
		log.Print("released")

	}
}
