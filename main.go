package main

import (
	"log"

	"github.com/alaa/consul-lock/cache"
	"github.com/alaa/consul-lock/consul"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	lock            = kingpin.Command("lock", "Accquire lock from consul")
	status          = kingpin.Command("status", "Check lock status")
	release         = kingpin.Command("release", "Release the lock, uses .consul_lock_id file")
	release_with_id = kingpin.Command("release-with-id", "Release the lock with explicit session ID argument. Should not be used normally")
	releaseID       = release_with_id.Arg("id", "session ID to release").Required().String()
)

func main() {
	lockPath, revisionPath := "terraform/lock", "terraform/revision"
	sessionDir := "./"

	consul, err := consul.New(lockPath, revisionPath)
	if err != nil {
		log.Fatal(err)
	}

	cache, err := cache.New(sessionDir)
	if err != nil {
		log.Fatal(err)
	}

	switch kingpin.Parse() {
	case "lock":
		sid, err := consul.AcquireLock()
		if err != nil {
			log.Fatal(err)
		}
		cache.UpdateSession(sid)
		log.Printf("Lock Session ID: %s", sid)

	case "status":
		err := consul.Status()
		if err != nil {
			log.Fatal("Lock is currently accquired. Please wait. %s", err)
		}
		log.Printf("Lock status is available")

	case "release":
		releaseID, err := cache.GetSession()
		if err != nil {
			log.Fatal(err)
		}

		if err := consul.ReleaseLock(releaseID); err != nil {
			log.Fatal(err)
		}
		log.Print("released")

	case "release-with-id":
		if err := consul.ReleaseLock(*releaseID); err != nil {
			log.Fatal(err)
		}
		log.Print("released")
	}
}
