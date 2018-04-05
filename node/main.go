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

// compago : Node Agent

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/osrg/gobgp/api"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"google.golang.org/grpc"
)

const (
	GOBGP_GRPC_SERVER = ":50051"
)

type NodeAgentInfo struct {
	routeListener *NetlinkRouteListener
	bgpClient     gobgpapi.GobgpApiClient
}

func NewNodeAgentInfo() *NodeAgentInfo {

	//FIXME : We need to get this value from Config.
	conn, err := grpc.Dial(GOBGP_GRPC_SERVER, grpc.WithInsecure())
	if err != nil {
		return nil
	}
	bgpcl := gobgpapi.NewGobgpApiClient(conn)

	ns := netns.None()
	r := NewNetlinkRouteListener(ns)

	agent := &NodeAgentInfo{routeListener: r, bgpClient: bgpcl}

	return agent
}

func (agent *NodeAgentInfo) Run() {
	agent.routeListener.Run()
}

type NetlinkRouteListener struct {
	ns   netns.NsHandle
	ch   chan netlink.RouteUpdate
	done chan struct{}
}

func NewNetlinkRouteListener(ns netns.NsHandle) *NetlinkRouteListener {

	var newns netns.NsHandle

	if ns == 0 {
		newns = netns.None()
	} else {
		newns = ns
	}

	ch := make(chan netlink.RouteUpdate)
	done := make(chan struct{})
	n := &NetlinkRouteListener{ch: ch, done: done, ns: newns}

	return n
}

func (nl *NetlinkRouteListener) Run() {

	if nl == nil {
		return
	}

	defer close(nl.done)
	err := netlink.RouteSubscribeAt(nl.ns, nl.ch, nl.done)
	if err != nil {
		os.Exit(-1)
	}

	for {
		select {
		case update := <-nl.ch:
			// FIXME: handle this route update
			fmt.Println("%q", update)
		default:
			break
		}
	}

}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)

	var opts struct {
		ConfigFile string `short:"f" long:"config-file" description:"specify a config file"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(-1)
	}

	agent := NewNodeAgentInfo()

	agent.Run()

	os.Exit(0)
}
