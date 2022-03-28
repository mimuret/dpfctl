package utils_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
	"github.com/mimuret/golang-iij-dpf/pkg/testtool"
)

var _ = Describe("client", func() {
	Context("NewClientDefalt", func() {
		var (
			cl  api.ClientInterface
			err error
		)
		When("context exists", func() {
			BeforeEach(func() {
				cl, err = utils.NewClientDefault(nil)
			})
			It("returns github.com/mimuret/golang-iij-dpf/pkg/api.Client", func() {
				Expect(cl).NotTo(BeNil())
			})
			It("not return err", func() {
				Expect(err).To(Succeed())
			})
		})
		When("context not eixst", func() {
			BeforeEach(func() {
				viper.Set("context", "NOTDOUND")
				cl, err = utils.NewClientDefault(nil)
			})
			It("returns err", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Context("Wait", func() {
		var (
			cl *testtool.TestClient
		)
		BeforeEach(func() {
			cl = testtool.NewTestClient("", "http://localhost", nil)
		})
		When("normal", func() {
			BeforeEach(func() {
				cl.ReadFunc = func(s api.Spec) (requestID string, err error) {
					job := s.(*core.Job)
					job.RequestID = "9BCFE2E9C10D4D9A8444CB0B48C72830"
					job.Status = core.JobStatusSuccessful
					return "ok", nil
				}
			})
			It("returns job", func() {
				job := &core.Job{
					RequestID: "9BCFE2E9C10D4D9A8444CB0B48C72830",
					Status:    core.JobStatusSuccessful,
				}
				Eventually(func() (*core.Job, error) {
					return utils.Wait(cl, "9BCFE2E9C10D4D9A8444CB0B48C72830", time.Second)
				}, time.Second).Should(Equal(job))
			})
		})
		When("fail job", func() {
			BeforeEach(func() {
				cl.ReadFunc = func(s api.Spec) (requestID string, err error) {
					job := s.(*core.Job)
					job.RequestID = "9BCFE2E9C10D4D9A8444CB0B48C72830"
					job.Status = core.JobStatusFailed
					return "ok", nil
				}
			})
			It("returns error", func() {
				Eventually(func() error {
					_, err := utils.Wait(cl, "9BCFE2E9C10D4D9A8444CB0B48C72830", time.Second)
					return err
				}, time.Second).Should(HaveOccurred())
			})
		})
		When("timeout", func() {
			BeforeEach(func() {
				cl.ReadFunc = func(s api.Spec) (requestID string, err error) {
					job := s.(*core.Job)
					job.RequestID = "9BCFE2E9C10D4D9A8444CB0B48C72830"
					job.Status = core.JobStatusRunning
					return "ok", nil
				}
				cl.WatchReadFunc = func(ctx context.Context, interval time.Duration, s api.Spec) error {
					time.Sleep(time.Second)
					return ctx.Err()
				}
			})
			It("returns error", func() {
				Eventually(func() error {
					_, err := utils.Wait(cl, "9BCFE2E9C10D4D9A8444CB0B48C72830", time.Second)
					return err
				}, time.Second).Should(HaveOccurred())
			})
		})
	})
})
