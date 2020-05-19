package cacli

import (
	"encoding/json"
	"fmt"
	"k8s-ca-websocket/cautils"
)

// WrokloadTemplateWordpressMock worpres mock
func WrokloadTemplateWordpressMock() *cautils.WorkloadTemplate {
	wordpressWT := `{
		"kind": "deployment",
		"name": "wordpress",
		"cluster": "system-test",
		"namespace": "default",
		"groupingLevel0": "system-test",
		"groupingLevel1": "native",
		"wlid": "wlid://cluster-system-test/namespace-default/deployment-wordpress",
		"metainfo": {
			"creationDate": "2019-12-19 11:50:16.852792971 +0000 UTC",
			"lastEdited": "2019-12-19 11:50:16.93203732 +0000 UTC",
			"workloadKind": "",
			"instances": {
				"active": [],
				"inactive": []
			},
			"categories": null
		},
		"autoAccessTokenUpdate": true,
		"containers": [
			{
				"name": "mysql",
				"os": "undefined",
				"architecture": "64bit",
				"imageHash": "sha256:5779c71a4730da36f013a23a437b5831198e68e634575f487d37a0639470e3a8",
				"imageTag": "mysql:5.7",
				"enableVisiblity": [
					{
						"ALL": true
					}
				]
			},
			{
				"name": "wordpress",
				"os": "debian",
				"architecture": "64bit",
				"imageHash": "docker-pullable://wordpress@sha256:6216f64ab88fc51d311e38c7f69ca3f9aaba621492b4f1fa93ddf63093768845",
				"imageTag": "wordpress:4.8",
				"enableVisiblity": [
					{
						"ALL": true
					}
				],
				"signingProfileName": "wordpress:4.8"
			}
		]
	}`
	wt := &cautils.WorkloadTemplate{}
	if err := json.Unmarshal([]byte(wordpressWT), wt); err != nil {
		fmt.Println(err)
		return nil
	}
	return wt
}

// CacliMock commands
type CacliMock struct {
}

// NewCacliMock -
func NewCacliMock() *CacliMock {
	return &CacliMock{}
}

// Login -
func (caclim *CacliMock) Login(globalLoginCredentials cautils.CredStruct) error {
	return nil
}

// Get wordpress wlid
func (caclim *CacliMock) Get(wlid string) (cautils.WorkloadTemplate, error) {
	return *WrokloadTemplateWordpressMock(), nil
}

// Sign wordpress wlid
func (caclim *CacliMock) Sign(wlid, user, password string) error {
	return nil
}
