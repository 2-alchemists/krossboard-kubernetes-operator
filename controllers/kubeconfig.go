/*
   Copyright (C) 2020  2ALCHEMISTS SAS.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as
   published by the Free Software Foundation, either version 3 of the
   License, or (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package controllers

import (
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"

	kclient "k8s.io/client-go/tools/clientcmd"
	kapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	AuthTypeUnknown     = 0
	AuthTypeBearerToken = 1
	AuthTypeX509Cert    = 2
	AuthTypeBasicToken  = 3
)

// KubeConfigManager holds an object describing a K8s Cluster.
type KubeConfigManager struct {
	Paths []string `json:"path,omitempty"`
}

// ManagedCluster holds an object describing managed clusters.
type ManagedCluster struct {
	Name        string         `json:"name,omitempty"`
	APIEndpoint string         `json:"apiEndpoint,omitempty"`
	AuthInfo    *kapi.AuthInfo `json:"authInfo,omitempty"`
	CaData      []byte         `json:"cacert,omitempty"`
	AuthType    int            `json:"authType,omitempty"`
}

// GetManagedClusters lists Kubernetes clusters available in KUBECONFIG.
func (m *KubeConfigManager) GetManagedClusters() map[string]*ManagedCluster {
	managedClusters := make(map[string]*ManagedCluster)
	for _, path := range m.Paths {
		config, err := kclient.LoadFromFile(path)
		if err != nil {
			log.WithError(err).Errorln("failed reading KUBECONFIG", path)
			continue
		}

		managedClusterSets := m.GetManagedClustersFromConfig(config)

		for k, v := range managedClusterSets {
			managedClusters[k] = v
		}
	}
	return managedClusters
}

// GetManagedClustersFromData lists Kubernetes clusters from a provided KUBECONFIG data.
func (m *KubeConfigManager) GetManagedClustersFromData(data []byte) (map[string]*ManagedCluster, error) {
	config, err := kclient.Load(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed loading KUBECONFIG data")
	}
	return m.GetManagedClustersFromConfig(config), nil
}

// GetManagedClustersFromConfig lists Kubernetes clusters from a provided KUBECONFIG.
func (m *KubeConfigManager) GetManagedClustersFromConfig(config *kapi.Config) map[string]*ManagedCluster {
	managedClusters := make(map[string]*ManagedCluster)

	// FIXME: check usefulness
	authInfos := make(map[string]string)
	for user, authInfo := range config.AuthInfos {
		authInfos[user] = authInfo.Token
	}

	for clusterName, clusterInfo := range config.Clusters {
		clusterNameEscaped := strings.ReplaceAll(clusterName, "/", "@")
		managedClusters[clusterNameEscaped] = &ManagedCluster{
			Name:        clusterNameEscaped,
			APIEndpoint: clusterInfo.Server,
			CaData:      clusterInfo.CertificateAuthorityData,
		}
	}
	for _, context := range config.Contexts {
		clusterNameEscaped := strings.ReplaceAll(context.Cluster, "/", "@")
		if cluster, found := managedClusters[clusterNameEscaped]; found {
			cluster.AuthInfo = config.AuthInfos[context.AuthInfo]
		}
	}

	return managedClusters
}

// GetAccessToken retrieves access token from AuthInfo.
func (m *KubeConfigManager) GetAccessToken(authInfo *kapi.AuthInfo) (string, error) {
	if authInfo == nil {
		return "", errors.New("no AuthInfo provided")
	}

	if authInfo.Token != "" {
		return authInfo.Token, nil // auth with Bearer token
	}

	authHookCmd := ""
	var args []string
	switch {
	case authInfo.AuthProvider != nil:
		authHookCmd = authInfo.AuthProvider.Config["cmd-path"]
		args = strings.Split(authInfo.AuthProvider.Config["cmd-args"], " ")
	case authInfo.Exec != nil:
		authHookCmd = authInfo.Exec.Command
		args = authInfo.Exec.Args
	default:
		return "", errors.New("no AuthInfo command provided")
	}

	cmd := exec.Command(authHookCmd, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, string(out))
	}

	token, err := jsonparser.GetString(out, "credential", "access_token") // GKE and alike
	if err != nil {
		errOut := errors.Wrap(err, "credentials string not compliant with GKE")
		token, err = jsonparser.GetString(out, "status", "token") // EKS and alike
		if err != nil {
			return "", errors.Wrap(errOut, "credentials string not compliant with EKS")
		}
	}

	return token, nil
}
