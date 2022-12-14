package zabbix_test

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	zabbix "github.com/XenoStar123/go-zabbix"
)

var (
	_host string
	_api  *zabbix.API
)

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	_host, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	_host += "-testing"

	if os.Getenv("TEST_ZABBIX_URL") == "" {
		log.Fatal("Set environment variables TEST_ZABBIX_URL (and optionally TEST_ZABBIX_USER and TEST_ZABBIX_PASSWORD)")
	}

	// set test host group name
	if os.Getenv("TEST_ZABBIX_HOST_GROUP") == "" {
		log.Fatal("Set environment variables TEST_ZABBIX_HOST_GROUP")
	}
}

func getHost() string {
	return _host
}

func getAPI(t *testing.T) *zabbix.API {
	if _api != nil {
		return _api
	}

	url, user, password := os.Getenv("TEST_ZABBIX_URL"), os.Getenv("TEST_ZABBIX_USER"), os.Getenv("TEST_ZABBIX_PASSWORD")
	_api = zabbix.NewAPI(url)
	_api.SetClient(http.DefaultClient)
	v := os.Getenv("TEST_ZABBIX_VERBOSE")
	if v != "" && v != "0" {
		_api.Logger = log.New(os.Stderr, "[zabbix] ", 0)
	}

	if user != "" {
		auth, err := _api.Login(user, password)
		if err != nil {
			t.Fatal(err)
		}
		if auth == "" {
			t.Fatal("Login failed")
		}
	}

	return _api
}

func TestBadCalls(t *testing.T) {
	api := getAPI(t)
	resp, err := api.Call("", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Error.Code != -32602 {
		t.Errorf("Expected code -32602, got %s", resp.Error)
	}
}

func TestVersion(t *testing.T) {
	api := getAPI(t)
	v, err := api.Version()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Zabbix version %s", v)
	if !regexp.MustCompile(`^\d\.\d\.\d+$`).MatchString(v) {
		t.Errorf("Unexpected version: %s", v)
	}
}

func ExampleAPI_Call() {
	api := zabbix.NewAPI("http://host/api_jsonrpc.php")
	api.Login("user", "password")
	resp, _ := api.Call("item.get", zabbix.Params{"itemids": "23970", "output": "extend"})
	log.Print(resp)
}
