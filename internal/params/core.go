package params

import (
	"fmt"

	. "github.com/mimuret/dpfctl/pkg/params"

	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/core"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/zones"
)

func init() {
	list := APISetSlice{}
	contractId := Param{Name: "ContractId", Type: ParamTypeString, Required: true}
	// 		core            contracts                       Contract                        Read    Update
	//		core            contracts                       ContractList                    List
	list = append(list, &APISet{
		Name:        "contracts",
		Description: "IIJ DNS Platform Service contract information",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &core.ContractList{}},
			api.ActionRead:   {Object: &core.Contract{}, Params: Params{contractId}},
			api.ActionUpdate: {Object: &core.Contract{}, Params: Params{contractId}},
		},
	})
	// core            delegations                     DelegationList                  List
	list = append(list, &APISet{
		Name:        "delegations",
		Description: "get domain name management contract and Line contract",
		Action: map[api.Action]API{
			api.ActionList: {Object: &core.DelegationList{}},
		},
	})
	// core            delegations_request             DelegationApply                 Apply
	list = append(list, &APISet{
		Name:        "delegations_request",
		Description: "update name server",
		Action: map[api.Action]API{
			api.ActionApply: {
				Object: &core.DelegationApply{},
				SetFunc: func(spec apis.Spec, args []string) error {
					apply, ok := spec.(*core.DelegationApply)
					if !ok {
						return fmt.Errorf("why this code running")
					}
					apply.ZoneIds = append(apply.ZoneIds, args...)
					return nil
				},
				Params: Params{{Name: "servicecode", Type: ParamTypeArrayString, Required: true}},
			},
		},
	})

	requestId := Param{Name: "RequestId", Type: ParamTypeString, Required: true}
	// core            jobs                            Job                             Read
	list = append(list, &APISet{
		Name:        "jobs",
		Description: "Asynchronous responses can be obtained during job processing. You can get it only once after the job is finished.",
		Action: map[api.Action]API{
			api.ActionRead: {Object: &core.Job{}, Params: Params{requestId}},
		},
	})

	zoneId := Param{Name: "ZoneId", Type: ParamTypeString, Required: true}
	// core            zones                           ZoneList                        List
	// core            zones                           Zone                            Read    Update  Cancel
	// core            zones_apply                     ZoneApply                       Apply
	list = append(list, &APISet{
		Name:        "zones",
		Description: "IIJ Managed DNS Service information",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &core.ZoneList{}},
			api.ActionRead:   {Object: &core.Zone{}, Params: Params{zoneId}},
			api.ActionUpdate: {Object: &core.Zone{}, Params: Params{zoneId}},
			api.ActionCancel: {Object: &zones.ZoneApply{}, Params: Params{zoneId}},
			api.ActionApply:  {Object: &zones.ZoneApply{}, Params: Params{zoneId}},
		},
	})
	SetGroup("core", list)

}
