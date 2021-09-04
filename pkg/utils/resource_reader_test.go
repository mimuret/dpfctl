package utils_test

import (
	"bytes"
	_ "embed"
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"

	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/mimuret/golang-iij-dpf/pkg/testtool"
	_ "github.com/mimuret/golang-iij-dpf/pkg/testtool"
)

//go:embed testdata/single-doc.yaml
var singleYamlDoc []byte

//go:embed testdata/multi-doc.yaml
var multiYamlDoc []byte

//go:embed testdata/single-doc.json
var jsonDoc []byte

//go:embed testdata/bad.yaml
var badYamlDoc []byte

//go:embed testdata/bad-schema.yaml
var badSchemaDoc []byte

var _ = Describe("ResourceReader", func() {
	var (
		reader         *utils.ResourceReader
		err            error
		s1, s2, s3, s4 *testtool.TestSpec
		list           *testtool.TestSpecList
		docs           []json.RawMessage
		res            []apis.Spec
	)
	BeforeEach(func() {
		s1 = &testtool.TestSpec{Id: "id1", Name: "apple", Number: 10}
		s2 = &testtool.TestSpec{Id: "id2", Name: "orange", Number: 20}
		s3 = &testtool.TestSpec{Id: "id3", Name: "pen", Number: 40}
		s4 = &testtool.TestSpec{Id: "id10", Name: "green", Number: 999}
		list = &testtool.TestSpecList{Items: []testtool.TestSpec{*s1, *s2, *s3}}
		reader = utils.NewResourceReader(nil)
	})

	Context("ReadYamlDocuments", func() {
		When("single yaml document", func() {
			BeforeEach(func() {
				docs, err = reader.ReadYamlDocuments(bytes.NewBuffer(singleYamlDoc))
			})
			It("not return error", func() {
				Expect(err).To(Succeed())
			})
			It("len 1", func() {
				Expect(len(docs)).To(Equal(1))
			})
		})
		When("multi yaml document", func() {
			BeforeEach(func() {
				docs, err = reader.ReadYamlDocuments(bytes.NewBuffer(multiYamlDoc))
			})
			It("not return error", func() {
				Expect(err).To(Succeed())
			})
			It("len 2", func() {
				Expect(len(docs)).To(Equal(2))
			})
		})
		When("json document", func() {
			BeforeEach(func() {
				docs, err = reader.ReadYamlDocuments(bytes.NewBuffer(jsonDoc))
			})
			It("not return error", func() {
				Expect(err).To(Succeed())
			})
			It("len 1", func() {
				Expect(len(docs)).To(Equal(1))
			})
		})
		When("bad yaml document", func() {
			BeforeEach(func() {
				docs, err = reader.ReadYamlDocuments(bytes.NewBuffer(badYamlDoc))
			})
			It("not return error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("len 0", func() {
				Expect(len(docs)).To(Equal(0))
			})
		})
	})
	Context("ParseResouress", func() {
		When("single yaml document", func() {
			BeforeEach(func() {
				docs, err = reader.ReadYamlDocuments(bytes.NewBuffer(singleYamlDoc))
				Expect(err).To(Succeed())
				res, err = reader.ParseResouress(docs)
			})
			It("not return error", func() {
				Expect(err).To(Succeed())
			})
			It("len 1", func() {
				Expect(len(res)).To(Equal(1))
				Expect(res[0]).To(Equal(list))
			})
		})
		When("multi yaml document", func() {
			BeforeEach(func() {
				docs, err = reader.ReadYamlDocuments(bytes.NewBuffer(multiYamlDoc))
				Expect(err).To(Succeed())
				res, err = reader.ParseResouress(docs)
			})
			It("not return error", func() {
				Expect(err).To(Succeed())
			})
			It("docs len 2", func() {
				Expect(len(docs)).To(Equal(2))
				Expect(len(res)).To(Equal(2))
				Expect(res[0]).To(Equal(s4))
				Expect(res[1]).To(Equal(list))
			})
		})
		When("json document", func() {
			BeforeEach(func() {
				docs, err = reader.ReadYamlDocuments(bytes.NewBuffer(jsonDoc))
				Expect(err).To(Succeed())
				res, err = reader.ParseResouress(docs)
			})
			It("not return error", func() {
				Expect(err).To(Succeed())
			})
			It("len 1", func() {
				Expect(len(res)).To(Equal(1))
				Expect(res[0]).To(Equal(list))
			})
		})
		When("bad schema document", func() {
			BeforeEach(func() {
				docs, err = reader.ReadYamlDocuments(bytes.NewBuffer(badSchemaDoc))
				Expect(err).To(Succeed())
				res, err = reader.ParseResouress(docs)
			})
			It("not return error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("docs len is 1", func() {
				Expect(len(docs)).To(Equal(1))
			})
			It("len 0", func() {
				Expect(len(res)).To(Equal(0))
			})
		})
	})
	Context("GetResources", func() {
		BeforeEach(func() {
			fs := afero.NewMemMapFs()
			fs.MkdirAll("testdata", 0755)
			afero.WriteFile(fs, "testdata/bad-schema.yaml", badYamlDoc, 0644)
			afero.WriteFile(fs, "testdata/bad.yaml", badYamlDoc, 0644)
			afero.WriteFile(fs, "testdata/multi-doc.yaml", multiYamlDoc, 0644)
			afero.WriteFile(fs, "testdata/single-doc.yaml", singleYamlDoc, 0644)
			afero.WriteFile(fs, "testdata/single-doc.json", jsonDoc, 0644)
			reader = utils.NewResourceReader(fs)
		})
		When("file not exist", func() {
			BeforeEach(func() {
				res, err = reader.GetResources("testdata/not-exist.yaml")
			})
			It("not return error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
		When("single yaml document", func() {
			BeforeEach(func() {
				res, err = reader.GetResources("testdata/single-doc.yaml")
			})
			It("not return error", func() {
				Expect(err).To(Succeed())
			})
			It("len 1", func() {
				Expect(len(res)).To(Equal(1))
				Expect(res[0]).To(Equal(list))
			})
		})
		When("multi yaml document", func() {
			BeforeEach(func() {
				res, err = reader.GetResources("testdata/multi-doc.yaml")
			})
			It("not return error", func() {
				Expect(err).To(Succeed())
			})
			It("len 2", func() {
				Expect(len(res)).To(Equal(2))
				Expect(res[0]).To(Equal(s4))
				Expect(res[1]).To(Equal(list))
			})
		})
		When("json document", func() {
			BeforeEach(func() {
				res, err = reader.GetResources("testdata/single-doc.json")
			})
			It("not return error", func() {
				Expect(err).To(Succeed())
			})
			It("len 1", func() {
				Expect(len(res)).To(Equal(1))
				Expect(res[0]).To(Equal(list))
			})
		})
		When("bad yaml format document", func() {
			BeforeEach(func() {
				res, err = reader.GetResources("testdata/bad.yaml")
			})
			It("not return error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("len 0", func() {
				Expect(len(res)).To(Equal(0))
			})
		})
		When("bad schema document", func() {
			BeforeEach(func() {
				res, err = reader.GetResources("testdata/bad-schema.yaml")
			})
			It("not return error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("len 0", func() {
				Expect(len(res)).To(Equal(0))
			})
		})
	})
})
