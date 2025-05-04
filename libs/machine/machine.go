package machine

import (
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func GetMachineID() string {
	if uuid := GetSystemUUID(); uuid != "" {
		return "uuid:" + uuid
	}

	if mac := GetMACAddress(); mac != "" {
		return "mac:" + mac
	}

	hostname, _ := os.Hostname()
	return "host:" + hostname
}

func GetSystemUUID() string {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("cat", "/sys/class/dmi/id/product_uuid")
	case "windows":
		cmd = exec.Command("wmic", "csproduct", "get", "UUID")
	default:
		return ""
	}

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return ""
	}

	// Windows returns two lines, first is "UUID", second is the value
	if runtime.GOOS == "windows" && len(lines) > 1 {
		return strings.TrimSpace(lines[1])
	}

	return strings.TrimSpace(lines[0])
}

func GetMACAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback == 0 && iface.HardwareAddr != nil {
			mac := iface.HardwareAddr.String()
			if mac != "" {
				return mac
			}
		}
	}

	return ""
}
