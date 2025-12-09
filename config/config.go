package config

import "fmt"

const (
	TargetHost = "192.168.0.92"
	TargetPort = 5010
	ProxyPort  = 5000
)

var TargetURL = fmt.Sprintf("http://%s:%d", TargetHost, TargetPort)
