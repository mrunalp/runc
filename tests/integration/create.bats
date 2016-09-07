#!/usr/bin/env bats

load helpers

function setup() {
	teardown_busybox
	setup_busybox
}

function teardown() {
	teardown_busybox
}

@test "runc create" {
	sed -i 's/"terminal": true,/"terminal": false,/' config.json

	runc create -d test_busybox
	[ "$status" -eq 0 ]

	testcontainer test_busybox created

	# start the command
	runc start test_busybox
	[ "$status" -eq 0 ]

	testcontainer test_busybox running
}

@test "runc create exec" {
	sed -i 's/"terminal": true,/"terminal": false,/' config.json

	runc create -d test_busybox
	[ "$status" -eq 0 ]

	testcontainer test_busybox created

	runc exec test_busybox true
	[ "$status" -eq 0 ]

	# start the command
	runc start test_busybox
	[ "$status" -eq 0 ]

	testcontainer test_busybox running
}
