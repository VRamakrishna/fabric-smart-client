/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package weaver

import (
	"github.com/hyperledger-labs/fabric-smart-client/integration/nwo/fabric/topology"
)

const (
	TopologyName = "weaver"
)

type Network struct {
	Type string
	Name string
}

type Driver struct {
	Name     string
	Hostname string
	Port     uint16 `yaml:"port,omitempty"`
}

type IINAgent struct {
	Name            string
	Hostname        string
	Port            uint16 `yaml:"port,omitempty"`
	Network         string
	Organization    string
}

type InteropChaincode struct {
	Label     string
	Channel   string
	Namespace string
	Path      string
}

type RelayServer struct {
	FabricTopology     *topology.Topology `yaml:"-"`
	FabricTopologyName string
	Name               string
	Hostname           string
	Port               uint16 `yaml:"port,omitempty"`
	Organization       string
	Networks           []*Network
	Drivers            []*Driver
	InteropChaincode   InteropChaincode
}

func (w *RelayServer) AddFabricNetwork(ft *topology.Topology) *RelayServer {
	w.Networks = append(w.Networks, &Network{
		Type: "Fabric",
		Name: ft.Name(),
	})
	return w
}

type Topology struct {
	TopologyName string `yaml:"name,omitempty"`
	TopologyType string `yaml:"type,omitempty"`
	Relays       []*RelayServer
	IINAgents    []*IINAgent
}

func NewTopology() *Topology {
	return &Topology{
		TopologyName: TopologyName,
		TopologyType: TopologyName,
		Relays:       []*RelayServer{},
		IINAgents:    []*IINAgent{},
	}
}

func (t *Topology) Name() string {
	return t.TopologyName
}

func (t *Topology) Type() string {
	return t.TopologyType
}

func (t *Topology) AddRelayServer(ft *topology.Topology, org string) *RelayServer {
	ft.EnableWeaver()
	r := &RelayServer{
		FabricTopology:     ft,
		FabricTopologyName: ft.Name(),
		Name:               ft.Name(),
		Hostname:           "relay-" + ft.Name(),
		Organization:       org,
		Drivers: []*Driver{
			{
				Name:     "Fabric",
				Hostname: "driver-" + ft.Name(),
			},
		},
		InteropChaincode: InteropChaincode{
			Label:     "interop",
			Channel:   ft.Channels[0].Name,
			Namespace: "interop",
			Path:      "github.com/hyperledger/cacti/weaver/core/network/fabric-interop-cc/contracts/interop/v2",
		},
	}
	t.Relays = append(t.Relays, r)
	r.AddFabricNetwork(ft)
	return r
}

func (t *Topology) AddIINAgent(ft *topology.Topology, org string) *IINAgent {
	ft.EnableWeaver()
	ia := &IINAgent{
		Name:           "Fabric-" + org,
		Hostname:       "iinagent-" + ft.Name() + "-" + org,
		Network:        ft.Name(),
		Organization:   org,
	}
	t.IINAgents = append(t.IINAgents, ia)
	return ia
}

func (t *Topology) AddIINAgents(ft *topology.Topology) []*IINAgent {
	for _, organization := range ft.Organizations {
		t.AddIINAgent(ft, organization.Name)
	}
	return t.IINAgents
}
