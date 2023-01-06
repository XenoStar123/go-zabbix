package zabbix

import "fmt"

// ValueType value_type
// (required)	integer	Type of information of the item.
// Possible values:
// 0 - numeric float;
// 1 - character;
// 2 - log;
// 3 - numeric unsigned;
// 4 - text.
type (
	ValueType string
)

const (
	Float     ValueType = "0"
	Character ValueType = "1"
	Log       ValueType = "2"
	Unsigned  ValueType = "3"
	Text      ValueType = "4"
)

// https://www.zabbix.com/documentation/3.4/en/manual/api/reference/item/object
type Item struct {
	ItemId    string    `json:"itemid"`
	HostId    string    `json:"hostid"`
	Key       string    `json:"key_"`
	Name      string    `json:"name"`
	ValueType ValueType `json:"value_type"`
	LastClock string    `json:"lastclock"`
	LastValue string    `json:"lastvalue"`
}

type Items []Item

// Converts slice to map by hostId.
func (items Items) ByHostId() map[string]Items {
	res := make(map[string]Items, 0)
	for _, i := range items {
		if _, ok := res[i.HostId]; !ok {
			res[i.HostId] = make(Items, 0)
		}
		res[i.HostId] = append(res[i.HostId], i)
	}
	return res
}

// Wrapper for item.get https://www.zabbix.com/documentation/3.4/en/manual/api/reference/item/get
func (api *API) ItemsGet(params Params) (Items, error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	resp, err := api.CallWithError("item.get", params)
	if err != nil {
		return nil, fmt.Errorf("api.CallWithError: %v", err)
	}

	res := Items{}
	err = resp.Bind(&res)
	if err != nil {
		return nil, fmt.Errorf("resp.Bind: %v", err)
	}
	return res, nil
}

// Gets items by host Ids.
func (api *API) ItemsGetByHostIds(ids []string) (Items, error) {
	return api.ItemsGet(
		Params{
			"hostids": ids,
			"output": []string{
				"itemid",
				"hostid",
				"key_",
				"name",
				"lastvalue",
				"lastclock",
				"value_type",
			}})
}
