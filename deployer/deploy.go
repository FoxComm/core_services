package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var Root, _ = os.Getwd()

type Service struct {
	Name   string
	Config string
}

func (service Service) ConfigName() string {
	paths := strings.Split(service.Config, "/")
	return paths[len(paths)-1]
}

func (service Service) Path() string {
	paths := strings.Split(service.Config, "/")
	return strings.Join(paths[0:len(paths)-1], "/")
}

func AllServices() (services []Service) {
	if matches, err := filepath.Glob("*/*@.service"); err == nil {
		for _, match := range matches {
			service := Service{Config: match}
			service.Name = strings.TrimSuffix(service.ConfigName(), "@.service")
			services = append(services, service)
		}
	}
	return
}

func (service Service) Compile() bool {
	os.Setenv("CGO_ENABLED", "0")
	os.Setenv("GOOS", "linux")

	currentDir, _ := os.Getwd()
	os.Chdir(service.Path())
	defer os.Chdir(currentDir)
	cOut, err := exec.Command("go", "build", "-a", "-installsuffix", "cgo").CombinedOutput()
	if *debug {
		fmt.Printf("Attempted to build %s module with following results: %s\n", currentDir, cOut)
	}
	return err == nil
}

func (service Service) Release(version string) bool {
	currentDir, _ := os.Getwd()
	os.Chdir(service.Path())
	defer os.Chdir(currentDir)

	if out, err := exec.Command("docker", "build", "-t", service.Name, ".").CombinedOutput(); err == nil {
		fmt.Printf("Building docker image for %s...\n", service.Name)
		out, _ := exec.Command("docker", "images").Output()
		for _, str := range strings.Split(string(out), "\n") {
			if matched, _ := regexp.MatchString(service.Name, str); matched {
				imageID := strings.Fields(str)[2]
				imageName := fmt.Sprintf("quay.io/foxcomm/%s:%s", service.Name, version)
				fmt.Printf("Pushing %s to %s...\n", imageID, imageName)
				if _, err := exec.Command("docker", "tag", imageID, imageName).Output(); err == nil {
					if _, err := exec.Command("docker", "push", imageName).Output(); err == nil {
						return true
					} else {
						fmt.Printf("Failure with docker push of image: %s\n", imageName)
						fmt.Printf("Failure Details: %s\n", err)
					}
				} else {
					fmt.Printf("Failure with docker tag of image: %s\n", imageName)
					fmt.Printf("Failure Details: %s\n", err)
				}
			}
		}
	} else {
		fmt.Println("Encountered an error building docker images.")
		fmt.Println(string(out))
	}
	return false
}

func (service Service) Destroy() bool {
	if out, err := exec.Command("fleetctl", "list-unit-files").Output(); err == nil {
		for _, str := range strings.Split(string(out), "\n") {
			deployName := strings.TrimSuffix(service.DeployName(), "@.service")
			// ignore other services
			if matched, _ := regexp.MatchString(service.Name, str); matched {
				// don't handle current deploy version
				if matched, _ := regexp.MatchString(deployName+"@", str); !matched || (matched && *overwrite) {
					name := strings.Fields(str)[0]

					// backup service before destroy
					if strings.HasSuffix(name, "@.service") {
						if out, err := exec.Command("fleetctl", "cat", name).Output(); err == nil {
							backupPath := path.Join(*backup, name)
							if err := ioutil.WriteFile(backupPath, out, 0644); err == nil {
								fmt.Printf("Backed up %s into %s\n", service.ConfigName(), backupPath)
							} else {
								fmt.Printf("Get error when write backup %s: %s\n", service.ConfigName(), err.Error())
								fmt.Printf(string(out))
							}
						}
					}

					// destroy service
					if _, err := exec.Command("fleetctl", "destroy", name).Output(); err == nil {
						fmt.Println("Destroyed " + name)
					} else {
						fmt.Println("Failed to destroy " + name)
						return false
					}
				}
			}
		}
	}
	return true
}

