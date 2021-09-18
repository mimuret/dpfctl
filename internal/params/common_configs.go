package params

import (
	. "github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/common_configs"
)

func init() {
	list := APISetSlice{}
	CommonConfigId := Param{Name: "CommonConfigId", Type: ParamTypeInt64, Required: true}
	ccPrimaryId := Param{Name: "CcPrimaryId", Type: ParamTypeInt64, Required: true}
	// common_configs  cc_primaries                    CcPrimary                       Create  Read    Update  Delete
	// common_configs  cc_primaries                    CcPrimaryList                   List
	list = append(list, &APISet{
		Name:        "cc_primaries",
		Description: "Primary DNS settings, using by zone proxy (DNS Secondary) mode.",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &common_configs.CcPrimaryList{}, Params: Params{CommonConfigId}},
			api.ActionCreate: {Object: &common_configs.CcPrimary{}, Params: Params{CommonConfigId}},
			api.ActionRead:   {Object: &common_configs.CcPrimary{}, Params: Params{CommonConfigId, ccPrimaryId}},
			api.ActionUpdate: {Object: &common_configs.CcPrimary{}, Params: Params{CommonConfigId, ccPrimaryId}},
			api.ActionDelete: {Object: &common_configs.CcPrimary{}, Params: Params{CommonConfigId, ccPrimaryId}},
		},
	})

	ccSecNotifiedServerId := Param{Name: "CcSecNotifiedServerId", Type: ParamTypeInt64, Required: true}
	// common_configs  cc_sec_notified_servers         CcSecNotifiedServer             Create  Read    Update  Delete
	// common_configs  cc_sec_notified_servers         CcSecNotifiedServerList         List

	list = append(list, &APISet{
		Name:        "cc_sec_notified_servers",
		Description: "Secondary DNS function notify setting. ",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &common_configs.CcSecNotifiedServerList{}, Params: Params{CommonConfigId}},
			api.ActionCreate: {Object: &common_configs.CcSecNotifiedServer{}, Params: Params{CommonConfigId}},
			api.ActionRead:   {Object: &common_configs.CcSecNotifiedServer{}, Params: Params{CommonConfigId, ccSecNotifiedServerId}},
			api.ActionUpdate: {Object: &common_configs.CcSecNotifiedServer{}, Params: Params{CommonConfigId, ccSecNotifiedServerId}},
			api.ActionDelete: {Object: &common_configs.CcSecNotifiedServer{}, Params: Params{CommonConfigId, ccSecNotifiedServerId}},
		},
	})

	CcSecTransferAclId := Param{Name: "CcSecTransferAclId", Type: ParamTypeInt64, Required: true}
	// common_configs  cc_sec_transfer_acls            CcSecTransferAcl                Create  Read    Update  Delete
	// common_configs  cc_sec_transfer_acls            CcSecTransferAclList            List
	list = append(list, &APISet{
		Name:        "cc_sec_transfer_acls",
		Description: "Secondary DNS function allow zone transer setting. ",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &common_configs.CcSecTransferAclList{}, Params: Params{CommonConfigId}},
			api.ActionCreate: {Object: &common_configs.CcSecTransferAcl{}, Params: Params{CommonConfigId}},
			api.ActionRead:   {Object: &common_configs.CcSecTransferAcl{}, Params: Params{CommonConfigId, CcSecTransferAclId}},
			api.ActionUpdate: {Object: &common_configs.CcSecTransferAcl{}, Params: Params{CommonConfigId, CcSecTransferAclId}},
			api.ActionDelete: {Object: &common_configs.CcSecTransferAcl{}, Params: Params{CommonConfigId, CcSecTransferAclId}},
		},
	})
	SetGroup("commoon_conifgs", list)
}
