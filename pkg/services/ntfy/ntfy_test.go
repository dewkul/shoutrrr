package ntfy_test

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"testing"

	"github.com/containrrr/shoutrrr/internal/testutils"
	. "github.com/containrrr/shoutrrr/pkg/services/ntfy"
	"github.com/containrrr/shoutrrr/pkg/types"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNtfy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Ntfy Suite")
}

var (
	service *Service
	logger  *log.Logger
	err     error
)

var _ = Describe("Ntfy service", func() {
	BeforeSuite(func() {
		service = &Service{}
		logger = log.New(GinkgoWriter, "Test", log.LstdFlags)
	})

	Describe("Generating Ntfy config from URL", func() {
		When("using default config value", func() {
			ntfyUrl, _ := url.Parse("ntfy://@/topicName")
			cfg, cfgErr := CreateConfigFromURL(ntfyUrl)
			It("should use https scheme", func() {
				Expect(cfgErr).NotTo(HaveOccurred())
				Expect(cfg.DisableTLS).To(Equal(false))
			})
			It("should use default hostname", func() {
				Expect(cfg.Host).To(Equal("ntfy.sh"))
			})
			It("should set topic", func() {
				Expect(cfg.Topic).To(Equal("topicName"))
			})
		})
		When("providing custom config values", func() {
			host := "ntfy.domain.tld"
			topic := "myTopic"
			attachment := ""
			clickUrl := "domain.tld"
			delay := "30m"
			email := "jone@doe.me"
			isTlsDisable := "Yes"
			priority := uint8(2)
			tags := "a,b,c"

			notificationUrl := fmt.Sprintf(
				"ntfy://@%v/%v?attach=%v&click=%v&delay=%v&disabletls=%v&email=%v&priority=%v&tags=%v",
				host, topic, attachment, clickUrl, delay, isTlsDisable, email, priority, tags,
			)
			ntfyUrl, _ := url.Parse(notificationUrl)
			cfg, cfgErr := CreateConfigFromURL(ntfyUrl)
			It("should not produced an error", func() {
				Expect(cfgErr).NotTo(HaveOccurred())
			})
			It("should use http scheme", func() {
				Expect(cfg.DisableTLS).To(Equal(isTlsDisable == "Yes"))
			})
			It("should use a given host", func() {
				Expect(cfg.Host).To(Equal(host))
			})
			It("should use a given topic", func() {
				Expect(cfg.Topic).To(Equal(topic))
			})
			It("should use a given attachment url", func() {
				Expect(cfg.Attachment).To(Equal(attachment))
			})
			It("should use a given click url", func() {
				Expect(cfg.ClickURL).To(Equal(clickUrl))
			})
			It("should use a given delay duration", func() {
				Expect(cfg.Delay).To(Equal(delay))
			})
			It("should use a given email", func() {
				Expect(cfg.Attachment).To(Equal(attachment))
			})
			It("should use a given priority", func() {
				Expect(cfg.Priority).To(Equal(priority))
			})
			It("should use given tags", func() {
				Expect(cfg.Tags).To(Equal(strings.Split(tags, ",")))
			})
		})
		When("no topic is provided", func() {
			It("should report an error", func() {
				ntfyUrl, _ := url.Parse("ntfy://@ntfy.sh")
				_, cfgErr := CreateConfigFromURL(ntfyUrl)
				Expect(cfgErr).To(HaveOccurred())
			})
		})
		When("publishing to a protected topic", func() {
			It("should set username and password", func() {
				ntfyUrl, _ := url.Parse("ntfy://username:password@ntfy.sh/topic")
				cfg, cfgErr := CreateConfigFromURL(ntfyUrl)
				Expect(cfgErr).NotTo(HaveOccurred())
				Expect(cfg.Username).To(Equal("username"))
				Expect(cfg.Password).To(Equal("password"))
			})
			It("should ignore username in case of empty password", func() {
				ntfyUrl, _ := url.Parse("ntfy://username@ntfy.sh/topic")
				cfg, cfgErr := CreateConfigFromURL(ntfyUrl)
				Expect(cfgErr).NotTo(HaveOccurred())
				Expect(cfg.Username).To(BeEmpty())
				Expect(cfg.Password).To(BeEmpty())
			})
			It("should ignore username if password is not set", func() {
				ntfyUrl, _ := url.Parse("ntfy://username:@ntfy.sh/topic")
				cfg, cfgErr := CreateConfigFromURL(ntfyUrl)
				Expect(cfgErr).NotTo(HaveOccurred())
				Expect(cfg.Username).To(BeEmpty())
				Expect(cfg.Password).To(BeEmpty())
			})
			It("should ignore password if username is not set", func() {
				ntfyUrl, _ := url.Parse("ntfy://:password@ntfy.sh/topic")
				cfg, cfgErr := CreateConfigFromURL(ntfyUrl)
				Expect(cfgErr).NotTo(HaveOccurred())
				Expect(cfg.Username).To(BeEmpty())
				Expect(cfg.Password).To(BeEmpty())
			})
		})
	})

	Describe("sending the payload", func() {
		When("sending to ntfy server", func() {
			AfterEach(func() {
				httpmock.DeactivateAndReset()
			})
			It("should not throw an error if the server accepts payload", func() {
				ntfyURL, _ := url.Parse("ntfy://@ntfy.test/topic")
				err = service.Initialize(ntfyURL, logger)
				Expect(err).NotTo(HaveOccurred())

				httpmock.ActivateNonDefault(service.GetHTTPClient())
				targetUrl := "https://ntfy.test"
				httpmock.RegisterResponder("POST", targetUrl, testutils.JSONRespondMust(200, MessageResponse{}))

				err = service.Send("Message", &types.Params{"title": ""})
				Expect(err).NotTo(HaveOccurred())
			})
			It("should not panic if an error occurs when sending the payload", func() {
				ntfyURL, _ := url.Parse("ntfy://@ntfy.sh/topic")
				err = service.Initialize(ntfyURL, logger)
				Expect(err).NotTo(HaveOccurred())

				httpmock.ActivateNonDefault(service.GetHTTPClient())
				targetUrl := "https://ntfy.sh"
				httpmock.RegisterResponder("POST", targetUrl, testutils.JSONRespondMust(401, ErrorResponse{
					Name:        "err",
					Code:        401,
					Description: "Not authorized",
				}))

				err = service.Send("Message", &types.Params{"title": ""})
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