func (service Service) Stop() bool {
	if out, err := exec.Command("fleetctl", "list-units").Output(); err == nil {
		for _, str := range strings.Split(string(out), "\n") {
			deployName := strings.TrimSuffix(service.DeployName(), "@.service")
			// ignore other services
			if matched, _ := regexp.MatchString(service.Name, str); matched {
				// don't handle current deploy version
				if matched, _ := regexp.MatchString(deployName+"@", str); !matched {
					name := strings.Fields(str)[0]

					// stop service
					if _, err := exec.Command("fleetctl", "stop", name).Output(); err == nil {
						fmt.Println("Stopped " + name)
					} else {
						fmt.Println("Failed to stop " + name)
						return false
					}
				}
			}
		}
	}
	return true
}

func (service Service) Unload() bool {
	if out, err := exec.Command("fleetctl", "list-units").Output(); err == nil {
		for _, str := range strings.Split(string(out), "\n") {
			deployName := strings.TrimSuffix(service.DeployName(), "@.service")
			// ignore other services
			if matched, _ := regexp.MatchString(service.Name, str); matched {
				// don't handle current deploy version
				if matched, _ := regexp.MatchString(deployName+"@", str); !matched {
					name := strings.Fields(str)[0]

					// stop service
					if _, err := exec.Command("fleetctl", "unload", name).Output(); err == nil {
						fmt.Println("Unloaded " + name)
					} else {
						fmt.Println("Failed to unload " + name)
						return false
					}
				}
			}
		}
	}
	return true
}

func (service Service) Submit() bool {
	timeoutSecs := strconv.Itoa(*timeout)
	timeoutStr := fmt.Sprintf("--request-timeout=%s", timeoutSecs)
	debugStr := "--debug"
	// cmdPrefix := fmt.Sprintf(" %s %s", timeoutStr, debugStr)

	fmt.Printf("Submitting new service file for: %s...\n", service.Name)
	name, _ := ioutil.TempDir(os.TempDir(), "FoxComm_"+*release)
	configFile := path.Join(name, service.DeployName())
	exec.Command("cp", "-f", service.Config, configFile).Run()
	exec.Command("sed", "-i", "-e", fmt.Sprintf(`s/\$VERSION\$/%s/g`, *deploy), configFile).Run()
	if *rackspace {
		if rackOut, rackErr := exec.Command("sed", "-i", "-e", fmt.Sprintf(`s/\$COREOS_PRIVATE_IPV4/%s/g`, "$RAX_SERVICENET_IPV4"), configFile).CombinedOutput(); rackErr != nil {
			fmt.Printf("Could not configure %s for RackSpace", configFile)
			if *debug {
				fmt.Printf("RackSpace Swap Output: %s", rackOut)
			}
		}

		if rackOut, rackErr := exec.Command("sed", "-i", "-e", fmt.Sprintf(`s/\$COREOS_PUBLIC_IPV4/%s/g`, "$RAX_PUBLICNET_IPV4"), configFile).CombinedOutput(); rackErr != nil {
			fmt.Printf("Could not configure %s for RackSpace", configFile)
			if *debug {
				fmt.Printf("RackSpace Swap Output: %s", rackOut)
			}
		}

	}

	if subOutput, err := exec.Command("fleetctl", timeoutStr, debugStr, "submit", configFile).CombinedOutput(); err == nil && !strings.Contains(string(subOutput), "differs from local unit file") && !strings.Contains(string(subOutput), "no need to recreate it") {
		fmt.Printf("FleetCtl submission output: %s\n", string(subOutput))
		for i := 1; i <= *instances; i++ {
			fmt.Printf("Starting instance %v for: %s...\n", i, service.DeployName())
			name := strings.Replace(service.DeployName(), "@.", fmt.Sprintf("@%v.", i), 1)
			cmd := exec.Command("fleetctl", timeoutStr, debugStr, "start", name)
			stdout, _ := cmd.StdoutPipe()
			if err := cmd.Start(); err == nil {
				in := bufio.NewScanner(stdout)
				for in.Scan() {
					log.Printf(in.Text())
				}
			} else {
				fmt.Println("Error starting 'fleetctl start' command : %s\n", err)
			}
		}
		return true
	} else {
		fmt.Printf("Failed to submit %s\n", service.DeployName())
		fmt.Printf("Service file might already exist.")
		if *debug {
			fmt.Printf("FleetCTL output : %s\n ", subOutput)
			fmt.Printf("CMD Error: %s\n", err)
		}
	}
	return false
}

func (service Service) DeployName() string {
	return strings.Replace(service.ConfigName(), "@", fmt.Sprintf("_%s@", *deploy), 1)
}

