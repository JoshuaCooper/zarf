// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package cluster contains Zarf-specific cluster management functions.
package cluster

import (
	"time"

	"github.com/defenseunicorns/zarf/src/config"
	"github.com/defenseunicorns/zarf/src/pkg/k8s"
	"github.com/defenseunicorns/zarf/src/pkg/message"
)

// Cluster is a wrapper for the k8s package that provides Zarf-specific cluster management functions.
type Cluster struct {
	Kube *k8s.K8s
}

const (
	defaultTimeout = 30 * time.Second
	agentLabel     = "zarf.dev/agent"
)

var labels = k8s.Labels{
	config.ZarfManagedByLabel: "zarf",
}

// NewClusterOrDie creates a new cluster instance and waits up to 30 seconds for the cluster to be ready or throws a fatal error.
func NewClusterOrDie() *Cluster {
	c, err := NewClusterWithWait(defaultTimeout)
	if err != nil {
		message.Fatalf(err, "Failed to connect to cluster")
	}

	return c
}

// NewClusterWithWait creates a new cluster instance and waits for the given timeout for the cluster to be ready.
func NewClusterWithWait(timeout time.Duration) (*Cluster, error) {
	c := &Cluster{}
	var err error
	c.Kube, err = k8s.New(message.Debugf, labels)
	if err != nil {
		return c, err
	}
	return c, c.Kube.WaitForHealthyCluster(timeout)
}

// NewCluster creates a new cluster instance without waiting for the cluster to be ready.
func NewCluster() (*Cluster, error) {
	c := &Cluster{}
	c.Kube, _ = k8s.New(message.Debugf, labels)
	return c, nil
}
