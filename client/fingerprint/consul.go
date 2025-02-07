package fingerprint

import (
	"fmt"
	"log"
	"strconv"
	"time"

	consul "github.com/hashicorp/consul/api"

	client "github.com/hashicorp/nomad/client/config"
	"github.com/hashicorp/nomad/nomad/structs"
)

// ConsulFingerprint is used to fingerprint the architecture
type ConsulFingerprint struct {
	logger *log.Logger
	client *consul.Client
}

// NewConsulFingerprint is used to create an OS fingerprint
func NewConsulFingerprint(logger *log.Logger) Fingerprint {
	f := &ConsulFingerprint{logger: logger}
	return f
}

func (f *ConsulFingerprint) Fingerprint(config *client.Config, node *structs.Node) (bool, error) {
	// Guard against uninitialized Links
	if node.Links == nil {
		node.Links = map[string]string{}
	}

	// Only create the client once to avoid creating too many connections to
	// Consul.
	if f.client == nil {
		address := config.ReadDefault("consul.address", "127.0.0.1:8500")
		timeout, err := time.ParseDuration(config.ReadDefault("consul.timeout", "10ms"))
		if err != nil {
			return false, fmt.Errorf("Unable to parse consul.timeout: %s", err)
		}

		consulConfig := consul.DefaultConfig()
		consulConfig.Address = address
		consulConfig.HttpClient.Timeout = timeout

		f.client, err = consul.NewClient(consulConfig)
		if err != nil {
			return false, fmt.Errorf("Failed to initialize consul client: %s", err)
		}
	}

	// We'll try to detect consul by making a query to to the agent's self API.
	// If we can't hit this URL consul is probably not running on this machine.
	info, err := f.client.Agent().Self()
	if err != nil {
		// Clear any attributes set by a previous fingerprint.
		f.clearConsulAttributes(node)
		return false, nil
	}

	node.Attributes["consul.server"] = strconv.FormatBool(info["Config"]["Server"].(bool))
	node.Attributes["consul.version"] = info["Config"]["Version"].(string)
	node.Attributes["consul.revision"] = info["Config"]["Revision"].(string)
	node.Attributes["consul.name"] = info["Config"]["NodeName"].(string)
	node.Attributes["consul.datacenter"] = info["Config"]["Datacenter"].(string)

	node.Links["consul"] = fmt.Sprintf("%s.%s",
		node.Attributes["consul.datacenter"],
		node.Attributes["consul.name"])

	return true, nil
}

// clearConsulAttributes removes consul attributes and links from the passed
// Node.
func (f *ConsulFingerprint) clearConsulAttributes(n *structs.Node) {
	delete(n.Attributes, "consul.server")
	delete(n.Attributes, "consul.version")
	delete(n.Attributes, "consul.revision")
	delete(n.Attributes, "consul.name")
	delete(n.Attributes, "consul.datacenter")
	delete(n.Links, "consul")
}

func (f *ConsulFingerprint) Periodic() (bool, time.Duration) {
	return true, 15 * time.Second
}
