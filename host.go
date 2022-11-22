package zabbix

// https://www.zabbix.com/documentation/3.4/en/manual/api/reference/host/object
type Host struct {
	HostId string `json:"hostid"`
	Host   string `json:"host"`
}

type Hosts []Host

// Wrapper for host.get: https://www.zabbix.com/documentation/3.4/en/manual/api/reference/host/get
func (api *API) HostsGet(params Params) (res Hosts, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("host.get", params)
	if err != nil {
		return
	}

	err = response.Bind(&res)
	if err != nil {
		return
	}
	return
}

// Gets hosts by host group Ids.
func (api *API) HostsGetByHostGroupIds(ids []string) (res Hosts, err error) {
	return api.HostsGet(
		Params{
			"groupids": ids,
			"output": []string{
				"hostid",
				"host",
			}})
}
