package manifest_test

import (
	"encoding/json"

	"code.cloudfoundry.org/cf-operator/pkg/bosh/manifest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const boshDNSAddOn = `
{
  "jobs": [
    {
      "name": "bosh-dns-aliases",
      "properties": {
        "aliases": [
          {
            "domain": "_.cell.service.cf.internal",
            "targets": [
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "diego-cell",
                "network": "default",
                "query": "_"
              },
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "windows2012R2-cell",
                "network": "default",
                "query": "_"
              },
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "windows2016-cell",
                "network": "default",
                "query": "_"
              },
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "windows1803-cell",
                "network": "default",
                "query": "_"
              },
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "windows2019-cell",
                "network": "default",
                "query": "_"
              },
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "isolated-diego-cell",
                "network": "default",
                "query": "_"
              }
            ]
          },
          {
            "domain": "auctioneer.service.cf.internal",
            "targets": [
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "scheduler",
                "network": "default",
                "query": "q-s4"
              }
            ]
          },
           {
            "domain": "bbs1.service.cf.internal",
            "targets": [
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "diego-api",
                "network": "default",
                "query": "q-s4"
              }
            ]
          },
         {
            "domain": "bbs.service.cf.internal",
            "targets": [
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "diego-api",
                "network": "default",
                "query": "q-s4"
              }
            ]
          },
          {
            "domain": "bits.service.cf.internal",
            "targets": [
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "bits",
                "network": "default",
                "query": "*"
              }
            ]
          },
          {
            "domain": "uaa.service.cf.internal",
            "targets": [
              {
                "deployment": "cf",
                "domain": "bosh",
                "instance_group": "uaa",
                "network": "default",
                "query": "*"
              }
            ]
          }
        ]
      },
      "release": "bosh-dns-aliases"
    }
  ],
  "name": "bosh-dns-aliases"
}
`

const logAPIJob = `
{
  "consumes": {
    "doppler": {
      "from": "doppler"
    }
  },
  "name": "loggregator_trafficcontroller",
  "properties": {
    "cc": {
      "internal_service_hostname": "cloud-controller-ng.service.cf.internal",
      "mutual_tls": {
        "ca_cert": "((service_cf_internal_ca.certificate))"
      },
      "tls_port": 9023
    },
    "loggregator": {
      "outgoing_cert": "((loggregator_trafficcontroller_tls.certificate))",
      "outgoing_key": "((loggregator_trafficcontroller_tls.private_key))",
      "tls": {
        "ca_cert": "((loggregator_ca.certificate))",
        "cc_trafficcontroller": {
          "cert": "((loggregator_tls_cc_tc.certificate))",
          "key": "((loggregator_tls_cc_tc.private_key))"
        },
        "trafficcontroller": {
          "cert": "((loggregator_tls_tc.certificate))",
          "key": "((loggregator_tls_tc.private_key))"
        }
      },
      "uaa": {
        "client_secret": "((uaa_clients_doppler_secret))"
      }
    },
    "ssl": {
      "skip_cert_verify": true
    },
    "system_domain": "((system_domain))",
    "uaa": {
      "ca_cert": "((uaa_ca.certificate))",
      "internal_url": "https://uaa.service.cf.internal:8443",
      "other_urls" : [
         "https://uaa.service.cf.internal:8443",
         { "host" : "https://uaa.service.cf.internal:8443" }
      ]
    }
  },
  "release": "loggregator"
}
`

func loadAddOn(data string) *manifest.AddOn {
	var addOn manifest.AddOn
	err := json.Unmarshal([]byte(data), &addOn)
	if err != nil {
		panic("Loading yaml failed")
	}
	return &addOn
}

var _ = Describe("kube converter", func() {

	Context("bosh-dns", func() {
		It("loads dns from addons correct", func() {
			dns, err := manifest.NewDomainNameService("default", loadAddOn(boshDNSAddOn))
			Expect(err).NotTo(HaveOccurred())
			Expect(dns.Aliases).To(HaveLen(6))
			Expect(dns.Aliases[5].Domain).To(Equal("uaa.service.cf.internal"))
		})
		It("replaces dns entries in jobs", func() {
			dns, err := manifest.NewDomainNameService("namespace1", loadAddOn(boshDNSAddOn))
			Expect(err).NotTo(HaveOccurred())
			var job manifest.Job
			err = json.Unmarshal([]byte(logAPIJob), &job)
			Expect(err).NotTo(HaveOccurred())
			dns.ReplaceProperties(job.Properties.Properties)
			value, ok := job.Property("uaa.internal_url")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("https://uaa.namespace1.svc.cluster.local:8443"))
			value, ok = job.Property("cc.internal_service_hostname")
			Expect(ok).To(BeTrue())
			Expect(value).To(ContainSubstring("cloud-controller-ng.service.cf.internal")) // No replacement for unknown dns entry
			value, ok = job.Property("uaa.other_urls")
			Expect(ok).To(BeTrue())
			value = value.([]interface{})[0].(string)
			Expect(value).To(Equal("https://uaa.namespace1.svc.cluster.local:8443"))
		})
		It("returns the correct service names", func() {
			dns, err := manifest.NewDomainNameService("default", loadAddOn(boshDNSAddOn))
			Expect(err).NotTo(HaveOccurred())
			diegoAPI := dns.FindServiceNames("diego-api", "cf")
			Expect(diegoAPI).To(ConsistOf("bbs", "bbs1"))
			uaa := dns.FindServiceNames("uaa", "cf")
			Expect(uaa).To(ConsistOf("uaa"))
			invalid := dns.FindServiceNames("invalid", "cf")
			Expect(invalid).To(ConsistOf())

		})
	})
})
