/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"sigs.k8s.io/gateway-api/gwctl/pkg/common/resourcehelpers"
	"sigs.k8s.io/gateway-api/gwctl/pkg/effectivepolicy"
	"sigs.k8s.io/gateway-api/gwctl/pkg/utils"
	"sigs.k8s.io/gateway-api/gwctl/pkg/utils/printer"
)

func NewGetCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "get {policies|policycrds|httproutes}",
		Short: "Display one or many resources",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			params := getParams(kubeConfigPath)
			runGet(cmd, args, params)
		},
	}
	cmd.Flags().StringVarP(&namespaceFlag, "namespace", "n", "default", "")
	cmd.Flags().BoolVarP(&allNamespacesFlag, "all-namespaces", "A", false, "If present, list requested resources from all namespaces.")

	return cmd
}

func runGet(cmd *cobra.Command, args []string, params *utils.CmdParams) {
	kind := args[0]
	ns, err := cmd.Flags().GetString("namespace")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read flag \"namespace\": %v\n", err)
		os.Exit(1)
	}

	allNs, err := cmd.Flags().GetBool("all-namespaces")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read flag \"all-namespaces\": %v\n", err)
		os.Exit(1)
	}

	if allNs {
		ns = ""
	}

	epc := effectivepolicy.NewCalculator(params.K8sClients, params.PolicyManager)
	policiesPrinter := &printer.PoliciesPrinter{Out: params.Out}
	httpRoutesPrinter := &printer.HTTPRoutesPrinter{Out: params.Out, EPC: epc}

	switch kind {
	case "policy", "policies":
		list := params.PolicyManager.GetPolicies()
		policiesPrinter.Print(list)

	case "policycrds":
		list := params.PolicyManager.GetCRDs()
		policiesPrinter.PrintCRDs(list)

	case "httproute", "httproutes":
		list, err := resourcehelpers.ListHTTPRoutes(context.TODO(), params.K8sClients, ns)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to list HTTPRoute resources: %v\n", err)
			os.Exit(1)
		}
		httpRoutesPrinter.Print(list)

	default:
		fmt.Fprintf(os.Stderr, "Unrecognized RESOURCE_TYPE\n")
	}
}
