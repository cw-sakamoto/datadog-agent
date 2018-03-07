// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/DataDog/datadog-agent/pkg/collector/corechecks/cluster"
	_ "github.com/DataDog/datadog-agent/pkg/collector/corechecks/network"
	_ "github.com/DataDog/datadog-agent/pkg/collector/corechecks/system"
	"github.com/DataDog/datadog-agent/pkg/config"
	log "github.com/cihub/seelog"

	"github.com/DataDog/datadog-agent/cmd/cluster-agent/app"
)

func main() {
	// go_expvar server
	go http.ListenAndServe(
		fmt.Sprintf("127.0.0.1:%d", config.Datadog.GetInt("clusteragent_expvar_port")),
		http.DefaultServeMux)

	if err := app.ClusterAgentCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
