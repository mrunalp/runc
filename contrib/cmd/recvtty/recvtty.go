package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/opencontainers/runc/libcontainer/utils"
)

func bail(err error) {
	fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
	os.Exit(1)
}

func handler(path string) error {
	// Open a socket.
	ln, err := net.Listen("unix", path)
	if err != nil {
		return err
	}
	defer ln.Close()

	// We only accept a single connection, since we can only really have
	// one reader for os.Stdin. Plus this is all a PoC.
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Close ln, to allow for other instances to take over.
	ln.Close()

	// Get the fd of the connection.
	unixconn, ok := conn.(*net.UnixConn)
	if !ok {
		bail(fmt.Errorf("failed to cast to unixconn"))
	}

	socket, err := unixconn.File()
	if err != nil {
		bail(err)
	}
	defer socket.Close()

	// Get the master file descriptor from runC.
	master, err := utils.RecvFd(socket)
	if err != nil {
		bail(err)
	}

	// TODO: Print information about what container the socket comes from.
	fmt.Fprintf(os.Stderr, "recv: masterfd(%d)\n", master.Fd())

	// Copy from our stdio to the master fd.
	quitChan := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, master)
		quitChan <- struct{}{}
	}()
	go func() {
		io.Copy(master, os.Stdin)
		quitChan <- struct{}{}
	}()

	// Only close the master fd once we've stopped copying.
	<-quitChan
	master.Close()
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: recvtty <path>\n")
		os.Exit(1)
	}

	path := os.Args[1]
	if err := handler(path); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		os.Exit(1)
	}
}
