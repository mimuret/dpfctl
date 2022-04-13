package params_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
)

func TestAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "params Suite")
}

var _ = Describe("GroupMap", func() {
	var (
		gm              *params.GroupMap
		bookSet, penSet *params.APISet
		list            params.APISetSlice
		invalidSet      *params.APISet
	)
	BeforeEach(func() {
		gm = params.NewGroupMap()
		bookSet = &params.APISet{
			Name:        "book",
			Description: "hoge",
			Action: map[api.Action]params.API{
				api.ActionRead: {
					nil,
					params.Params{
						{Name: "id", Type: params.ParamTypeInt64, Required: true},
					},
					nil, "",
				},
				api.ActionList: {
					nil,
					params.Params{},
					nil, "",
				},
			},
		}
		penSet = &params.APISet{
			Name:        "pen",
			Description: "hoge",
			Action: map[api.Action]params.API{
				api.ActionRead: {
					nil,
					params.Params{
						{Name: "id", Type: params.ParamTypeArrayString, Required: true},
					},
					nil, "",
				},
			},
		}
		invalidSet = &params.APISet{
			Name:        "invalid",
			Description: "invalid param set",
			Action: map[api.Action]params.API{
				api.ActionRead: {
					nil,
					params.Params{
						params.Param{Name: "IDS", Type: params.ParamTypeArrayInt64, Required: true},
						params.Param{Name: "Value", Type: params.ParamTypeString, Required: false},
					},
					nil, "",
				},
			},
		}
		list = params.APISetSlice{penSet, bookSet}
		gm.SetGroup("test", list)
	})
	Context("NewGroupMap", func() {
		It("returns *GroupMap", func() {
			Expect(params.NewGroupMap()).ShouldNot(BeNil())
		})
	})
	Context("IterateGroup", func() {
		var (
			names []string
		)
		BeforeEach(func() {
			gm.IterateGroup(func(name string, list params.APISetSlice) {
				names = append(names, name)
			})
		})
		It("can iterate apiSet", func() {
			Expect(names).Should(Equal([]string{"test"}))
		})
	})
	Context("UpdateCmdMap", func() {
		When("duplicate commnand name", func() {
			It("raise panic", func() {
				Expect(func() { gm.SetGroup("newtest", list) }).Should(Panic())
			})
		})
	})
	Context("ValidateParams", func() {
		When("params are valid", func() {
			It("not raise panic", func() {
				Expect(func() { gm.ValidateParams() }).ShouldNot(Panic())
			})
		})
		When("any of params is not valid", func() {
			It(" raise panic", func() {
				list3 := params.APISetSlice{invalidSet}
				gm.SetGroup("invalid", list3)
				Expect(func() { gm.ValidateParams() }).Should(Panic())
			})
		})
	})
	Context("GetAPISetfromCmdName", func() {
		var (
			apiSet *params.APISet
		)
		When("exist", func() {
			BeforeEach(func() {
				apiSet = gm.GetAPISetfromCmdName("book")
			})
			It("return *params.APISet", func() {
				Expect(apiSet).To(Equal(bookSet))
			})
		})
		When("not exist", func() {
			BeforeEach(func() {
				apiSet = gm.GetAPISetfromCmdName("book2")
			})
			It("return nil", func() {
				Expect(apiSet).To(BeNil())
			})
		})
	})
	Context("SetGroup", func() {
		var (
			testList params.APISetSlice
			list2    params.APISetSlice
		)
		BeforeEach(func() {
			list2 = params.APISetSlice{penSet}
		})
		It("can iterate apiSet", func() {
			testList = gm.GetAPISlice("test")
			Expect(testList).Should(Equal(list))
		})
		When("already set group", func() {
			BeforeEach(func() {
				gm.SetGroup("test", list2)
			})
			It("can replace group", func() {
				testList = gm.GetAPISlice("test")
				Expect(testList).Should(Equal(list2))
			})
		})
		When("duplicate commnand name", func() {
			It("raise panic", func() {
				Expect(func() { gm.SetGroup("test2", list2) }).Should(Panic())
			})
		})
	})
	Context("GetAPISlice", func() {
		var (
			setSlice params.APISetSlice
		)
		When("exist", func() {
			BeforeEach(func() {
				setSlice = gm.GetAPISlice("test")
			})
			It("return *params.APISet", func() {
				Expect(setSlice).To(Equal(list))
			})
		})
		When("not exist", func() {
			BeforeEach(func() {
				setSlice = gm.GetAPISlice("hogehoge")
			})
			It("return nil", func() {
				Expect(setSlice).To(BeNil())
			})
		})
	})
	Context("GetValidArgs", func() {
		var (
			vars []string
		)
		When("ActionRead", func() {
			BeforeEach(func() {
				vars = gm.GetValidArgs(api.ActionRead)
			})
			It("returns sorted readable command names", func() {
				Expect(vars).To(Equal([]string{"book", "pen"}))
			})
		})
		When("ActionList", func() {
			BeforeEach(func() {
				vars = gm.GetValidArgs(api.ActionList)
			})
			It("returns sorted readable command names", func() {
				Expect(vars).To(Equal([]string{"book"}))
			})
		})
	})
})

