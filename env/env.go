package env

import "os"

var (
	psm    = UNDEFINED
	dsAddr = UNDEFINED
	psmIP  = UNDEFINED
)

const (
	EnvKeyPSM = "SERVICE_PSM"

	EnvKeyDSAddr = "DS_ADDR"

	EnvKeyPSMIP = "SERVICE_PSM_IP"

	UNDEFINED = "undefined"
)

func init() {
	psm = os.Getenv(EnvKeyPSM)
	dsAddr = os.Getenv(EnvKeyDSAddr)
	psmIP = os.Getenv(EnvKeyPSMIP)
}

func PSM() string {
	if IsLocal() {
		return "whalethinker.test.consul"
	}
	return psm
}

func DSAddr() string {
	if IsLocal() {
		return "localhost:8500"
	}
	return dsAddr
}

func IsLocal() bool {
	return os.Getenv("IS_LOCAL") == "true"
}

func PSMIP() string {
	if IsLocal() {
		return "localhost"
	}
	return psmIP
}
