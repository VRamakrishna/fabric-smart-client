/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package id

import (
	"github.com/pkg/errors"

	"github.com/hyperledger-labs/fabric-smart-client/platform/view/driver"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/services/flogging"
	kms "github.com/hyperledger-labs/fabric-smart-client/platform/view/services/kms"
	"github.com/hyperledger-labs/fabric-smart-client/platform/view/view"
)

var logger = flogging.MustGetLogger("view-sdk.id")

//go:generate counterfeiter -o mock/config_provider.go -fake-name ConfigProvider . ConfigProvider

type ConfigProvider interface {
	GetPath(s string) string
	GetStringSlice(key string) []string
	TranslatePath(path string) string
}

//go:generate counterfeiter -o mock/sig_service.go -fake-name SigService . SigService

type SigService interface {
	RegisterSigner(identity view.Identity, signer driver.Signer, verifier driver.Verifier) error
}

type EndpointService interface {
	GetIdentity(label string, pkid []byte) (view.Identity, error)
}

type provider struct {
	configProvider  ConfigProvider
	sigService      SigService
	endpointService EndpointService
	defaultID       view.Identity
	admins          []view.Identity
	clients         []view.Identity
	kms             *kms.KMS
}

func NewProvider(configProvider ConfigProvider, sigService SigService, endpointService EndpointService, kms *kms.KMS) *provider {
	return &provider{
		configProvider:  configProvider,
		sigService:      sigService,
		endpointService: endpointService,
		kms:             kms,
	}
}

func (p *provider) Load() error {
	if err := p.loadDefaultIdentity(); err != nil {
		return errors.WithMessagef(err, "failed loading default identity")
	}

	if err := p.loadClientIdentities(); err != nil {
		return errors.WithMessagef(err, "failed loading client identities")
	}

	return nil
}

func (p *provider) DefaultIdentity() view.Identity {
	return p.defaultID
}

func (p *provider) Identity(label string) view.Identity {
	id, err := p.endpointService.GetIdentity(label, nil)
	if err != nil {
		logger.Warningf("failed to get identity for label %s: %s", label, err)
		return nil
	}
	return id
}

func (p *provider) Admins() []view.Identity {
	return p.admins
}

func (p *provider) Clients() []view.Identity {
	return p.clients
}

func (p *provider) loadDefaultIdentity() error {
	id, signer, verifier, err := p.kms.Load(p.configProvider)
	if err != nil {
		return errors.Wrapf(err, "failed loading default signer")
	}

	if err := p.sigService.RegisterSigner(id, signer, verifier); err != nil {
		return errors.Wrapf(err, "failed registering default identity signer")
	}
	p.defaultID = id
	return nil
}

func (p *provider) loadClientIdentities() error {
	certs := p.configProvider.GetStringSlice("fsc.client.certs")
	var clients []view.Identity
	for _, cert := range certs {
		// TODO: support cert as a folder
		certPath := p.configProvider.TranslatePath(cert)
		client, err := LoadIdentity(certPath)
		if err != nil {
			logger.Errorf("failed loading client cert at [%s]: [%s]", certPath, err)
			continue
		}
		logger.Infof("loaded client cert at [%s]: [%s]", certPath, err)
		clients = append(clients, client)
	}
	logger.Infof("loaded [%d] client identities", len(clients))
	p.clients = clients
	return nil
}
