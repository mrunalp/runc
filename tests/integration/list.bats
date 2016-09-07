#!/usr/bin/env bats

load helpers

function setup() {
	teardown_running_container_inroot test_box1 $HELLO_BUNDLE
	teardown_running_container_inroot test_box2 $HELLO_BUNDLE
	teardown_running_container_inroot test_box3 $HELLO_BUNDLE
	teardown_busybox
	setup_busybox
}

function teardown() {
	teardown_running_container_inroot test_box1 $HELLO_BUNDLE
	teardown_running_container_inroot test_box2 $HELLO_BUNDLE
	teardown_running_container_inroot test_box3 $HELLO_BUNDLE
	teardown_busybox
}

@test "list" {
	sed -i 's/"terminal": true,/"terminal": false,/' config.json

	# run a few busyboxes detached
	ROOT=$HELLO_BUNDLE runc run -d test_box1
	[ "$status" -eq 0 ]
	wait_for_container_inroot 15 1 test_box1 $HELLO_BUNDLE

	ROOT=$HELLO_BUNDLE runc run -d test_box2
	[ "$status" -eq 0 ]
	wait_for_container_inroot 15 1 test_box2 $HELLO_BUNDLE

	ROOT=$HELLO_BUNDLE runc run -d test_box3
	[ "$status" -eq 0 ]
	wait_for_container_inroot 15 1 test_box3 $HELLO_BUNDLE

	ROOT=$HELLO_BUNDLE runc list
	[ "$status" -eq 0 ]
	[[ ${lines[0]} =~ ID\ +PID\ +STATUS\ +BUNDLE\ +CREATED+ ]]
	[[ "${lines[1]}" == *"test_box1"*[0-9]*"running"*$BUSYBOX_BUNDLE*[0-9]* ]]
	[[ "${lines[2]}" == *"test_box2"*[0-9]*"running"*$BUSYBOX_BUNDLE*[0-9]* ]]
	[[ "${lines[3]}" == *"test_box3"*[0-9]*"running"*$BUSYBOX_BUNDLE*[0-9]* ]]

	ROOT=$HELLO_BUNDLE runc list --format table
	[ "$status" -eq 0 ]
	[[ ${lines[0]} =~ ID\ +PID\ +STATUS\ +BUNDLE\ +CREATED+ ]]
	[[ "${lines[1]}" == *"test_box1"*[0-9]*"running"*$BUSYBOX_BUNDLE*[0-9]* ]]
	[[ "${lines[2]}" == *"test_box2"*[0-9]*"running"*$BUSYBOX_BUNDLE*[0-9]* ]]
	[[ "${lines[3]}" == *"test_box3"*[0-9]*"running"*$BUSYBOX_BUNDLE*[0-9]* ]]

	ROOT=$HELLO_BUNDLE runc list --format json
	[ "$status" -eq 0 ]
	[[ "${lines[0]}" == [\[][\{]"\"id\""[:]"\"test_box1\""[,]"\"pid\""[:]*[0-9][,]"\"status\""[:]*"\"running\""[,]"\"bundle\""[:]*$BUSYBOX_BUNDLE*[,]"\"created\""[:]*[0-9]*[\}]* ]]
	[[ "${lines[0]}" == *[,][\{]"\"id\""[:]"\"test_box2\""[,]"\"pid\""[:]*[0-9][,]"\"status\""[:]*"\"running\""[,]"\"bundle\""[:]*$BUSYBOX_BUNDLE*[,]"\"created\""[:]*[0-9]*[\}]* ]]
	[[ "${lines[0]}" == *[,][\{]"\"id\""[:]"\"test_box3\""[,]"\"pid\""[:]*[0-9][,]"\"status\""[:]*"\"running\""[,]"\"bundle\""[:]*$BUSYBOX_BUNDLE*[,]"\"created\""[:]*[0-9]*[\}][\]] ]]
}
