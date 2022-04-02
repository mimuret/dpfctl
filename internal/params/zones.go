package params

import (
	. "github.com/mimuret/dpfctl/pkg/params"
	"github.com/mimuret/golang-iij-dpf/pkg/api"
	"github.com/mimuret/golang-iij-dpf/pkg/apis/dpf/v1/zones"
)

func init() {
	list := APISetSlice{}
	zoneID := Param{Name: "ZoneID", Type: ParamTypeString, Required: true}

	// zones           current_records                 CurrentRecordList               List
	list = append(list, &APISet{
		Name:        "current_records",
		Description: "A list of current running records.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.CurrentRecordList{}, Params: Params{zoneID}},
		},
	})

	// zones           dnssec                          Dnssec                          Read    Update
	list = append(list, &APISet{
		Name:        "dnssec",
		Description: "DNSSEC setting",
		Action: map[api.Action]API{
			api.ActionRead:   {Object: &zones.Dnssec{}, Params: Params{zoneID}},
			api.ActionUpdate: {Object: &zones.Dnssec{}, Params: Params{zoneID}},
		},
	})

	// zones           ds_records                      DsRecordList                    List
	list = append(list, &APISet{
		Name:        "ds_records",
		Description: "DS records information",
		Action: map[api.Action]API{
			api.ActionRead: {Object: &zones.DsRecordList{}, Params: Params{zoneID}},
		},
	})

	// zones           ksk_roll_over                   DnssecKskRollover               Apply
	list = append(list, &APISet{
		Name:        "ksk_roll_over",
		Description: "KSK rollover",
		Action: map[api.Action]API{
			api.ActionApply: {Object: &zones.DnssecKskRollover{}, Params: Params{zoneID}},
		},
	})

	// zones           managed_dns_servers             ManagedDnsList                  List
	list = append(list, &APISet{
		Name:        "managed_dns_servers",
		Description: "Managed DNS server",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.ManagedDnsList{}, Params: Params{zoneID}},
		},
	})

	// zones           record_diffs                    RecordDiffList                  List
	list = append(list, &APISet{
		Name:        "record_diffs",
		Description: "records diff between running records and editing records.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.RecordDiffList{}, Params: Params{zoneID}},
		},
	})

	recordID := Param{Name: "RecordID", Type: ParamTypeString, Required: true}

	// zones           records                         Record                          Create  Read    Update  Delete  Cancel
	// zones           records                         RecordList                      List
	list = append(list, &APISet{
		Name:        "records",
		Description: "A list of editing records.",
		Action: map[api.Action]API{
			api.ActionList:   {Object: &zones.RecordList{}, Params: Params{zoneID}},
			api.ActionCreate: {Object: &zones.Record{}, Params: Params{zoneID}},
			api.ActionRead:   {Object: &zones.Record{}, Params: Params{zoneID, recordID}},
			api.ActionUpdate: {Object: &zones.Record{}, Params: Params{zoneID, recordID}},
			api.ActionDelete: {Object: &zones.Record{}, Params: Params{zoneID, recordID}},
		},
	})

	// zones           zone_default_ttl                DefaultTTL                      Read    Update  Cancel
	list = append(list, &APISet{
		Name:        "zone_default_ttl",
		Description: "Default TTL setting.",
		Action: map[api.Action]API{
			api.ActionRead:   {Object: &zones.DefaultTTL{}, Params: Params{zoneID}},
			api.ActionUpdate: {Object: &zones.DefaultTTL{}, Params: Params{zoneID}},
			api.ActionDelete: {Object: &zones.DefaultTTL{}, Params: Params{zoneID}},
		},
	})

	// zones           zone_default_ttl_diffs          DefaultTTLDiffList              List
	list = append(list, &APISet{
		Name:        "zone_default_ttl_diffs",
		Description: "Default TTL diff between running and editing.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.DefaultTTLDiffList{}, Params: Params{zoneID}},
		},
	})

	// zones           zone_histories                  HistoryList                     List
	list = append(list, &APISet{
		Name:        "zone_histories",
		Description: "History of zone committed.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.HistoryList{}, Params: Params{zoneID}},
		},
	})

	// zones           zone_logs                       LogList                         List
	list = append(list, &APISet{
		Name:        "zone_logs",
		Description: "IIJ Managed DNS Service log.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.HistoryList{}, Params: Params{zoneID}},
		},
	})

	// zones           zone_proxy                      ZoneProxy                       Read    Update
	list = append(list, &APISet{
		Name:        "zone_proxy",
		Description: "Zone proxy mode setting.",
		Action: map[api.Action]API{
			api.ActionRead:   {Object: &zones.ZoneProxy{}, Params: Params{zoneID}},
			api.ActionUpdate: {Object: &zones.ZoneProxy{}, Params: Params{zoneID}},
		},
	})

	// zones           zone_proxy_health_check         ZoneProxyHealthCheckList        List
	list = append(list, &APISet{
		Name:        "zone_proxy_health_check",
		Description: "Health check information for primary servers.",
		Action: map[api.Action]API{
			api.ActionList: {Object: &zones.ZoneProxyHealthCheckList{}, Params: Params{zoneID}},
		},
	})

	// zones           zones_contract                  Contract                        Read
	list = append(list, &APISet{
		Name:        "zones_contract",
		Description: "The contract belongs to the zone.",
		Action: map[api.Action]API{
			api.ActionRead: {Object: &zones.Contract{}, Params: Params{zoneID}},
		},
	})

	SetGroup("zones", list)
}
