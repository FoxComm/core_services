package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var enviroments = map[string]string{
	"production": "foxcomm123",
	"staging":    "foxcomm-stage",
}

var (
	// env is the given env arg
	env string

	// host filters based on tags (for ansible)
	host string

	withCoreOsInfo bool
	unitRp         *regexp.Regexp
	spaceRp        *regexp.Regexp
)

func init() {
	flag.StringVar(&env, "env", "staging", "Specify gcloud --project (environment). Defaults to 'staging'. INVENTORY env var overrides this flag.")
	flag.StringVar(&host, "host", "", "")
	flag.Bool("list", true, "")
	flag.BoolVar(&withCoreOsInfo, "with-coreos", false, "Include CoreOs info")
	if withCoreOsEnv := os.Getenv("WITH_COREOS"); withCoreOsEnv != "" {
		withCoreOsInfo = true
	}

	unitRp = regexp.MustCompile(`(.+)_((release-[0-9-]+)|(v[0-9.]+))@([0-9]+)\.service`)
	spaceRp = regexp.MustCompile(`\s+`)
}

type Instances []Instance

type Instance struct {
	ID                string
	Kind              string
	Name              string
	Status            string
	MachineType       string
	NetworkInterfaces []NetworkInterface
	TagItems          struct {
		Items []string
	} `json:"tags"`
	Zone       string
	Containers []Container
}

type NetworkInterface struct {
	Name          string
	Network       string
	IP            string `json:"networkIP"`
	AccessConfigs []AccessConfig
}

type AccessConfig struct {
	NatIP string
}

// PublicIP returns the NAT/publically accessible IP for this instance
func (i Instance) PublicIP() string {
	if len(i.NetworkInterfaces) > 0 && len(i.NetworkInterfaces[0].AccessConfigs) > 0 {
		return i.NetworkInterfaces[0].AccessConfigs[0].NatIP
	} else {
		return ""
	}
}

// PrivateIP returns the internal/private IP for this instance
func (i Instance) PrivateIP() string {
	if len(i.NetworkInterfaces) > 0 {
		return i.NetworkInterfaces[0].IP
	} else {
		return ""
	}
}

func (i Instance) Tags() []string {
	return i.TagItems.Items
}

// IsActive is it RUNNING
func (i Instance) IsActive() bool {
	return strings.ToUpper(i.Status) == "RUNNING"
}

func (i Instances) ToJSON() ([]byte, error) {
	capRoles := map[string]interface{}{}
	ansibleHostVars := map[string]map[string]interface{}{}

	for _, instance := range i {
		ip := instance.PrivateIP()
		ansibleHostVars[ip] = map[string]interface{}{
			"private_ip_address": instance.PrivateIP(),
			"public_ip_address":  instance.PublicIP(),
			"ansible_node_name":  instance.Name,
		}
		if withCoreOsInfo {
			units := []map[string]interface{}{}
			for _, unit := range instance.Containers {
				units = append(units, map[string]interface{}{
					"cid":       unit.ConcurrentId,
					"name":      unit.Name,
					"version":   unit.Version,
					"endpoints": unit.Endpoints,
				})
			}
			ansibleHostVars[ip]["containers"] = units

			if stringInSlice("cluster", instance.Tags()) {
				ansibleHostVars[ip]["ansible_ssh_user"] = "core"
				ansibleHostVars[ip]["ansible_python_interpreter"] = "/home/core/bin/python"
			}
		}

		for _, role := range instance.Tags() {
			if _, ok := capRoles[role]; !ok {
				capRoles[role] = []string{}
			}

			capRoles[role] = append(capRoles[role].([]string), instance.PrivateIP())
		}
	}

	output := capRoles
	output["_meta"] = map[string]interface{}{"hostvars": ansibleHostVars}

	return json.Marshal(output)
}

func List(environment string) (Instances, error) {
	var project string
	var ok bool

	if project, ok = enviroments[environment]; !ok {
		return nil, fmt.Errorf("Could not find environment %s", environment)
	}

	cmd := exec.Command("gcloud", "compute", "instances", "list", "--format", "json", "--project", project)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	return decode(out)
}

type Endpoint struct {
	Name string
	Port int
}

type Container struct {
	ConcurrentId int
	Version      string
	Name         string
	Ip           string
	Endpoints    []Endpoint
}

