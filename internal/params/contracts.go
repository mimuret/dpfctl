package params

import (
	. "github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/contracts"
)

func init() {
	list := APISetSlice{}
	contractId := Param{Name: "contractId", Type: ParamTypeString, Required: true}
	commonConfigId := Param{Name: "CommonConfigId", Type: ParamTypeInt64, Required: true}
	// contracts       common_configs                  CommonConfig                    Create  Read    Update  Delete
	// contracts       common_configs                  CommonConfigList                List
	list = append(list, &APISet{
		Name:        "common_configs",
		Description: "Common config setting.",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &contracts.CommonConfigList{}, Params: Params{contractId}},
			api.ActionCreate: {Object: &contracts.CommonConfig{}, Params: Params{contractId}},
			api.ActionRead:   {Object: &contracts.CommonConfig{}, Params: Params{contractId, commonConfigId}},
			api.ActionUpdate: {Object: &contracts.CommonConfig{}, Params: Params{contractId, commonConfigId}},
			api.ActionDelete: {Object: &contracts.CommonConfig{}, Params: Params{contractId, commonConfigId}},
		},
	})
	// contracts       common_configs_default          CommonConfigDefault             Apply
	list = append(list, &APISet{
		Name:        "common_configs_default",
		Description: "change new contract initial common_config",
		Action: map[api.Action]API{
			api.ActionApply: {Object: &contracts.CommonConfigDefault{}, Params: Params{contractId, commonConfigId}},
		},
	})

	// contracts       common_configs_managed_dns      CommonConfigManagedDns          Apply
	list = append(list, &APISet{
		Name:        "common_configs_managed_dns",
		Description: "managed dns function setting that is part of common config.",
		Action: map[api.Action]API{
			api.ActionApply: {Object: &contracts.CommonConfigManagedDns{}, Params: Params{contractId, commonConfigId}},
		},
	})

	// contracts       contract_logs                   LogList                         List
	list = append(list, &APISet{
		Name:        "contract_logs",
		Description: "IIJ DNS Platform Service log.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &contracts.LogList{}, Params: Params{contractId}},
		},
	})

	// contracts       contract_partners               CommonConfigList                List
	list = append(list, &APISet{
		Name:        "contract_partners",
		Description: "list of linked services.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &contracts.ContractPartnerList{}, Params: Params{contractId}},
		},
	})

	// contracts       contract_qps                    QpsHistoryList                  List
	list = append(list, &APISet{
		Name:        "contract_qps",
		Description: "query per sec information.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &contracts.QpsHistoryList{}, Params: Params{contractId}},
		},
	})

	// contracts       contract_zone                   ContractZoneList                List
	list = append(list, &APISet{
		Name:        "contract_zone",
		Description: " list of zones belong to contract",
		Action: map[api.Action]API{
			api.ActionList: {Object: &contracts.ContractZoneList{}, Params: Params{contractId}},
		},
	})

	zoneIds := Param{Name: "ZoneIds", Type: ParamTypeArrayString, Required: true}
	// contracts       contract_zone_common_configs    ContractZoneCommonConfig        Apply
	list = append(list, &APISet{
		Name:        "contract_zone_common_configs",
		Description: "managed dns function setting that is part of common config.",
		Action: map[api.Action]API{
			api.ActionApply: {Object: &contracts.ContractZoneCommonConfig{}, Params: Params{contractId, commonConfigId, zoneIds}},
		},
	})

	tsigId := Param{Name: "TsigId", Type: ParamTypeInt64, Required: true}

	// contracts       tsigs                           Tsig                            Create  Read    Update  Delete
	// contracts       tsigs                           TsigList                        List
	list = append(list, &APISet{
		Name:        "tsigs",
		Description: "TSIG setting.",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &contracts.TsigList{}, Params: Params{contractId}},
			api.ActionCreate: {Object: &contracts.Tsig{}, Params: Params{contractId}},
			api.ActionRead:   {Object: &contracts.Tsig{}, Params: Params{contractId, tsigId}},
			api.ActionUpdate: {Object: &contracts.Tsig{}, Params: Params{contractId, tsigId}},
			api.ActionDelete: {Object: &contracts.Tsig{}, Params: Params{contractId, tsigId}},
		},
	})

	// contracts       tsigs_common_configs            TsigCommonConfigList            List
	list = append(list, &APISet{
		Name:        "tsigs_common_configs",
		Description: " common_configs of using to TSIG",
		Action: map[api.Action]API{
			api.ActionList: {Object: &contracts.TsigCommonConfigList{}, Params: Params{contractId, tsigId}},
		},
	})

	SetGroup("contracts", list)
}
