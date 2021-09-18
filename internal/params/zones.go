package params

import (
	. "github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/zones"
)

func init() {
	list := APISetSlice{}
	zoneId := Param{Name: "ZoneId", Type: ParamTypeString, Required: true}

	// zones           current_records                 CurrentRecordList               List
	list = append(list, &APISet{
		Name:        "current_records",
		Description: "current running records.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.CurrentRecordList{}, Params: Params{zoneId}},
		},
	})

	// zones           dnssec                          Dnssec                          Read    Update
	list = append(list, &APISet{
		Name:        "dnssec",
		Description: "dnssec setting",
		Action: map[api.Action]API{
			api.ActionRead:   {Object: &zones.Dnssec{}, Params: Params{zoneId}},
			api.ActionUpdate: {Object: &zones.Dnssec{}, Params: Params{zoneId}},
		},
	})

	// zones           ds_records                      DsRecordList                    List
	list = append(list, &APISet{
		Name:        "ds_records",
		Description: "get DS record",
		Action: map[api.Action]API{
			api.ActionRead: {Object: &zones.DsRecordList{}, Params: Params{zoneId}},
		},
	})

	// zones           ksk_roll_over                   DnssecKskRollover               Apply
	list = append(list, &APISet{
		Name:        "ksk_roll_over",
		Description: "Running KSK roll over",
		Action: map[api.Action]API{
			api.ActionApply: {Object: &zones.DnssecKskRollover{}, Params: Params{zoneId}},
		},
	})

	// zones           managed_dns_servers             ManagedDnsList                  List
	list = append(list, &APISet{
		Name:        "managed_dns_servers",
		Description: "get managed dns server",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.ManagedDnsList{}, Params: Params{zoneId}},
		},
	})

	// zones           record_diffs                    RecordDiffList                  List
	list = append(list, &APISet{
		Name:        "record_diffs",
		Description: "records diff between running records and edit records.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.RecordDiffList{}, Params: Params{zoneId}},
		},
	})

	recordId := Param{Name: "RecordId", Type: ParamTypeString, Required: true}

	// zones           records                         Record                          Create  Read    Update  Delete  Cancel
	// zones           records                         RecordList                      List
	list = append(list, &APISet{
		Name:        "records",
		Description: "record",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &zones.RecordList{}, Params: Params{zoneId}},
			api.ActionCreate: {Object: &zones.Record{}, Params: Params{zoneId}},
			api.ActionRead:   {Object: &zones.Record{}, Params: Params{zoneId, recordId}},
			api.ActionUpdate: {Object: &zones.Record{}, Params: Params{zoneId, recordId}},
			api.ActionDelete: {Object: &zones.Record{}, Params: Params{zoneId, recordId}},
		},
	})

	// zones           zone_default_ttl                DefaultTTL                      Read    Update  Cancel
	list = append(list, &APISet{
		Name:        "zone_default_ttl",
		Description: "default ttl setting.",
		Action: map[api.Action]API{
			api.ActionRead:   {Object: &zones.DefaultTTL{}, Params: Params{zoneId}},
			api.ActionUpdate: {Object: &zones.DefaultTTL{}, Params: Params{zoneId}},
			api.ActionDelete: {Object: &zones.DefaultTTL{}, Params: Params{zoneId}},
		},
	})

	// zones           zone_default_ttl_diffs          DefaultTTLDiffList              List
	list = append(list, &APISet{
		Name:        "zone_default_ttl_diffs",
		Description: "default ttl diff between running and edited.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.DefaultTTLDiffList{}, Params: Params{zoneId}},
		},
	})

	// zones           zone_histories                  HistoryList                     List
	list = append(list, &APISet{
		Name:        "zone_histories",
		Description: "zone record histories.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.HistoryList{}, Params: Params{zoneId}},
		},
	})

	// zones           zone_logs                       LogList                         List
	list = append(list, &APISet{
		Name:        "zone_logs",
		Description: "IIJ Managed DNS Service log.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.HistoryList{}, Params: Params{zoneId}},
		},
	})

	// zones           zone_proxy                      ZoneProxy                       Read    Update
	list = append(list, &APISet{
		Name:        "zone_proxy",
		Description: "IIJ Managed DNS Service log.",
		Action: map[api.Action]API{
			api.ActionRead:   {Object: &zones.ZoneProxy{}, Params: Params{zoneId}},
			api.ActionUpdate: {Object: &zones.ZoneProxy{}, Params: Params{zoneId}},
		},
	})

	// zones           zone_proxy_health_check         ZoneProxyHealthCheckList        List
	list = append(list, &APISet{
		Name:        "zone_proxy_health_check",
		Description: "zone proxy health check information",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.ZoneProxyHealthCheckList{}, Params: Params{zoneId}},
		},
	})

	// zones           zones_contract                  Contract                        Read
	list = append(list, &APISet{
		Name:        "zones_contract",
		Description: "the contract belong to zone",
		Action: map[api.Action]API{
			api.ActionRead: {Object: &zones.Contract{}, Params: Params{zoneId}},
		},
	})

	SetGroup("zones", list)
}