func (c *Container) FoxCommCalculateEndpoints() error {
	var portBase string
	c.Endpoints = []Endpoint{}

	type tmpEndpoint struct {
		Name     string
		PortBase string
	}
	tmpInfo := []tmpEndpoint{}

	switch c.Name {
	case "backups":
		portBase = "1700"
		tmpInfo = append(tmpInfo, tmpEndpoint{c.Name, portBase})
	case "catalog_cache":
		portBase = "1300"
		tmpInfo = append(tmpInfo, tmpEndpoint{c.Name, portBase})
	case "causes_cache":
		portBase = "1600"
		tmpInfo = append(tmpInfo, tmpEndpoint{c.Name, portBase})
	case "feature_manager":
		portBase = "1100"
		tmpInfo = append(tmpInfo, tmpEndpoint{"core", portBase})
	case "router":
		return nil
	case "social_analytics":
		portBase = "1000"
		tmpInfo = append(tmpInfo, tmpEndpoint{c.Name, portBase})
		tmpInfo = append(tmpInfo, tmpEndpoint{"loyalty_engine", portBase})
	case "social_shopping":
		portBase = "1200"
		tmpInfo = append(tmpInfo, tmpEndpoint{c.Name, portBase})
	case "ui", "foxcomm-ui":
		portBase = "600"
		tmpInfo = append(tmpInfo, tmpEndpoint{c.Name, portBase})
	case "user":
		portBase = "1500"
		tmpInfo = append(tmpInfo, tmpEndpoint{c.Name, portBase})
	default:
		return fmt.Errorf("Unknown service: %s", c.Name)
	}

	for _, tmpEp := range tmpInfo {
		portStr := fmt.Sprintf("%s%d", tmpEp.PortBase, c.ConcurrentId)
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}
		ep := Endpoint{Name: tmpEp.Name, Port: port}
		c.Endpoints = append(c.Endpoints, ep)
	}

	return nil
}

func parseFleetUnit(unit string, machine string) (*Container, error) {
	// parse unit
	matches := unitRp.FindAllStringSubmatch(unit, -1)

	if len(matches) < 1 {
		return nil, fmt.Errorf("Can't parse unit: %s", unit)
	}
	submatches := matches[0]
	if len(submatches) < 4 {
		return nil, fmt.Errorf("Can't parse unit: %s", unit)
	}

	version, err := strconv.Atoi(submatches[5])
	if err != nil {
		return nil, err
	}

	// parse machine
	ms := strings.Split(machine, "/")
	if len(ms) < 2 {
		return nil, fmt.Errorf("Can't parse machine: %s", machine)
	}
	return &Container{
		Version:      submatches[2], // 3 also
		Name:         submatches[1],
		ConcurrentId: version,
		Ip:           ms[1],
	}, nil
}

// If we decide to expose fleetd we can use go client github.com/coreos/fleet/client
// Now, we use cmd line tool which connects to fleetd by ssh via FLEETCTL_TUNNEL
// and parse given results
func getFleetContainers() ([]Container, error) {
	containers := []Container{}

	if os.Getenv("FLEETCTL_TUNNEL") == "" {
		return containers, fmt.Errorf("Please set FLEETCTL_TUNNEL to a machine in the cluster and ensure you have access to that machine in order to continue.")
	}
	cmd := exec.Command("fleetctl", "list-units", "-no-legend", "-fields=unit,machine")
	if out, err := cmd.CombinedOutput(); err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			if line == "" {
				continue
			}

			ss := spaceRp.Split(line, -1)
			if len(ss) < 2 {
				return containers, fmt.Errorf("Can't parse unit line: %s", line)
			}

			unit, err := parseFleetUnit(ss[0], ss[1])
			if err == nil {
				err = unit.FoxCommCalculateEndpoints()
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "skip unit, %s\n", err.Error())
				continue
			}
			containers = append(containers, *unit)

		}
	} else {
		return containers, fmt.Errorf("%s: %s", err.Error(), out)
	}
	return containers, nil
}

func (is Instances) LoadUnits() error {
	units, err := getFleetContainers()
	if err != nil {
		return err
	}
	m := map[string][]Container{}
	for _, u := range units {
		if _, ok := m[u.Ip]; !ok {
			m[u.Ip] = []Container{}
		}
		m[u.Ip] = append(m[u.Ip], u)
	}

	for i, it := range is {
		ip := it.PrivateIP()
		if units, ok := m[ip]; ok {
			instance := &is[i]
			instance.Containers = units
		}
	}

	return nil
}

func decode(r io.Reader) (Instances, error) {
	instances := Instances{}
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&instances); err != nil {
		return nil, err
	}

	running := make(Instances, 0, len(instances))
	for _, inst := range instances {
		if inst.IsActive() {
			running = append(running, inst)
		}
	}
	return running, nil
}

func main() {
	flag.Parse()
	var result []byte
	var err error

	instances, err := List(env)
	if err != nil {
		panic(err)
	}

	if host != "" {
		filtered := make(Instances, 0, 0)

		for _, i := range instances {
			for _, tag := range i.Tags() {
				if tag == host {
					filtered = append(filtered, i)
				}
			}
		}

		instances = filtered
	}

	if withCoreOsInfo {
		if err := instances.LoadUnits(); err != nil {
			panic(err)
		}
	}

	result, err = instances.ToJSON()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(result))
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
