#!/bin/bash


# functions for setting up following topology

##############################################################################
#
#
#   +-------+           +-------+    IPv6     +-------+            +-------+
#   | node1 |<--------->| ctrl1 |<----------->| ctrl2 |<---------->| node2 |
#   +-------+           +-------+             +-------+            +-------+
#
##############################################################################


IP=ip

function _do_create_controller {

	ctrller=$1

	${IP} netns add ${ctrller}
}

function controller_controller_link {

	ctrller1=$1
	ctrller2=$2

	# Controller controller link
	${IP} link add ${ctrller1}-eth type veth peer name ${ctrller2}-eth

	for i in `seq 1 2`; do
		if [ $i -eq 1 ]; then
			ctrller=${ctrller1}
		else
			ctrller=${ctrller2}
		fi
		${IP} link set ${ctrller}-eth netns ${ctrller}
		${IP} netns exec ${ctrller} ${IP} link set ${ctrller}-eth name eth0
		${IP} netns exec ${ctrller} ${IP} link set eth0 up
		${IP} netns exec ${ctrller} sysctl -w net.ipv6.conf.eth0.disable_ipv6=0

		sleep 2
		# ping to find neighbours ping6 doesn't work inside a netns easily, so we do this trick
		${IP} netns exec ${ctrller} bash -c 'ping6 -c 2 ff02::1\%eth0' &
	done
}

function setup_controllers {

	for i in `seq 1 2`; do
		_do_create_controller ctrl${i}
	done
	controller_controller_link ctrl1 ctrl2
}

function del_controllers {
	for i in `seq 1 2`; do
		${IP} netns del ctrl${i}
	done
}