var release = flag.String("release", "", "Publish a new release")
var backup = flag.String("backup", path.Join(Root, "deployer", "fleet-backup"), "Set backup directory")
var target = flag.String("target", "all", "Set target to 'all' to target all services")
var instances = flag.Int("instances", 2, "Set instances count when deploy")
var deploy = flag.String("deploy", "", "Deploy a release")
var overwrite = flag.Bool("overwrite", false, "Overwrite existing service file.")
var timeout = flag.Int("timeout", 7, "Seconds before timing out on a request to Fleet via Tunnel")
var debug = flag.Bool("debug", false, "Show Debug Output")
var rackspace = flag.Bool("rackspace", false, "Swap IP values to be relevant for RackSpace")

func init() {
	flag.Parse()

	if _, err := os.Stat(path.Join(Root, "Procfile")); err != nil {
		fmt.Println("The script should be run from the root of Foxcomm")
		os.Exit(0)
	}

	if tunnel := os.Getenv("FLEETCTL_TUNNEL"); tunnel == "" {
		fmt.Println("Please set FLEETCTL_TUNNEL to a machine in the cluster and ensure you have access to that machine in order to continue.")
		os.Exit(0)
	}
}

func main() {
	// Ensure that user passes in release or deploy tag
	if *release == "" && *deploy == "" {
		fmt.Println("This script either publishes a release or deploys one to the fleet.   Please set the 'release' or 'deploy' parameters to continue.")
	}

	// Release a new version: create git tag
	if *release != "" {

		githubTag := fmt.Sprintf("%s_%s", *target, *release)

		if _, err := exec.Command("git", "tag", githubTag).Output(); err == nil {
			fmt.Printf("Pushing git tag %v to github...\n", githubTag)
			if out, err := exec.Command("git", "push", "origin", githubTag).Output(); err == nil {
				fmt.Println(string(out))
			} else {
				fmt.Printf("Error pushing release tag to GitHub: %s\n", err)
			}
		} else {
			fmt.Printf("Error tagging release in Github: %s\n", err)
		}

		foundReleaseTarget := false
		for _, service := range AllServices() {
			if service.Name == "solr-cluster" {
				continue
			}
			if *target == "all" || *target == service.Name {
				if !(service.Compile() && service.Release(*release)) {
					fmt.Println("Failed to publish a new release for " + service.Name)
					os.Exit(0)
				}
				foundReleaseTarget = true
			}
		}
		if !foundReleaseTarget {
			fmt.Printf("Target service with name %s was not found for release.\n", *target)
		}
	}

	// Deploy a version
	if *deploy != "" {
		hasFound := false
		anythingSubmitted := false
		for _, service := range AllServices() {
			//Skip the router if all targets are listed, because it's the only service that causes downtime.
			//TODO: Implement fancy router-deploy handling.
			if (*target == "all" && service.Name != "router") || *target == service.Name {
				hasFound = true

				if *overwrite {
					if service.Destroy() {
						fmt.Sprintf("Destroyed service file for overwriting: %s", service.Name)
					} else {
						fmt.Sprintf("Failed to destroy service file for overwriting: %s", service.Name)
					}
				}

				if service.Submit() {
					anythingSubmitted = true
					fmt.Printf("Successfully Submitted service: %s\n", service.Name)
					if service.Stop() {
						fmt.Printf("Successfully Stopped service: %s\n", service.Name)
						if service.Unload() {
							fmt.Printf("Successfully Unloaded service: %s\n", service.Name)
						} else {
							fmt.Printf("Failed to Unload Service: %s\n", service.Name)
						}
					} else {
						fmt.Printf("Failed to STOP Service: %s\n", service.Name)
					}

				} else {
					fmt.Printf("FAILED Submitted service: %s\n", service.Name)
				}
			}
		}
		if anythingSubmitted {
			fmt.Println("Deployment change, need to update router endpoints configuration")
			fmt.Println("")
			fmt.Println(`WITH_COREOS=1 ansible-playbook -i provisioner/hosts/staging/inventory provisioner/foxcomm.yml --vault-password-file provisioner/hosts/staging/vault_pass`)
		}
		if !hasFound {
			fmt.Printf("Target service with name %s was not found for deployment.\n", *target)
		}
	}
}
