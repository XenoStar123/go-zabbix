package zabbix_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestItems(t *testing.T) {
	groupName := os.Getenv("TEST_ZABBIX_HOST_GROUP")

	api := getAPI(t)

	groups, err := api.HostGroupsGetByNames([]string{groupName})
	if err != nil {
		t.Fatal(err)
	}

	groupIds := make([]string, 0)
	for _, g := range groups {
		groupIds = append(groupIds, g.GroupId)
	}

	hosts, err := api.HostsGetByHostGroupIds(groupIds)
	if err != nil {
		t.Fatal(err)
	}

	hostIds := make([]string, 0)
	for _, h := range hosts {
		hostIds = append(hostIds, h.HostId)
	}

	items, err := api.ItemsGetByHostIds(hostIds)
	if err != nil {
		t.Fatal(err)
	}
	for _, i := range items {
		b, err := json.Marshal(i)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("item: %s", b)
	}
}
