package zabbix

import "fmt"

// https://www.zabbix.com/documentation/3.4/en/manual/api/reference/host/object
type Host struct {
	HostId string `json:"hostid"`
	Host   string `json:"host"`
}

type Hosts []Host

// Wrapper for host.get: https://www.zabbix.com/documentation/3.4/en/manual/api/reference/host/get
func (api *API) HostsGet(params Params) (Hosts, error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	resp, err := api.CallWithError("host.get", params)
	if err != nil {
		return nil, fmt.Errorf("api.CallWithError: %v", err)
	}

	res := Hosts{}
	err = resp.Bind(&res)
	if err != nil {
		return nil, fmt.Errorf("resp.Bind: %v", err)
	}
	return res, nil
}

// Gets hosts by host group Ids.
func (api *API) HostsGetByHostGroupIds(ids []string) (Hosts, error) {
	return api.HostsGet(
		Params{
			"groupids": ids,
			"output": []string{
				"hostid",
				"host",
			}})
}
