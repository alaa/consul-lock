package consul

import (
	"errors"
	"fmt"
	"log"
	"os"

	consulapi "github.com/hashicorp/consul/api"
	multierror "github.com/hashicorp/go-multierror"
)

type Consul struct {
	Client       *consulapi.Client
	LockPath     string
	RevisionPath string
}

func New(lockPath, revisionPath string) (*Consul, error) {
	consulAddr := os.Getenv("CONSUL_ADDR")
	if consulAddr == "" {
		log.Print("CONSUL_ADDR is not set, using localhost:8500")
		consulAddr = "localhost:8500"
	}

	if lockPath == "" || revisionPath == "" {
		log.Fatal("Please set lockPath and revisionPath correctly")
	}

	config := &consulapi.Config{
		Address: consulAddr,
		Scheme:  "http",
	}

	client, err := consulapi.NewClient(config)
	if err != nil {
		return &Consul{}, err
	}

	return &Consul{
		Client:       client,
		LockPath:     lockPath,
		RevisionPath: revisionPath,
	}, nil
}

func (c *Consul) createSession() (string, error) {
	session := c.Client.Session()
	sessionID, _, err := session.Create(&consulapi.SessionEntry{
		Behavior: consulapi.SessionBehaviorDelete,
		TTL:      "8h",
	}, nil)
	return sessionID, err
}

func (c *Consul) destroySession(sessionID string) error {
	session := c.Client.Session()
	_, err := session.Destroy(sessionID, nil)
	return err
}

func (c *Consul) Status() error {
	err := c.isLocked()
	if err != nil {
		return err
	}
	return nil
}

func (c *Consul) AcquireLock() (string, error) {
	err := c.isLocked()
	if err != nil {
		return "", err
	}

	sessionID, err := c.createSession()
	if err != nil {
		return "", err
	}

	kv := c.Client.KV()
	kvpair := &consulapi.KVPair{
		Key:     c.LockPath,
		Session: sessionID,
		Value:   []byte(sessionID),
	}
	kv.Acquire(kvpair, nil)

	return sessionID, nil
}

func (c *Consul) isLocked() error {
	kv := c.Client.KV()
	pair, _, err := kv.Get(c.LockPath, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not fetch lock %s", pair))
	}

	if pair != nil && pair.Session != "" {
		return errors.New(fmt.Sprintf("Lock session is locked with sessionID: %s", pair.Session))
	}

	return nil
}

func (c *Consul) GetRevision() (string, error) {
	kv := c.Client.KV()
	pair, _, err := kv.Get(c.RevisionPath, nil)
	if err != nil || pair == nil {
		return "", err
	}
	return string(pair.Value), nil
}

func (c *Consul) UpdateRevision(id string) error {
	kv := c.Client.KV()
	_, err := kv.Put(&consulapi.KVPair{
		Key:   c.RevisionPath,
		Value: []byte(id),
	}, nil)
	return err
}

func (c *Consul) ReleaseLock(sessionID string) error {
	var result error

	kv := c.Client.KV()
	_, _, err := kv.Release(&consulapi.KVPair{
		Key:     c.LockPath,
		Session: sessionID,
	}, nil)
	if err != nil {
		result = multierror.Append(result, err)
	}

	err = c.destroySession(sessionID)
	if err != nil {
		result = multierror.Append(result, err)
	}

	return result
}
