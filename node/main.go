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

	"os"
	"os/signal"
	"syscall"

	"github.com/jessevdk/go-flags"
)

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

	os.Exit(0)
}
