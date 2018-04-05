// Copyright (c) 2018 Abhijit Gadgil <gabhijit@iitbombay.org>. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io"

	"github.com/vishvananda/netns"

	"github.com/osrg/gobgp/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	GOBGP_GRPC_SERVER = ":50051"
	MONITOR_FAMILY    = 0
)

type nodeAgentManager struct {
	routeListener *netlinkRouteListener
	bgpClient     gobgpapi.GobgpApiClient
	ctx           context.Context
}

func newNodeAgentManager() *nodeAgentManager {

	ns := netns.None()
	r := newNetlinkRouteListener(ns)

	//FIXME : We need to get this value from Config.
	conn, err := grpc.Dial(GOBGP_GRPC_SERVER, grpc.WithInsecure())
	if err != nil {
		return nil
	}
	bgpcl := gobgpapi.NewGobgpApiClient(conn)

	//FIXME: is this correct or should I get a context.Context()
	ctx := context.Background()

	agent := &nodeAgentManager{routeListener: r, bgpClient: bgpcl, ctx: ctx}

	return agent
}

func (n *nodeAgentManager) run() {

	// open netlink Writer

	// subscribe to RIB updates
	stream, err := n.bgpClient.MonitorRib(n.ctx, &gobgpapi.MonitorRibRequest{
		Table: &gobgpapi.Table{
			Type:   gobgpapi.Resource_GLOBAL,
			Family: MONITOR_FAMILY,
		},
		Current: false,
	})

	if err != nil {
		return
	}

	go n.ribMonitor(stream)

	// run netlink listener
	n.routeListener.run()
}

func (n *nodeAgentManager) ribMonitor(cl gobgpapi.GobgpApi_MonitorRibClient) {
	for {
		dst, err := cl.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			// FIXME: Error handling
			return

		} else {
			fmt.Println("%+v", dst)
		}
	}
}
