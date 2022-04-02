package params

import (
	. "github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/contracts"
)

func init() {
	list := APISetSlice{}
	contractID := Param{Name: "ContractID", Type: ParamTypeString, Required: true}
	commonConfigID := Param{Name: "CommonConfigID", Type: ParamTypeInt64, Required: true}
	// contracts       common_configs                  CommonConfig                    Create  Read    Update  Delete
	// contracts       common_configs                  CommonConfigList                List
	list = append(list, &APISet{
		Name:        "common_configs",
		Description: "Common config setting.",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &contracts.CommonConfigList{}, Params: Params{contractID}},
			api.ActionCreate: {Object: &contracts.CommonConfig{}, Params: Params{contractID}},
			api.ActionRead:   {Object: &contracts.CommonConfig{}, Params: Params{contractID, commonConfigID}},
			api.ActionUpdate: {Object: &contracts.CommonConfig{}, Params: Params{contractID, commonConfigID}},
			api.ActionDelete: {Object: &contracts.CommonConfig{}, Params: Params{contractID, commonConfigID}},
		},
	})
	// contracts       common_configs_default          CommonConfigDefault             Apply
	list = append(list, &APISet{
		Name:        "common_configs_default",
		Description: "Common config setting an initial value for IIJ Managed DNS Service.",
		Action: map[api.Action]API{
			api.ActionApply: {Object: &contracts.CommonConfigDefault{}, Params: Params{contractID, commonConfigID}},
		},
	})

	// contracts       common_configs_managed_dns      CommonConfigManagedDns          Apply
	list = append(list, &APISet{
		Name:        "common_configs_managed_dns",
		Description: "Managed DNS Server setting (advanced configuration.",
		Action: map[api.Action]API{
			api.ActionApply: {Object: &contracts.CommonConfigManagedDns{}, Params: Params{contractID, commonConfigID}},
		},
	})

	// contracts       contract_logs                   LogList                         List
	list = append(list, &APISet{
		Name:        "contract_logs",
		Description: "IIJ DNS Platform Service operation log.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &contracts.LogList{}, Params: Params{contractID}},
		},
	})

	// contracts       contract_partners               CommonConfigList                List
	list = append(list, &APISet{
		Name:        "contract_partners",
		Description: "List of linked services.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &contracts.ContractPartnerList{}, Params: Params{contractID}},
		},
	})

	// contracts       contract_qps                    QpsHistoryList                  List
	list = append(list, &APISet{
		Name:        "contract_qps",
		Description: "QPS information.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &contracts.QpsHistoryList{}, Params: Params{contractID}},
		},
	})

	// contracts       contract_zone                   ContractZoneList                List
	list = append(list, &APISet{
		Name:        "contract_zone",
		Description: "A list of common configs belongs to contract",
		Action: map[api.Action]API{
			api.ActionList: {Object: &contracts.ContractZoneList{}, Params: Params{contractID}},
		},
	})

	zoneIDs := Param{Name: "ZoneIDs", Type: ParamTypeArrayString, Required: true}
	// contracts       contract_zone_common_configs    ContractZoneCommonConfig        Apply
	list = append(list, &APISet{
		Name:        "contract_zone_common_configs",
		Description: "Settings common config to IIJ Managed DNS.",
		Action: map[api.Action]API{
			api.ActionApply: {Object: &contracts.ContractZoneCommonConfig{}, Params: Params{contractID, commonConfigID, zoneIDs}},
		},
	})

	tsigID := Param{Name: "TsigID", Type: ParamTypeInt64, Required: true}

	// contracts       tsigs                           Tsig                            Create  Read    Update  Delete
	// contracts       tsigs                           TsigList                        List
	list = append(list, &APISet{
		Name:        "tsigs",
		Description: "TSIG settings.",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &contracts.TsigList{}, Params: Params{contractID}},
			api.ActionCreate: {Object: &contracts.Tsig{}, Params: Params{contractID}},
			api.ActionRead:   {Object: &contracts.Tsig{}, Params: Params{contractID, tsigID}},
			api.ActionUpdate: {Object: &contracts.Tsig{}, Params: Params{contractID, tsigID}},
			api.ActionDelete: {Object: &contracts.Tsig{}, Params: Params{contractID, tsigID}},
		},
	})

	// contracts       tsigs_common_configs            TsigCommonConfigList            List
	list = append(list, &APISet{
		Name:        "tsigs_common_configs",
		Description: "A list of common configs belongs to TSIG.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &contracts.TsigCommonConfigList{}, Params: Params{contractID, tsigID}},
		},
	})

	SetGroup("contracts", list)
}
