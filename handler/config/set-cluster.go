/*
Copyright © 2020 Portworx

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
package configcli

import (
	"fmt"
	"io/ioutil"

	"github.com/portworx/pxc/pkg/commander"
	"github.com/portworx/pxc/pkg/config"
	"github.com/portworx/pxc/pkg/util"
	"github.com/spf13/cobra"
)

// setClusterCmd represents the config command
var (
	setClusterCmd *cobra.Command
	setCluster    *config.Cluster
)

var _ = commander.RegisterCommandVar(func() {
	setCluster = config.NewCluster()
	setClusterCmd = &cobra.Command{
		Use:   "set-cluster",
		Short: "Setup pxc cluster configuration",
		RunE:  setClusterExec,
	}
})

var _ = commander.RegisterCommandInit(func() {
	ConfigAddCommand(setClusterCmd)

	setClusterCmd.Flags().StringVar(&setCluster.TunnelServiceNamespace,
		"portworx-service-namespace", "kube-system", "Kubernetes namespace for the Portworx service")
	setClusterCmd.Flags().StringVar(&setCluster.TunnelServiceName,
		"portworx-service-name", "portworx-api", "Kubernetes name for the Portworx service")
	setClusterCmd.Flags().StringVar(&setCluster.TunnelServicePort,
		"portworx-service-port", "9020", "Port for the Portworx SDK endpoint in the Kubernetes service")
	setClusterCmd.Flags().BoolVar(&setCluster.Secure,
		"tls", false, "Enable if using TLS. Passing a CA will enable this automatically.")
	setClusterCmd.Flags().StringVar(&setCluster.CACert,
		"cafile", "", "Path to CA certificate")
	setClusterCmd.Flags().StringVar(&setCluster.Endpoint,
		"endpoint", "", "Direct connection to a Portworx node gRPC endpoint. "+
			"This endpoint would be used instead of the Kubernetes Portworx API service. "+
			"Example: 1.1.1.1:9020")
})

func setClusterExec(cmd *cobra.Command, args []string) error {
	cc := config.KM().ToRawKubeConfigLoader()

	// This is the raw kubeconfig which may have been overridden by CLI args
	kconfig, err := cc.RawConfig()
	if err != nil {
		return err
	}

	// Get the current context
	currentContextName, err := config.GetKubernetesCurrentContext()
	if err != nil {
		return err
	}

	// Get the current context object
	currentContext := kconfig.Contexts[currentContextName]

	// Initialize cluster object
	setCluster.Name = currentContext.Cluster

	// Validate the value of endpoint if provided
	if len(setCluster.Endpoint) != 0 {
		var err error
		setCluster.Endpoint, err = util.ValidateEndpoint(setCluster.Endpoint)
		if err != nil {
			return fmt.Errorf("Invalid endpoint: %s", setCluster.Endpoint)
		}
	}

	// Check if CA was provided
	if len(setCluster.CACert) != 0 {
		var err error
		if !util.IsFileExists(setCluster.CACert) {
			return fmt.Errorf("CA file: %s does not exists", setCluster.CACert)
		}
		setCluster.CACertData, err = ioutil.ReadFile(setCluster.CACert)
		if err != nil {
			return fmt.Errorf("Unable to read CA file %s", setCluster.CACert)
		}
		// If CA cert is provided, make "secure" default to true
		setCluster.Secure = true
	}

	// Get the location of the kubeconfig for this specific object. This is necessary
	// because KUBECONFIG can have many kubeconfigs, example: KUBECONFIG=kube1.conf:kube2.conf
	location := kconfig.Clusters[currentContext.Cluster].LocationOfOrigin

	// Storage the information to the appropriate kubeconfig
	if err := config.SaveClusterInKubeconfig(currentContext.Cluster, location, setCluster); err != nil {
		return err
	}

	util.Printf("Portworx server information saved in %s for Kubernetes cluster %s\n",
		location,
		currentContext.Cluster)
	return nil
}
