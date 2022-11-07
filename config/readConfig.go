package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var ipv6_routing []interface{}
var srv6_end map[string]interface{}
var srv6_insert []interface{}

var Ethernet_dstAddr string
var Ipv6_routing []map[string]interface{}
var Srv6_end map[string]interface{}
var Srv6_insert []map[string]interface{}

type Cfg interface {
	GetEthernetAddr()
	GetIpv6Routing()
	GetSrv6End()
	GetSrv6Insert()
}

func ReadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	Ethernet_dstAddr = viper.GetString("ethernet_dstAddr")
	ipv6_routing = viper.Get("ipv6_routing").([]interface{})
	for i := 0; i < len(ipv6_routing); i++ {
		x := ipv6_routing[i].(map[string]interface{})
		Ipv6_routing = append(Ipv6_routing, x)
	}
	Srv6_end = viper.Get("srv6_end").(map[string]interface{})
	//fmt.Println(Srv6_end["ipv6_dstaddr"])
	srv6_insert = viper.Get("srv6_insert").([]interface{})
	for i := 0; i < len(srv6_insert); i++ {
		x := srv6_insert[i].(map[string]interface{})
		//fmt.Println(x["ipv6_srcAddr"])
		Srv6_insert = append(Srv6_insert, x)
	}

}
