package libcontainer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/opencontainers/runc/libcontainer/utils"
)

type syncType uint8

// Constants that are used for synchronisation between the parent and child
// during container setup. They come in pairs (with procError being a generic
// response which is followed by a &genericError).
//
// [  child  ] <-> [   parent   ]
//
// procHooks   --> [run hooks]
//             <-- procResume
//
// procReady   --> [final setup]
//             <-- procRun
const (
	procError syncType = iota
	procReady
	procRun
	procHooks
	procResume
)

type syncT struct {
	Type syncType `json:"type"`
}

// Used to write to a synchronisation pipe. An error is returned if there was
// a problem writing the payload.
func writeSync(pipe io.Writer, sync syncType) error {
	if err := utils.WriteJSON(pipe, syncT{sync}); err != nil {
		return err
	}
	return nil
}

// Used to read from a synchronisation pipe. An error is returned if we got a
// genericError, the pipe was closed, or we got an unexpected flag.
func readSync(pipe io.Reader, expected syncType) error {
	var procSync syncT
	if err := json.NewDecoder(pipe).Decode(&procSync); err != nil {
		if err == io.EOF {
			return fmt.Errorf("parent closed synchronisation channel")
		}

		if procSync.Type == procError {
			var ierr genericError

			if err := json.NewDecoder(pipe).Decode(&ierr); err != nil {
				return fmt.Errorf("failed reading error from parent: %v", err)
			}

			return &ierr
		}

		if procSync.Type != expected {
			return fmt.Errorf("invalid synchronisation flag from parent")
		}
	}
	return nil
}
