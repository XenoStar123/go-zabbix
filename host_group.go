package zabbix

import "fmt"

// https://www.zabbix.com/documentation/2.2/manual/appendix/api/hostgroup/definitions
type HostGroup struct {
	GroupId string `json:"groupid"`
	Name    string `json:"name"`
}

type HostGroups []HostGroup

// Wrapper for hostgroup.get: https://www.zabbix.com/documentation/3.4/en/manual/api/reference/hostgroup/get
func (api *API) HostGroupsGet(params Params) (HostGroups, error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	resp, err := api.CallWithError("hostgroup.get", params)
	if err != nil {
		return nil, fmt.Errorf("api.CallWithError: %v", err)
	}

	res := HostGroups{}
	err = resp.Bind(&res)
	if err != nil {
		return nil, fmt.Errorf("resp.Bind: %v", err)
	}
	return res, nil
}

// Gets host groups by names.
func (api *API) HostGroupsGetByNames(names []string) (HostGroups, error) {
	return api.HostGroupsGet(Params{"filter": map[string]interface{}{"name": names}})
}
