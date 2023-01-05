package zabbix_test

import (
	"encoding/json"
	"os"
	"testing"
)

func TestHostGroups(t *testing.T) {
	groupName := os.Getenv("TEST_ZABBIX_HOST_GROUP")

	api := getAPI(t)

	groups, err := api.HostGroupsGetByNames([]string{groupName})
	if err != nil {
		t.Fatal(err)
	}
	for _, g := range groups {
		b, err := json.Marshal(g)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("hostGroup: %s", b)
	}
}
