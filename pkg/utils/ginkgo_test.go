package utils_test

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/testtool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	paramsGroupMapDefault *params.GroupMap
)

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "utils Suite")
}

var _ = BeforeSuite(func() {
	utils.NewClient = func(logger api.Logger) (api.ClientInterface, error) {
		return testtool.NewTestClient("token", "http://localhost", logger), nil
	}
	paramsGroupMapDefault = params.GroupMapDefault
	params.GroupMapDefault = params.NewGroupMap()
	list := params.APISetSlice{}
	Id := params.Param{Name: "Id", Type: params.ParamTypeString, Required: true}
	list = append(list, &params.APISet{
		Name:        "test1",
		Description: "test1 description",
		Action: map[api.Action]params.API{
			api.ActionList:   {Object: &testtool.TestSpecList{}},
			api.ActionCreate: {Object: &testtool.TestSpec{}},
			api.ActionRead:   {Object: &testtool.TestSpec{}, Params: params.Params{Id}},
			api.ActionUpdate: {Object: &testtool.TestSpec{}, Params: params.Params{Id}},
			api.ActionDelete: {Object: &testtool.TestSpec{}, Params: params.Params{Id}},
			api.ActionApply:  {Object: &testtool.TestSpec{}, Params: params.Params{Id}},
			api.ActionCancel: {Object: &testtool.TestSpec{}, Params: params.Params{Id}},
		},
	})
	list = append(list, &params.APISet{
		Name:        "test2",
		Description: "test2 description",
		Action: map[api.Action]params.API{
			api.ActionList: {Object: &testtool.TestSpecList{}},
		},
	})
	list = append(list, &params.APISet{
		Name:        "test3",
		Description: "test3 description",
		Action: map[api.Action]params.API{
			api.ActionRead: {Object: &testtool.TestSpec{}, Params: params.Params{Id}},
		},
	})
	params.SetGroup("tests", list)

	httpmock.Activate()
})

var _ = BeforeEach(func() {
	httpmock.Reset()
	params.GroupMapDefault = paramsGroupMapDefault
})

var _ = AfterSuite(func() {
	utils.NewClient = utils.NewClientDefault
	httpmock.DeactivateAndReset()
})
