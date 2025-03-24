package main

import (
	"hostsparser/hosts"
	"hostsparser/logger"
	"log"

	"go.uber.org/zap"
)

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger.SetLogger(zapLogger.Sugar())
	hostsFile, err := hosts.LoadHostsFile(hosts.GetSystemHostsPath())
	if err != nil {
		log.Fatalf("Failed to load hosts file: %v", err)
	}

	ipDomainPairs := hostsFile.GetAllIPDomainPairs()

	for _, ipDomainPair := range ipDomainPairs {
		logger.Infof("ipDomainPair: %v\n", ipDomainPair)
	}
}
