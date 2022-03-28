package utils_test

import (
	_ "embed"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/dpfctl/pkg/utils"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/testtool"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var (
	paramsGroupMapDefault *params.GroupMap
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

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "utils Suite")
}

var _ = BeforeSuite(func() {
	utils.DefaultFS = afero.NewMemMapFs()
	afero.WriteFile(utils.DefaultFS, "testdata/bad-schema.yaml", badYamlDoc, 0644)
	afero.WriteFile(utils.DefaultFS, "testdata/bad.yaml", badYamlDoc, 0644)
	afero.WriteFile(utils.DefaultFS, "testdata/multi-doc.yaml", multiYamlDoc, 0644)
	afero.WriteFile(utils.DefaultFS, "testdata/single-doc.yaml", singleYamlDoc, 0644)
	afero.WriteFile(utils.DefaultFS, "testdata/single-doc.json", jsonDoc, 0644)

	utils.NewClient = func(logger api.Logger) (api.ClientInterface, error) {
		return testtool.NewTestClient("token", "http://localhost", logger), nil
	}
	paramsGroupMapDefault = params.GroupMapDefault
	params.GroupMapDefault = params.NewGroupMap()
	list := params.APISetSlice{}
	ID := params.Param{Name: "ID", Type: params.ParamTypeString, Required: true}
	list = append(list, &params.APISet{
		Name:        "test1",
		Description: "test1 description",
		Action: map[api.Action]params.API{
			api.ActionList:   {Object: &testtool.TestSpecList{}},
			api.ActionCreate: {Object: &testtool.TestSpec{}},
			api.ActionRead:   {Object: &testtool.TestSpec{}, Params: params.Params{ID}},
			api.ActionUpdate: {Object: &testtool.TestSpec{}, Params: params.Params{ID}},
			api.ActionDelete: {Object: &testtool.TestSpec{}, Params: params.Params{ID}},
			api.ActionApply:  {Object: &testtool.TestSpec{}, Params: params.Params{ID}},
			api.ActionCancel: {Object: &testtool.TestSpec{}, Params: params.Params{ID}},
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
			api.ActionRead: {Object: &testtool.TestSpec{}, Params: params.Params{ID}},
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
	utils.DefaultFS = afero.NewOsFs()
	utils.NewClient = utils.NewClientDefault
	httpmock.DeactivateAndReset()
})
