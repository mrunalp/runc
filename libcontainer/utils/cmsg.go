// +build linux
package utils

/*
#include <errno.h>
#include <stdlib.h>
#include "cmsg.h"
*/
import "C"
import "os"

// RecvFd waits for a file descriptor to be sent over the given AF_UNIX
// socket. The file name of the remote file descriptor will be recreated
// locally (it is sent as non-auxilliary data in the same payload).
func RecvFd(socket *os.File) (*os.File, error) {
	file, err := C.recvfd(C.int(socket.Fd()))
	if err != nil {
		return nil, err
	}
	defer C.free(file.tag)
	return os.NewFile(uintptr(file.fd), C.GoString(file.tag)), nil
}

// SendFd sends a file descriptor over the given AF_UNIX socket. In
// addition, the file.Name() of the given file will also be sent as
// non-auxilliary data in the same payload (allowing to send contextual
// information for a file descriptor).
func SendFd(socket, file *os.File) error {
	var cfile C.struct_file_t
	cfile.fd = C.int(file.Fd())
	cfile.tag = C.CString(file.Name())
	defer C.free(cfile.tag)

	_, err := C.sendfd(C.int(socket.Fd()), cfile)
	return err
}
