/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package relay_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/hyperledger-labs/fabric-smart-client/integration"
)

func TestEndToEnd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Two Fabric Networks Suite with Cacti Weaver Relay")
}

func StartPort() int {
	return integration.TwoFabricNetworksWithWeaverRelayPort.StartPortForNode()
}
