/*
 * Copyright 2016 The Linux Foundation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>

#include "cmsg.h"

#define error(fmt, ...)							\
	({								\
		fprintf(stderr, "nsenter: " fmt ": %m\n", ##__VA_ARGS__); \
		errno = ECOMM;						\
		-1; /* return value */					\
	})

/*
 * Sends a file descriptor along the sockfd provided. Returns the return
 * value of sendmsg(2). Any synchronisation and preparation of state
 * should be done external to this (we expect the other side to be in
 * recvfd() in the code).
 */
ssize_t sendfd(int sockfd, int fd)
{
	struct msghdr msg = {0};
	struct cmsghdr *cmsg;
	int *fdptr;

	union {
		char buf[CMSG_SPACE(sizeof(&fd))];
		struct cmsghdr align;
	} u;

	msg.msg_control = u.buf;
	msg.msg_controllen = sizeof(u.buf);

	cmsg = CMSG_FIRSTHDR(&msg);
	cmsg->cmsg_level = SOL_SOCKET;
	cmsg->cmsg_type = SCM_RIGHTS;
	cmsg->cmsg_len = CMSG_LEN(sizeof(int));

	fdptr = (int *) CMSG_DATA(cmsg);
	memcpy(fdptr, &fd, sizeof(int));

	return sendmsg(sockfd, &msg, 0);
}

/*
 * Receives a file descriptor from the sockfd provided. Returns the file
 * descriptor as sent from sendfd(). It will return the file descriptor
 * or die (literally) trying. Any synchronisation and preparation of
 * state should be done external to this (we expect the other side to be
 * in sendfd() in the code).
 */
int recvfd(int sockfd)
{
	struct msghdr msg = {0};
	struct cmsghdr *cmsg;
	int *fdptr;

	ssize_t ret = recvmsg(sockfd, &msg, 0);
	if (ret < 0)
		return ret;

	cmsg = CMSG_FIRSTHDR(&msg);
	if (cmsg->cmsg_level != SOL_SOCKET)
		return error("recvfd: expected SOL_SOCKET in cmsg: %d", cmsg->cmsg_level);
	if (cmsg->cmsg_type != SCM_RIGHTS)
		return error("recvfd: expected SCM_RIGHTS in cmsg: %d", cmsg->cmsg_type);
	if (cmsg->cmsg_len != CMSG_LEN(sizeof(int)))
		return error("recvfd: expected correct CMSG_LEN in cmsg: %d", cmsg->cmsg_len);

	fdptr = (int *) CMSG_DATA(cmsg);
	if (!fdptr || *fdptr < 0)
		return error("recvfd: unexpected failure to get fd");

	return *fdptr;
}
