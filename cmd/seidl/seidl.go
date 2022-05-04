package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

const VERSION = "0.1"

type Config struct {
	Region string
	Filter string // filtering string
}

var config Config

type Regions struct {
	Regions []Region `json:"regions"`
}
type Region struct {
	Name string `json:"name"`
}
type Environment struct {
	Name string `json:"name"`
}
type Image struct {
	Name        string `json:"name"`
	URN         string `json:"urn"`
	Id          string `json:"id"`
	State       string `json:"state"`
	ChangeInfo  string `json:"changeinfo"`
	PublishedOn string `json:"publishedon"`
	DeprecateOn string `json:"deprecatedon"`
	DeleteOn    string `json:"deletedon"`
	Environment string `json:"environment"`
	Region      string `json:"region"`
	Project     string `json:"project"`
}
type Images struct {
	Images []Image `json:"images"`
}

// Filter all images with the given filter. Returns the number of deleted entries
func (i *Images) filter(filter string) int {
	n := 0
	deleted := 0
	for _, val := range i.Images {
		if val.filter(filter) {
			i.Images[n] = val
			n++
		} else {
			deleted++
		}
	}
	i.Images = i.Images[:n]
	return deleted
}

func (i *Images) filterRegion(region string) {
	if region == "" {
		return
	}
	n := 0
	for _, val := range i.Images {
		if val.Region == region {
			i.Images[n] = val
			n++
		}
	}
	i.Images = i.Images[:n]
}

// Filter all images without the "suse-" prefix
func (i *Images) filterSUSE() {
	n := 0
	for _, val := range i.Images {
		if strings.HasPrefix(val.Name, "suse-") {
			i.Images[n] = val
			n++
		}
	}
	i.Images = i.Images[:n]
}

func (i *Image) filter(filter string) bool {
	if filter == "" {
		return true
	}
	// Multiple filters can be given by comma. All filters must be present
	filters := strings.Split(strings.ToLower(filter), ",")
	name := strings.ToLower(i.Name)
	for _, f := range filters {
		if !strings.Contains(name, strings.TrimSpace(f)) {
			return false
		}
	}
	return true
}

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] CSP...\n", os.Args[0])
	fmt.Println("CSP (cloud service providers)   gce|aws|azure")
	fmt.Println("OPTIONS:")
	fmt.Println("  -f, --filter FILTER           Filter results based on the given strings (comma-separated)")
	fmt.Println("  --region                      Set region (required for AWS)")
	fmt.Println("  --list-aws-regions            List AWS regions")
	fmt.Println("  --az-env                      Set environment (for Azure)")
	fmt.Println("  --list-az-envs                List possible Azure environments")
	fmt.Println("  --version                     Show program version")
	fmt.Println("")
	fmt.Println("Arguments are processed sequentially and a query is executed, once a CSP string is identified")
	fmt.Println("Consequently an argument folloing a CSP won't be considered in the query.")
	fmt.Printf("  right:  %s -f 'sles,15-sp2' gce\n", os.Args[0])
	fmt.Printf("  wrong:  %s gce -f 'sles,15-sp2'\n\n", os.Args[0])
}

func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return make([]byte, 0), err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// Fetch images from one of the following CSP: "microsoft", "google", "amazon"
func FetchImages(csp string) (Images, error) {
	var ret Images
	url := fmt.Sprintf("https://susepubliccloudinfo.suse.com/v1/%s/images.json", csp)
	body, err := fetch(url)
	if err != nil {
		return ret, err
	}
	if err := json.Unmarshal(body, &ret); err != nil {
		return ret, err
	}

	// Filter deprected images (in-place)
	n := 0
	for _, image := range ret.Images {
		if image.DeleteOn != "" {
			continue
		} else if image.DeprecateOn != "" {
			continue
		}
		ret.Images[n] = image
		n++
	}
	ret.Images = ret.Images[:n]
	// Sort by name
	sort.Slice(ret.Images, func(p, q int) bool { return ret.Images[p].Name < ret.Images[q].Name })
	return ret, err
}

func GetAzureEnvironments() ([]Environment, error) {
	// For now, just do it static
	envs := []string{"Blackforest", "Fairfax", "Mooncake", "PublicAzure"}

	ret := make([]Environment, 0)
	for _, name := range envs {
		env := Environment{Name: name}
		ret = append(ret, env)
	}
	return ret, nil
}

func GetAWSRegions() (Regions, error) {
	regions := Regions{}
	body, err := fetch("https://susepubliccloudinfo.suse.com/v1/amazon/regions.json")
	if err != nil {
		return regions, err
	}
	if err := json.Unmarshal(body, &regions); err != nil {
		return regions, err
	}
	return regions, nil
}

func isGCE(csp string) bool {
	return (csp == "g" || csp == "gce" || csp == "gcp" || csp == "google")
}
func isAzure(csp string) bool {
	return (csp == "m" || csp == "az" || csp == "azure" || csp == "microsoft")
}
func isAWS(csp string) bool {
	return (csp == "a" || csp == "aws" || csp == "ec2" || csp == "amazon")
}

