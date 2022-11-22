package zabbix_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestHosts(t *testing.T) {
	groupName := os.Getenv("TEST_ZABBIX_HOST_GROUP")

	api := getAPI(t)

	groups, err := api.HostGroupsGetByNames([]string{groupName})
	if err != nil {
		t.Fatal(err)
	}

	ids := make([]string, 0)
	for _, g := range groups {
		ids = append(ids, g.GroupId)
	}

	hosts, err := api.HostsGetByHostGroupIds(ids)
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range hosts {
		b, err := json.Marshal(h)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("host: %s", b)
	}
}
