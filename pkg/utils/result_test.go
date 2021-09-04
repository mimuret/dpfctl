package utils_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
	"github.com/mimuret/golang-iij-dpf/pkg/testtool"
)

var _ = Describe("CommandResults", func() {
	var (
		s      utils.CommandResults
		s1, s2 utils.CommandResult
	)
	BeforeEach(func() {
		s1 = utils.CommandResult{
			RequestId: "D44A2C1138B84CD0A20E858F6D57C17D",
			Err:       nil,
		}
		s2 = utils.CommandResult{
			RequestId: "3FC8732752844F029E4BC7DB4A60799C",
			Err:       nil,
		}
		s = utils.CommandResults{
			Items: []utils.CommandResult{s1, s2},
		}
	})
	Context("DeepCopyObject", func() {
		var (
			c *utils.CommandResults
		)
		BeforeEach(func() {
			c = s.DeepCopyObject().(*utils.CommandResults)
		})
		It("returns []interface{}", func() {
			Expect(c).To(Equal(&s))
		})
	})
	Context("Add", func() {
		var (
			c utils.CommandResults
		)
		When("no error", func() {
			BeforeEach(func() {
				c = utils.CommandResults{}
				c.Add("D44A2C1138B84CD0A20E858F6D57C17D", nil)
				c.Add("3FC8732752844F029E4BC7DB4A60799C", nil)
			})
			It("not set error", func() {
				Expect(c.Err).To(Succeed())
			})
			It("returns []interface{}", func() {
				Expect(c).To(Equal(s))
			})
		})
		When("with error", func() {
			BeforeEach(func() {
				c = utils.CommandResults{}
				c.Add("D44A2C1138B84CD0A20E858F6D57C17D", fmt.Errorf("error"))
				c.Add("3FC8732752844F029E4BC7DB4A60799C", nil)
			})
			It("not set error", func() {
				Expect(c.Err).To(HaveOccurred())
			})
			It("returns []interface{}", func() {
				Expect(c.Items[0].RequestId).To(Equal("D44A2C1138B84CD0A20E858F6D57C17D"))
				Expect(c.Items[0].Err).To(Equal(fmt.Errorf("error")))
				Expect(c.Items[1].RequestId).To(Equal("3FC8732752844F029E4BC7DB4A60799C"))
				Expect(c.Items[1].Err).To(BeNil())
			})
		})
	})
	Context("WaitJob", func() {
		var (
			c  utils.CommandResults
			cl *testtool.TestClient
			v  *viper.Viper
		)
		BeforeEach(func() {
			cl = testtool.NewTestClient("", "http://localhost", nil)
			c = utils.CommandResults{}
			c.Add("D44A2C1138B84CD0A20E858F6D57C17D", nil)
			c.Add("3FC8732752844F029E4BC7DB4A60799C", nil)
			v = viper.New()
			v.SetDefault("wait-timeout", time.Second)
		})
		When("wait=false", func() {
			BeforeEach(func() {
				v.SetDefault("wait", false)
				c.WaitJob(cl, v)
			})
			It("no operation anything", func() {
				Expect(c).To(Equal(s))
			})
		})
		When("wait=true", func() {
			BeforeEach(func() {
				v.SetDefault("wait", true)
			})
			When("all job succeed", func() {
				BeforeEach(func() {
					cl.ReadFunc = func(s api.Spec) (requestId string, err error) {
						job := s.(*core.Job)
						job.Status = core.JobStatusSuccessful
						return job.RequestId, nil
					}
					c.WaitJob(cl, v)
				})
				It("not set err", func() {
					Expect(c.Err).To(Succeed())
				})
				It("set job with JobStatusSuccessful", func() {
					Expect(c.Items[0].Job).NotTo(BeNil())
					Expect(c.Items[1].Job).NotTo(BeNil())
					Expect(c.Items[0].Job.Status).To(Equal(core.JobStatusSuccessful))
					Expect(c.Items[1].Job.Status).To(Equal(core.JobStatusSuccessful))
				})
			})
			When("job failed", func() {
				BeforeEach(func() {
					cl.ReadFunc = func(s api.Spec) (requestId string, err error) {
						job := s.(*core.Job)
						job.Status = core.JobStatusFailed
						return job.RequestId, nil
					}
					c.WaitJob(cl, v)
				})
				It("not set err", func() {
					Expect(c.Err).To(HaveOccurred())
				})
				It("set job with JobStatusFailed", func() {
					Expect(c.Items[0].Job).NotTo(BeNil())
					Expect(c.Items[1].Job).NotTo(BeNil())
					Expect(c.Items[0].Job.Status).To(Equal(core.JobStatusFailed))
					Expect(c.Items[1].Job.Status).To(Equal(core.JobStatusFailed))
				})
			})
		})
	})
	Context("not implement funcs", func() {
		It("return empty", func() {
			Expect(s.GetName()).To(BeEmpty())
			Expect(s.GetGroup()).To(BeEmpty())
			Expect(s.GetPathMethod(api.ActionApply)).To(BeEmpty())
		})
	})
})