func main() {
	// Parse program arguments, one after another
	args := os.Args[1:]
	if len(args) < 1 {
		usage()
		os.Exit(1)
	} else {

		// Convencience check: Ensure there are no dangling configuration parameters, as this is a common mistake
		dangling := ""
		for i := 0; i < len(args); i++ {
			arg := args[i]
			if arg == "" {
				continue
			}
			// Check for argument skips
			if arg == "-f" || arg == "--filter" || arg == "-r" || arg == "--region" {
				dangling = arg
				i++
			}
			// Check for CSP
			if isGCE(arg) || isAzure(arg) || isAWS(arg) {
				dangling = ""
			}
		}
		if dangling != "" {
			fmt.Fprintf(os.Stderr, "dangling argument: %s\nProgram arguments need to be BEFORE the CSP.\n", dangling)
			os.Exit(1)
		}

		for i := 0; i < len(args); i++ {
			arg := args[i]
			if arg == "" {
				continue
			}
			if arg[0] == '-' { // Configuration parameter
				if arg == "-h" || arg == "--help" {
					usage()
					return
				} else if arg == "--version" {
					fmt.Printf("seidl v%s -- https://github.com/grisu48/seidl/\n", VERSION)
					return
				} else if arg == "-f" || arg == "--filter" {
					i += 1
					config.Filter = args[i]
				} else if arg == "-r" || arg == "--region" {
					i += 1
					config.Region = args[i]
				} else if arg == "--list-az-envs" {
					envs, err := GetAzureEnvironments()
					if err != nil {
						fmt.Fprintf(os.Stderr, "error fetching environments: %s\n", err)
						os.Exit(1)
					}
					for _, env := range envs {
						fmt.Printf("%s\n", env.Name)
					}
				} else if arg == "--list-aws-regions" {
					envs, err := GetAWSRegions()
					if err != nil {
						fmt.Fprintf(os.Stderr, "error fetching regions: %s\n", err)
						os.Exit(1)
					}
					for _, env := range envs.Regions {
						fmt.Printf("%s\n", env.Name)
					}
				} else {
					fmt.Fprintf(os.Stderr, "invalid parameter: %s\n", arg)
					os.Exit(1)
				}
			} else {
				// Expect CSP argument
				if isGCE(arg) {
					images, err := FetchImages("google")
					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %s\n", err)
						os.Exit(1)
					}
					n := images.filter(config.Filter)
					if len(images.Images) == 0 {
						if n == 0 {
							fmt.Fprintf(os.Stderr, "no images found\n")
						} else {
							fmt.Fprintf(os.Stderr, "no images found (too restrictive filter?)\n")
						}
						os.Exit(1)
					}
					fmt.Printf("| %-58s | %-40s | %-20s |\n", "Name", "Project", "State")
					for _, image := range images.Images {
						fmt.Printf("%-60s | %-40s | %-20s\n", image.Name, image.Project, image.State)
					}
				} else if isAWS(arg) {
					images, err := FetchImages("amazon")
					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %s\n", err)
						os.Exit(1)
					}

					n := images.filter(config.Filter)
					if len(images.Images) == 0 {
						if n == 0 {
							fmt.Fprintf(os.Stderr, "no images found\n")
						} else {
							fmt.Fprintf(os.Stderr, "no images found (too restrictive filter?)\n")
						}
						os.Exit(1)
					}
					if config.Region == "" {
						fmt.Printf("| %-23s | %-60s | %-20s | %-20s |\n", "ID", "Name", "Region", "State")
						for _, image := range images.Images {
							fmt.Printf("%-25s | %-60s | %-20s | %-20s\n", image.Id, image.Name, image.Region, image.State)
						}
					} else {
						images.filterRegion(config.Region)
						fmt.Printf("| %-23s | %-60s | %-20s |\n", "ID", "Name", "State")
						for _, image := range images.Images {
							fmt.Printf("%-25s | %-60s | %-20s\n", image.Id, image.Name, image.State)
						}
					}
				} else if isAzure(arg) {
					images, err := FetchImages("microsoft")
					if err != nil {
						fmt.Fprintf(os.Stderr, "error: %s\n", err)
						os.Exit(1)
					}
					// Filter out weird entries
					images.filterSUSE()
					n := images.filter(config.Filter)
					if len(images.Images) == 0 {
						if n == 0 {
							fmt.Fprintf(os.Stderr, "no images found\n")
						} else {
							fmt.Fprintf(os.Stderr, "no images found (too restrictive filter?)\n")
						}
						os.Exit(1)
					}
					fmt.Printf("| %-58s | %-60s | %-20s\n", "URN", "Name", "State")
					for _, image := range images.Images {
						fmt.Printf("%-60s | %-60s | %-20s\n", image.URN, image.Name, image.State)
					}
				} else {
					fmt.Fprintf(os.Stderr, "error: invalid CSP\n")
					os.Exit(1)
				}
			}
		}
	}
}
