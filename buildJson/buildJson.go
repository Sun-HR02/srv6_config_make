package buildJson

import (
	"SRv6.Config.Builder/config"
	"encoding/json"
	"strconv"
)

type MatchItem interface {
}
type MatchItemImpl struct {
	Ethernet  string `json:"hdr.ethernet.dst_addr,omitempty"`
	Ipv6      string `json:"hdr.ipv6.dst_addr,omitempty"`
	MatchItem `json:"match_item,omitempty"`
}
type ActionParamItem interface {
}
type ActionParamItemImpl struct {
	S1              string `json:"s1,omitempty"`
	S2              string `json:"s2,omitempty"`
	S3              string `json:"s3,omitempty"`
	S4              string `json:"s4,omitempty"`
	S5              string `json:"s5,omitempty"`
	S6              string `json:"s6,omitempty"`
	S7              string `json:"s7,omitempty"`
	S8              string `json:"s8,omitempty"`
	Dmac            string `json:"dmac,omitempty"`
	Port            string `json:"port,omitempty"`
	ActionParamItem `json:"actionParamItem,omitempty"`
}

type Tables interface {
}
type TablesImpl struct {
	Table         string          `json:"table"`
	Match         MatchItem       `json:"match"`
	Action_name   string          `json:"action_name"`
	Action_params ActionParamItem `json:"action_params"`
	Tables        `json:"tables,omitempty"`
}
type Json interface {
}
type JsonImpl struct {
	Target        string   `json:"target"`
	P4Info        string   `json:"p4info"`
	Bmv2_json     string   `json:"bmv2_json"`
	Table_entries []Tables `json:"table_entries，omitempty"`
	Json          `json:"Json,omitempty"`
}

func WriteJson() ([]byte, error) {
	config.ReadConfig()
	jsontmp := JsonImpl{}
	jsontmp.Target = "bmv2"
	// NOTE: 可能有问题，需要注意
	jsontmp.P4Info = "build/srv6.p4.p4info.txt"
	jsontmp.Bmv2_json = "build/srv6.json"
	jsontmp.Table_entries = BuildTables()
	output, err := json.Marshal(jsontmp)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func BuildTables() []Tables {
	tables := make([]Tables, 0)
	ethernetAddr := config.Ethernet_dstAddr
	Srv6_end_all := config.Srv6_end
	Srv6_insert_all := config.Srv6_insert
	Ipv6_all := config.Ipv6_routing

	// ethernet
	ethernetTable := TablesImpl{
		Table: "ingress.local_mac_table",
		Match: MatchItemImpl{
			Ethernet: ethernetAddr,
		},
		Action_name:   "NoAction",
		Action_params: ActionParamItemImpl{},
	}
	tables = append(tables, ethernetTable)

	// end table
	srv6_ends := Srv6_end_all["ipv6_dstaddr"].([]interface{})
	for _, dstAddr := range srv6_ends {
		dst := dstAddr
		srv6EndTable := TablesImpl{
			Table: "ingress.local_sid_table",
			Match: MatchItemImpl{
				Ipv6: dst.(string),
			},
			Action_name:   "ingress.end",
			Action_params: ActionParamItemImpl{},
		}
		tables = append(tables, srv6EndTable)
	}

	// insert table
	for _, table := range Srv6_insert_all {
		dstAddr := table["ipv6_dstaddr"].(string)
		paramsSlice := table["params"].([]interface{})

		params := ActionParamItemImpl{}
		//fmt.Println(params)
		for i, segment := range paramsSlice {
			switch i {
			case 0:
				params.S1 = segment.(string)
			case 1:
				params.S2 = segment.(string)
			case 2:
				params.S3 = segment.(string)
			case 3:
				params.S4 = segment.(string)
			case 4:
				params.S5 = segment.(string)
			case 5:
				params.S6 = segment.(string)
			case 6:
				params.S7 = segment.(string)
			case 7:
				params.S8 = segment.(string)
			}
		}
		srv6InsertTable := TablesImpl{
			Table:         "ingress.transit_table",
			Match:         dstAddr,
			Action_name:   "ingress.insert_segment_list_" + strconv.Itoa(len(paramsSlice)),
			Action_params: params,
		}
		tables = append(tables, srv6InsertTable)
	}

	// ipv6 table
	for _, table := range Ipv6_all {
		dstAddr := table["ipv6_dstaddr"].(string)
		dstMac := table["dstmac"].(string)
		port := table["port"].(int)
		portStr := strconv.Itoa(port)
		ipv6Table := TablesImpl{
			Table:       "ingress.routing_v6_table",
			Match:       MatchItemImpl{Ipv6: dstAddr},
			Action_name: "ingress.set_next_hop",
			Action_params: ActionParamItemImpl{
				Dmac: dstMac,
				Port: portStr,
			},
		}
		tables = append(tables, ipv6Table)
	}

	return tables
}
