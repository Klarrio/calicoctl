// Copyright (c) 2019 Tigera, Inc. All rights reserved.

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
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"

	"github.com/projectcalico/calicoctl/calicoctl/commands"
	"github.com/projectcalico/calicoctl/calicoctl/util"
)

func main() {
	name, desc := util.NameAndDescription()
	doc := fmt.Sprintf(`Usage:
  <BINARY_NAME> [options] <command> [<args>...]

    create       Create a resource by file, directory or stdin.
    replace      Replace a resource by file, directory or stdin.
    apply        Apply a resource by file, directory or stdin.  This creates a resource
                 if it does not exist, and replaces a resource if it does exists.
    patch        Patch a pre-exisiting resource in place.
    delete       Delete a resource identified by file, directory, stdin or resource type and
                 name.
    get          Get a resource identified by file, directory, stdin or resource type and
                 name.
    label        Add or update labels of resources.
    convert      Convert config files between different API versions.
    ipam         IP address management.
    node         Calico node management.
    version      Display the version of this binary.
    export       Export the Calico datastore objects for migration
    import       Import the Calico datastore objects for migration
    datastore    Calico datastore management.

Options:
  -h --help               Show this screen.
  -l --log-level=<level>  Set the log level (one of panic, fatal, error,
                          warn, info, debug) [default: panic]

Description:
  The %s is used to manage Calico network and security
  policy, to view and manage endpoint configuration, and to manage a Calico
  node instance.

  See '<BINARY_NAME> <command> --help' to read about a specific subcommand.
`, desc)

	// Replace all instances of BINARY_NAME with the name of the binary.
	doc = strings.ReplaceAll(doc, "<BINARY_NAME>", name)

	arguments, err := docopt.Parse(doc, nil, true, commands.VERSION_SUMMARY, true, false)
	if err != nil {
		if _, ok := err.(*docopt.UserError); ok {
			// the user gave us bad input
			fmt.Printf("Invalid option: '%s'. Use flag '--help' to read about a specific subcommand.\n", strings.Join(os.Args[1:], " "))
		}
		os.Exit(1)
	}

	if logLevel := arguments["--log-level"]; logLevel != nil {
		parsedLogLevel, err := log.ParseLevel(logLevel.(string))
		if err != nil {
			fmt.Printf("Unknown log level: %s, expected one of: \n"+
				"panic, fatal, error, warn, info, debug.\n", logLevel)
			os.Exit(1)
		} else {
			log.SetLevel(parsedLogLevel)
			log.Infof("Log level set to %v", parsedLogLevel)
		}
	}

	if arguments["<command>"] != nil {
		command := arguments["<command>"].(string)
		args := append([]string{command}, arguments["<args>"].([]string)...)

		var err error

		switch command {
		case "create":
			err = commands.Create(args)
		case "replace":
			err = commands.Replace(args)
		case "apply":
			err = commands.Apply(args)
		case "patch":
			err = commands.Patch(args)
		case "delete":
			err = commands.Delete(args)
		case "get":
			err = commands.Get(args)
		case "label":
			err = commands.Label(args)
		case "convert":
			err = commands.Convert(args)
		case "version":
			err = commands.Version(args)
		case "node":
			err = commands.Node(args)
		case "ipam":
			err = commands.IPAM(args)
		case "datastore":
			err = commands.Datastore(args)
		default:
			err = fmt.Errorf("Unknown command: %q\n%s", command, doc)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
}
