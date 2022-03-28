package params

import (
	. "github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/common_configs"
)

func init() {
	list := APISetSlice{}
	commonConfigID := Param{Name: "CommonConfigID", Type: ParamTypeInt64, Required: true}
	ccPrimaryID := Param{Name: "CcPrimaryID", Type: ParamTypeInt64, Required: true}
	// common_configs  cc_primaries                    CcPrimary                       Create  Read    Update  Delete
	// common_configs  cc_primaries                    CcPrimaryList                   List
	list = append(list, &APISet{
		Name:        "cc_primaries",
		Description: "Primary DNS settings, using by zone proxy (DNS Secondary) mode.",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &common_configs.CcPrimaryList{}, Params: Params{commonConfigID}},
			api.ActionCreate: {Object: &common_configs.CcPrimary{}, Params: Params{commonConfigID}},
			api.ActionRead:   {Object: &common_configs.CcPrimary{}, Params: Params{commonConfigID, ccPrimaryID}},
			api.ActionUpdate: {Object: &common_configs.CcPrimary{}, Params: Params{commonConfigID, ccPrimaryID}},
			api.ActionDelete: {Object: &common_configs.CcPrimary{}, Params: Params{commonConfigID, ccPrimaryID}},
		},
	})

	ccSecNotifiedServerID := Param{Name: "CcSecNotifiedServerID", Type: ParamTypeInt64, Required: true}
	// common_configs  cc_sec_notified_servers         CcSecNotifiedServer             Create  Read    Update  Delete
	// common_configs  cc_sec_notified_servers         CcSecNotifiedServerList         List

	list = append(list, &APISet{
		Name:        "cc_sec_notified_servers",
		Description: "Secondary DNS function notify setting. ",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &common_configs.CcSecNotifiedServerList{}, Params: Params{commonConfigID}},
			api.ActionCreate: {Object: &common_configs.CcSecNotifiedServer{}, Params: Params{commonConfigID}},
			api.ActionRead:   {Object: &common_configs.CcSecNotifiedServer{}, Params: Params{commonConfigID, ccSecNotifiedServerID}},
			api.ActionUpdate: {Object: &common_configs.CcSecNotifiedServer{}, Params: Params{commonConfigID, ccSecNotifiedServerID}},
			api.ActionDelete: {Object: &common_configs.CcSecNotifiedServer{}, Params: Params{commonConfigID, ccSecNotifiedServerID}},
		},
	})

	CcSecTransferAclID := Param{Name: "CcSecTransferAclID", Type: ParamTypeInt64, Required: true}
	// common_configs  cc_sec_transfer_acls            CcSecTransferAcl                Create  Read    Update  Delete
	// common_configs  cc_sec_transfer_acls            CcSecTransferAclList            List
	list = append(list, &APISet{
		Name:        "cc_sec_transfer_acls",
		Description: "Secondary DNS function allow zone transer setting. ",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &common_configs.CcSecTransferAclList{}, Params: Params{commonConfigID}},
			api.ActionCreate: {Object: &common_configs.CcSecTransferAcl{}, Params: Params{commonConfigID}},
			api.ActionRead:   {Object: &common_configs.CcSecTransferAcl{}, Params: Params{commonConfigID, CcSecTransferAclID}},
			api.ActionUpdate: {Object: &common_configs.CcSecTransferAcl{}, Params: Params{commonConfigID, CcSecTransferAclID}},
			api.ActionDelete: {Object: &common_configs.CcSecTransferAcl{}, Params: Params{commonConfigID, CcSecTransferAclID}},
		},
	})
	SetGroup("commoon_configs", list)
}
