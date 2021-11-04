package main

import (
	"net"
	"os"

	"gioui.org/app"
	"github.com/katzenpost/katzenpost/catshadow"
	"github.com/katzenpost/katzenpost/core/log"
	"github.com/katzenpost/katzenpost/client"
	"github.com/katzenpost/katzenpost/client/config"
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

func setupCatShadow(passphrase []byte, result chan interface{}) {
	// XXX: if the catshadowClient already exists, shut it down
	// FIXME: figure out a better way to toggle connected/disconnected
	// states and allow to retry attempts on a timeout or other failure.
	var catshadowClient *catshadow.Client
	var stateWorker *catshadow.StateWriter
	var state *catshadow.State
	var err error

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

	var cfg *config.Config
	if len(*clientConfigFile) != 0 {
		cfg, err = config.LoadFile(*clientConfigFile)
		if err != nil {
			result <- err
			return
		}
	} else {
		// use the baked in configuration defaults if a configuration is not specified
		if hasTor() {
			cfg, err = getDefaultConfig()
		} else {
			cfg, err = getConfigNoTor()
		}

		if err != nil {
			result <- err
			return
		}
	}

	// initialize logging
        backendLog, err := log.New(cfg.Logging.File, cfg.Logging.Level, cfg.Logging.Disable)
        if err != nil {
		result <- err
                return
        }

	// automatically create a statefile if one does not already exist
	stateLogger := backendLog.GetLogger("catshadow_state")
	if _, err = os.Stat(statefile); os.IsNotExist(err) {
		stateWorker, err = catshadow.NewStateWriter(stateLogger, statefile, passphrase)
	} else {
		stateWorker, state, err = catshadow.LoadStateWriter(stateLogger, statefile, passphrase)
	}

	// catches any err above
	if err != nil {
		result <- err
		return
	}

	// NewEphemeralClientConfig requires network connectivity to fetch a
	// pki.Document and select a provider.
	// TODO: fixme so that clients cache a consensus or set of provider descriptors
	// which can be used to populate the config.Account section, or otherwise
	// allow a client to be started without having config.Account set...

	// create a client
	c, err := client.New(cfg)
	if err != nil {
		result <- err
		return
	}

	// Start the stateworker
	stateWorker.Start()

	catshadowClient, err = catshadow.New(backendLog, c, stateWorker, state)
	if err != nil {
		result <- err
		c.Shutdown()
		stateWorker.Halt()
		return
	}
	result <- catshadowClient
}