var _ = Describe("APISetSlice", func() {
	var (
		setSlice params.APISetSlice
	)
	BeforeEach(func() {
		bookSet := &params.APISet{
			Name:        "book",
			Description: "hoge",
		}
		penSet := &params.APISet{
			Name:        "pen",
			Description: "hoge",
		}
		setSlice = params.APISetSlice{penSet, bookSet}
	})
	Context("IterateAPISet", func() {
		var (
			names []string
		)
		BeforeEach(func() {
			setSlice.IterateAPISet(func(apiSet *params.APISet) {
				names = append(names, apiSet.Name)
			})
		})
		It("can iterate apiSet", func() {
			Expect(names).Should(Equal([]string{"book", "pen"}))
		})
	})
})

var _ = Describe("Params", func() {
	var (
		err        error
		paramSlice params.Params
	)
	Context("Validate", func() {
		When("Any of Param is empty", func() {
			It("returns nil", func() {
				paramSlice = params.Params{}
				err = paramSlice.Validate()
				Expect(err).To(Succeed())
			})
		})
		When("All of Param is valid", func() {
			It("returns nil", func() {
				paramSlice = params.Params{
					params.Param{Name: "id", Type: params.ParamTypeString, Required: true},
				}
				err = paramSlice.Validate()
				Expect(err).To(Succeed())
			})
		})
		When("Any of Param is not valid", func() {
			It("returns error", func() {
				paramSlice = params.Params{
					params.Param{},
				}
				err = paramSlice.Validate()
				Expect(err).To(HaveOccurred())
			})
		})
		When("variable length arguments it not last Param", func() {
			It("returns error", func() {
				paramSlice = params.Params{
					params.Param{Name: "ids", Type: params.ParamTypeArrayInt64, Required: true},
					params.Param{Name: "value", Type: params.ParamTypeString},
				}
				err = paramSlice.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("param\\[0\\]: ParamTypeArrayInt64 must be last item"))
				paramSlice = params.Params{
					params.Param{Name: "ids", Type: params.ParamTypeArrayString, Required: true},
					params.Param{Name: "value", Type: params.ParamTypeString},
				}
				err = paramSlice.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("param\\[0\\]: ParamTypeArrayString must be last item"))
			})
		})
		When("Required is false. but it not last Param", func() {
			It("returns error", func() {
				paramSlice = params.Params{
					params.Param{Name: "id", Type: params.ParamTypeInt64, Required: false},
					params.Param{Name: "value", Type: params.ParamTypeString, Required: true},
				}
				err = paramSlice.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("param\\[0\\]: Required false must be last item"))
			})
		})
	})
	Context("MakeArgs", func() {
		var (
			paramSlice params.Params
			setArgs    []interface{}
			err        error
		)
		When("last Param is not variable length arguments. it is Required.", func() {
			BeforeEach(func() {
				paramSlice = params.Params{
					params.Param{Name: "id", Type: params.ParamTypeInt64, Required: true},
					params.Param{Name: "value", Type: params.ParamTypeString, Required: true},
				}
			})
			When("enough cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"1", "book"})
				})
				It("is successful", func() {
					Expect(err).To(Succeed())
					Expect(setArgs).To(Equal([]interface{}{int64(1), "book"}))
				})
			})
			When("missing cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"1"})
				})
				It("returns err", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not enough parameters"))
				})
			})
			When("can't parse args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"one", "book"})
				})
				It("returns err", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("param\\[0\\]: int64 parse error"))
				})
			})
		})
		When("last Param is not variable length arguments. it is not Required", func() {
			BeforeEach(func() {
				paramSlice = params.Params{
					params.Param{Name: "id", Type: params.ParamTypeInt64, Required: true},
					params.Param{Name: "value", Type: params.ParamTypeString, Required: false},
				}
			})
			When("enough cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"1", "book"})
				})
				It("is successful", func() {
					Expect(err).To(Succeed())
					Expect(setArgs).To(Equal([]interface{}{int64(1), "book"}))
				})
			})
			When("complete cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"1", "book"})
				})
				It("is successful", func() {
					Expect(err).To(Succeed())
					Expect(setArgs).To(Equal([]interface{}{int64(1), "book"}))
				})
			})
			When("enough cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"1"})
				})
				It("is successful", func() {
					Expect(err).To(Succeed())
					Expect(setArgs).To(Equal([]interface{}{int64(1)}))
				})
			})
			When("missing cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{})
				})
				It("returns err", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not enough parameters"))
				})
			})
			When("can't parse args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"one"})
				})
				It("returns err", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("param\\[0\\]: int64 parse error"))
				})
			})
		})
		When("last Param is variable length arguments. it is Required", func() {
			BeforeEach(func() {
				paramSlice = params.Params{
					params.Param{Name: "id", Type: params.ParamTypeInt64, Required: true},
					params.Param{Name: "value", Type: params.ParamTypeArrayString, Required: true},
				}
			})
			When("complete cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"1", "book", "pen"})
				})
				It("is successful", func() {
					Expect(err).To(Succeed())
					Expect(setArgs).To(Equal([]interface{}{int64(1), "book", "pen"}))
				})
			})
			When("enough cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"1", "book"})
				})
				It("is successful", func() {
					Expect(err).To(Succeed())
					Expect(setArgs).To(Equal([]interface{}{int64(1), "book"}))
				})
			})
			When("missing cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"1"})
				})
				It("returns err", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not enough parameters"))
				})
			})
			When("can't parse args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"one"})
				})
				It("returns err", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("param\\[0\\]: int64 parse error"))
				})
			})
		})
		When("last Param is variable length arguments. it is not Required", func() {
			BeforeEach(func() {
				paramSlice = params.Params{
					params.Param{Name: "id", Type: params.ParamTypeInt64, Required: true},
					params.Param{Name: "value", Type: params.ParamTypeArrayString, Required: false},
				}
			})
			When("complete cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"1", "book", "pen"})
				})
				It("is successful", func() {
					Expect(err).To(Succeed())
					Expect(setArgs).To(Equal([]interface{}{int64(1), "book", "pen"}))
				})
			})
			When("enough cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"1"})
				})
				It("is successful", func() {
					Expect(err).To(Succeed())
					Expect(setArgs).To(Equal([]interface{}{int64(1)}))
				})
			})
			When("missing cmd args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{})
				})
				It("returns err", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("not enough parameters"))
				})
			})
			When("can't parse args", func() {
				BeforeEach(func() {
					setArgs, err = paramSlice.MakeArgs([]string{"one"})
				})
				It("returns err", func() {
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("param\\[0\\]: int64 parse error"))
				})
			})
		})
	})
})

