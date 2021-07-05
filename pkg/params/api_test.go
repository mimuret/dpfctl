package params_test

import (
	"testing"

	"github.com/mimuret/dpfctl/pkg/params"
	"github.com/stretchr/testify/assert"
)

func TestGroupMap(t *testing.T) {
	a := assert.New(t)

	gm := params.NewGroupMap()
	bookSet := &params.APISet{
		Name:        "book",
		Description: "hoge",
	}
	list := params.APISetSlice{nil, bookSet}
	gm.SetGroup("test", list)
	a.Equal(gm.GetAPISlice("test"), list)
	a.Equal(gm.GetAPISetfromCmdName("book"), bookSet)
}
