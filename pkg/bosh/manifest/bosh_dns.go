package manifest

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Target of domain alias
type Target struct {
	Query         string `json:"query"`
	InstanceGroup string `json:"instance_group"`
	Deployment    string `json:"deployment"`
	Network       string `json:"network"`
	Domain        string `json:"domain"`
}

// Alias of domain alias
type Alias struct {
	Domain  string   `json:"domain"`
	Targets []Target `json:"targets"`
}

// DomainNameService is used to emulate Bosh DNS
type DomainNameService struct {
	Namespace string
	Aliases   []Alias
}

// BoshDNSAddOnName name of bosh dns add on
const BoshDNSAddOnName = "bosh-dns-aliases"

var domainRegexp = regexp.MustCompile("\\.service\\.cf\\.internal")

// NewDomainNameService create a new DomainNameService
func NewDomainNameService(namespace string, addOn *AddOn) (*DomainNameService, error) {
	dns := DomainNameService{Namespace: namespace}
	for _, job := range addOn.Jobs {
		aliases := job.Properties.Properties["aliases"]
		if aliases != nil {
			aliasesBytes, err := json.Marshal(aliases)
			if err != nil {
				return nil, errors.Wrapf(err, "Loading aliases from manifest")
			}
			var a = make([]Alias, 0)
			err = json.Unmarshal(aliasesBytes, &a)
			if err != nil {
				return nil, errors.Wrapf(err, "Loading aliases from manifest")
			}
			dns.Aliases = append(dns.Aliases, a...)
		}
	}
	return &dns, nil
}

// ReplaceProperties replaces DNS entries job properties
func (dns *DomainNameService) ReplaceProperties(properties map[string]interface{}) error {
	replace := make(map[string][]byte)
	kubeBaseDomain := "." + dns.Namespace + ".svc.cluster.local"
	expr := "("
	for _, alias := range dns.Aliases {
		replace[alias.Domain] = []byte(domainRegexp.ReplaceAllString(alias.Domain, kubeBaseDomain))
		if len(expr) > 1 {
			expr = expr + "|"
		}
		expr = expr + regexp.QuoteMeta(alias.Domain)
	}
	expr = expr + ")"
	re := regexp.MustCompile(expr)

	bytes, err := json.Marshal(properties)
	if err != nil {
		return errors.Wrapf(err, "Marshal job properties")
	}
	bytes = re.ReplaceAllFunc(bytes, func(search []byte) []byte {
		return replace[string(search)]
	})
	err = json.Unmarshal(bytes, &properties)
	if err != nil {
		return errors.Wrapf(err, "Unmarshal job properties")
	}
	return nil
}

// FindServiceNames determines how a service should be named in accordance with the 'bosh-dns'-addon
func (dns *DomainNameService) FindServiceNames(instanceGroupName string, deploymentName string) []string {
	result := make([]string, 0)
	for _, alias := range dns.Aliases {
		for _, target := range alias.Targets {
			if target.InstanceGroup == instanceGroupName /* && target.Deployment == deploymentName */ {
				parts := strings.Split(alias.Domain, ".")
				if len(parts) == 0 {
					panic("Should never happen")
				}
				result = append(result, strings.Split(alias.Domain, ".")[0])
			}
		}
	}
	return result
}
