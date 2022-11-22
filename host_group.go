package zabbix

// https://www.zabbix.com/documentation/2.2/manual/appendix/api/hostgroup/definitions
type HostGroup struct {
	GroupId string `json:"groupid"`
	Name    string `json:"name"`
}

type HostGroups []HostGroup

// Wrapper for hostgroup.get: https://www.zabbix.com/documentation/3.4/en/manual/api/reference/hostgroup/get
func (api *API) HostGroupsGet(params Params) (res HostGroups, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("hostgroup.get", params)
	if err != nil {
		return
	}

	err = response.Bind(&res)
	if err != nil {
		return
	}
	return
}

// Gets host groups by names.
func (api *API) HostGroupsGetByNames(names []string) (res HostGroups, err error) {
	return api.HostGroupsGet(Params{"filter": map[string]interface{}{"name": names}})
}
