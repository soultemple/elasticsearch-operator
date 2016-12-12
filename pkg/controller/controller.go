/*
Copyright (c) 2016, UPMC Enterprises
All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name UPMC Enterprises nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL UPMC ENTERPRISES BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
*/

package controller

import (
	"fmt"
	"traefik/log"

	"github.com/upmc-enterprises/elasticsearch-operator/pkg/cluster"
	"github.com/upmc-enterprises/elasticsearch-operator/util/k8sutil"
	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/rest"
	"k8s.io/client-go/1.4/tools/clientcmd"
)

type Config struct {
	Namespace  string
	MasterHost string
}

type Controller struct {
	Config
	kclient  *kubernetes.Clientset
	clusters map[string]*cluster.ElasticSearchCluster
}

func New(name, ns, kubeCfgFile string) (*Controller, error) {
	var (
		client     *kubernetes.Clientset
		masterHost string
	)

	// Should we use in cluster or out of cluster config
	if len(kubeCfgFile) == 0 {
		log.Info("Using InCluster k8s config")
		cfg, err := rest.InClusterConfig()

		if err != nil {
			return nil, err
		}

		masterHost = cfg.Host
		client, err = kubernetes.NewForConfig(cfg)

		if err != nil {
			return nil, err
		}
	} else {
		log.Infof("Using OutOfCluster k8s config with kubeConfigFile: %s", kubeCfgFile)
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeCfgFile)

		if err != nil {
			log.Error("Got error trying to create client: ", err)
			return nil, err
		}

		masterHost = cfg.Host
		client, err = kubernetes.NewForConfig(cfg)

		if err != nil {
			return nil, err
		}
	}

	c := &Controller{
		kclient: client,
		Config: Config{
			Namespace:  ns,
			MasterHost: masterHost,
		},
		clusters: make(map[string]*cluster.ElasticSearchCluster),
	}

	return c, nil
}

func (c *Controller) Run() error {

	_, err := c.init()

	if err != nil {
		log.Error("Error in init(): ", err)
	}
	return nil
}

func (c *Controller) init() (string, error) {
	err := k8sutil.CreateKubernetesThirdPartyResource(c.MasterHost)
	if err != nil {
		return "", err
	}

	return "", fmt.Errorf("fail to create TPR: %v", err)
}
