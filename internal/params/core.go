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
	contractID := Param{Name: "ContractID", Type: ParamTypeString, Required: true}
	// 		core            contracts                       Contract                        Read    Update
	//		core            contracts                       ContractList                    List
	list = append(list, &APISet{
		Name:        "contracts",
		Description: "IIJ DNS Platform Service information.",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &core.ContractList{}},
			api.ActionRead:   {Object: &core.Contract{}, Params: Params{contractID}},
			api.ActionUpdate: {Object: &core.Contract{}, Params: Params{contractID}},
		},
	})
	// core            delegations                     DelegationList                  List
	list = append(list, &APISet{
		Name:        "delegations",
		Description: "Get domain name management contract and Line contract.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &core.DelegationList{}},
		},
	})
	// core            delegations_request             DelegationApply                 Apply
	list = append(list, &APISet{
		Name:        "delegations_request",
		Description: "Update name server.",
		Action: map[api.Action]API{
			api.ActionApply: {
				Object: &core.DelegationApply{},
				SetFunc: func(spec apis.Spec, args []string) error {
					apply, ok := spec.(*core.DelegationApply)
					if !ok {
						return fmt.Errorf("why this code running")
					}
					apply.ZoneIDs = append(apply.ZoneIDs, args...)
					return nil
				},
				Params: Params{{Name: "servicecode", Type: ParamTypeArrayString, Required: true}},
			},
		},
	})

	requestID := Param{Name: "RequestID", Type: ParamTypeString, Required: true}
	// core            jobs                            Job                             Read
	list = append(list, &APISet{
		Name:        "jobs",
		Description: "Asynchronous responses can be obtained during job processing.",
		Action: map[api.Action]API{
			api.ActionRead: {Object: &core.Job{}, Params: Params{requestID}},
		},
	})

	zoneID := Param{Name: "ZoneID", Type: ParamTypeString, Required: true}
	// core            zones                           ZoneList                        List
	// core            zones                           Zone                            Read    Update  Cancel
	// core            zones_apply                     ZoneApply                       Apply
	list = append(list, &APISet{
		Name:        "zones",
		Description: "IIJ Managed DNS Service information",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &core.ZoneList{}},
			api.ActionRead:   {Object: &core.Zone{}, Params: Params{zoneID}},
			api.ActionUpdate: {Object: &core.Zone{}, Params: Params{zoneID}},
			api.ActionCancel: {Object: &zones.ZoneApply{}, Params: Params{zoneID}, Desc: "Reset edited records."},
			api.ActionApply:  {Object: &zones.ZoneApply{}, Params: Params{zoneID}, Desc: "Apply edited records."},
		},
	})
	SetGroup("core", list)

}
