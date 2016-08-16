package utils

/*
#include <errno.h>
#include "cmsg.h"
*/
import "C"
import "os"

func RecvFd(socket *os.File) (*os.File, error) {
	fd, err := C.recvfd(C.int(socket.Fd()))
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), "[recvfd]"), nil
}

func SendFd(socket, file *os.File) error {
	_, err := C.sendfd(C.int(socket.Fd()), C.int(file.Fd()))
	return err
}