var _ = Describe("Param", func() {
	var (
		err error
	)
	Context("Validate", func() {
		When("name is empty", func() {
			BeforeEach(func() {
				err = params.Param{Type: params.ParamTypeInt64}.Validate()
			})
			It("returns err", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("name must not be empty"))
			})
		})
		When("name include special character", func() {
			BeforeEach(func() {
				err = params.Param{Name: "[ID]", Type: params.ParamTypeInt64}.Validate()
			})
			It("returns err", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("name must not include special character"))
			})
		})
		When("type is invalid", func() {
			BeforeEach(func() {
				err = params.Param{Name: "id", Type: params.ParamType(0)}.Validate()
			})
			It("returns err", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("type error"))
			})
		})
	})
	Context("Parse", func() {
		var (
			v   interface{}
			err error
		)
		When("type is ParamTypeArrayInt64, input is interger string", func() {
			BeforeEach(func() {
				p := params.Param{Type: params.ParamTypeArrayInt64}
				v, err = p.Parse("1")
			})
			It("returns nil", func() {
				Expect(err).To(Succeed())
				Expect(v).To(Equal(interface{}(int64(1))))
			})
		})
		When("type is ParamTypeInt64, input is interger string", func() {
			BeforeEach(func() {
				p := params.Param{Type: params.ParamTypeInt64}
				v, err = p.Parse("1")
			})
			It("returns nil", func() {
				Expect(err).To(Succeed())
				Expect(v).To(Equal(interface{}(int64(1))))
			})
		})
		When("type is ParamTypeArrayInt64, input is not interger string", func() {
			BeforeEach(func() {
				p := params.Param{Type: params.ParamTypeArrayInt64}
				v, err = p.Parse("one")
			})
			It("returns nil", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("int64 parse error"))
			})
		})
		When("type is ParamTypeInt64, input is not interger string", func() {
			BeforeEach(func() {
				p := params.Param{Type: params.ParamTypeInt64}
				v, err = p.Parse("one")
			})
			It("returns nil", func() {
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("int64 parse error"))
			})
		})
		When("type is ParamTypeArrayString, ParamTypeString", func() {
			It("returns string", func() {
				p := params.Param{Type: params.ParamTypeString}
				v, err = p.Parse("1")
				Expect(err).To(Succeed())
				Expect(v).To(Equal(interface{}("1")))
				p = params.Param{Type: params.ParamTypeArrayString}
				v, err = p.Parse("1")
				Expect(err).To(Succeed())
				Expect(v).To(Equal(interface{}("1")))
			})
		})
		When("type is invalid", func() {
			It("returns err", func() {
				p := params.Param{Type: 0}
				v, err = p.Parse("1")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(MatchRegexp("type error"))
			})
		})
	})
	Context("Parse", func() {
		When("type is ParamTypeInt64, ParamTypeString. it is Required", func() {
			It(`returns ".Name"`, func() {
				p := params.Param{Name: "int", Type: params.ParamTypeInt64, Required: true}
				Expect(p.String()).To(Equal("int"))
				p = params.Param{Name: "string", Type: params.ParamTypeString, Required: true}
				Expect(p.String()).To(Equal("string"))
			})
		})
		When("type is ParamTypeInt64, ParamTypeString. it is not Required", func() {
			It(`returns "[ .Name ]"`, func() {
				p := params.Param{Name: "int", Type: params.ParamTypeInt64, Required: false}
				Expect(p.String()).To(Equal("[ int ]"))
				p = params.Param{Name: "string", Type: params.ParamTypeString, Required: false}
				Expect(p.String()).To(Equal("[ string ]"))
			})
		})
		When("type is ParamTypeArrayInt64, ParamTyperrayString. it is Required", func() {
			It(`returns ".Name ..."`, func() {
				p := params.Param{Name: "int", Type: params.ParamTypeArrayInt64, Required: true}
				Expect(p.String()).To(Equal("int ..."))
				p = params.Param{Name: "string", Type: params.ParamTypeArrayString, Required: true}
				Expect(p.String()).To(Equal("string ..."))
			})
		})
		When("type is ParamTypeArrayInt64, ParamTyperrayString. it is not Required", func() {
			It(`returns ".Name ..."`, func() {
				p := params.Param{Name: "int", Type: params.ParamTypeArrayInt64, Required: false}
				Expect(p.String()).To(Equal("[ int ... ]"))
				p = params.Param{Name: "string", Type: params.ParamTypeArrayString, Required: false}
				Expect(p.String()).To(Equal("[ string ... ]"))
			})
		})
	})
})
