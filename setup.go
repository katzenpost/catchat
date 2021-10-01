package main

import (
	"fmt"
	"net"
	"os"

	"gioui.org/app"
	"github.com/katzenpost/katzenpost/catshadow"
	catconfig "github.com/katzenpost/katzenpost/catshadow/config"
	"github.com/katzenpost/katzenpost/client"
	"path/filepath"
)

// checks to see if the local system has a listener on port 9050
func hasTor() bool {
	c, err := net.Dial("tcp", "127.0.0.1:9050")
	if err != nil {
		return false
	}
	c.Close()
	return true
}

func setupCatShadow(catshadowCfg *catconfig.Config, passphrase []byte, result chan interface{}) {
	// XXX: if the catshadowClient already exists, shut it down
	// FIXME: figure out a better way to toggle connected/disconnected
	// states and allow to retry attempts on a timeout or other failure.
	var stateWorker *catshadow.StateWriter
	var state *catshadow.State
	cfg, err := catshadowCfg.ClientConfig()
	if err != nil {
		result <- err
		return
	}

	// obtain the default data location
	dir, err := app.DataDir()
	if err != nil {
		result <- err
		return
	}

	// dir does not appear to point to ~/.config/catchat but rather ~/.config on linux?
	// create directory for application data
	datadir := filepath.Join(dir, dataDirName)
	_, err = os.Stat(datadir)
	if os.IsNotExist(err) {
		// create the application data directory
		err := os.Mkdir(datadir, os.ModeDir|os.FileMode(0700))
		if err != nil {
			result <- err
			return
		}
	}

	// if the statefile doesn't exist, try the default datadir
	var statefile string
	if _, err := os.Stat(*stateFile); os.IsNotExist(err) {
		statefile = filepath.Join(datadir, *stateFile)
	} else {
		statefile = *stateFile
	}

	// initialize logging
	backendLog, err := catshadowCfg.InitLogBackend()
	if err != nil {
		result <- err
		return
	}

	var catshadowClient *catshadow.Client
	cfg, linkKey, err := client.NewEphemeralClientConfig(cfg)
	if err != nil {
		result <- err
		return
	}

	c, err := client.New(cfg)
	if err != nil {
		result <- err
		return
	}

	// automatically create a statefile if one does not already exist
	if _, err := os.Stat(statefile); os.IsNotExist(err) {
		// Create statefile.
		stateWorker, err = catshadow.NewStateWriter(c.GetLogger("catshadow_state"), statefile, passphrase)
		if err != nil {
			result <- err
			c.Shutdown()
			return
		}
		// Start the stateworker
		stateWorker.Start()

		// create a catshadow client and message spool
		fmt.Println("creating remote message receiver spool")
		user := fmt.Sprintf("%x", linkKey.PublicKey().Bytes())
		catshadowClient, err = catshadow.NewClientAndRemoteSpool(backendLog, c, stateWorker, user, linkKey)
		if err != nil {
			result <- err
			stateWorker.Halt()
			c.Shutdown()
			return
		}
		fmt.Println("catshadow client successfully created")
	} else {
		// Load previous state to setup our current client state.
		stateWorker, state, err = catshadow.LoadStateWriter(c.GetLogger("catshadow_state"), statefile, passphrase)
		if err != nil {
			result <- err
			return
		}

		// Start the stateworker
		stateWorker.Start()

		// Make a catshadow Client.
		catshadowClient, err = catshadow.New(backendLog, c, stateWorker, state)
		if err != nil {
			stateWorker.Halt()
			result <- err
			return
		}
	}

	// Start catshadow client.
	result <- catshadowClient
}
