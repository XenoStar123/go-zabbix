package zabbix_test

import (
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
	m := items.ByHostId()
	for hostId, items := range m {
		t.Logf("hostId: %s item len: %v", hostId, len(items))
	}
}
