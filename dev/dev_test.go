package dev_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Context("When the gospell image is available", func() {
		It("can create a service", func() {
			_, err := s.Up(ctx)
			Expect(err).NotTo(HaveOccurred())
			// wait for resources to be created
			inProgress := true
			for inProgress {
				time.Sleep(time.Second * 2)
				upStatus, err := s.Info(ctx)
				Expect(err).NotTo(HaveOccurred())
				inProgress = upStatus.UpdateInProgress
			}
			outputMap, err := s.Outputs(ctx)
			Expect(err).NotTo(HaveOccurred())
			serviceURL := outputMap["srvurl"].Value.(string)
			resp, err = http.Get(serviceURL + "/v1/spell?input=hi")
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
			Expect(err).NotTo(HaveOccurred())
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			var actual []string
			json.Unmarshal(body, &actual)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).Should(Equal(expected))
		})
	})
})
