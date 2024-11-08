package collector

import "testing"

var (
	MDTTotalSpaceCommandOutput = `
osd-ldiskfs.sjtu-MDT0000.kbytestotal=11714622144
osd-ldiskfs.sjtu-MDT0001.kbytestotal=11714622144
`
	MDTFreeSpaceCommandOutput = `
osd-ldiskfs.sjtu-MDT0000.kbytesfree=10401788800
osd-ldiskfs.sjtu-MDT0001.kbytesfree=10758664440
`
)

/*
TODO:
1. 优化collector的结构：现在的collector有一个metric全局变量，用来存储所有的指标，这种写法不利于编写测试和维护
2. 重写单元测试：现在单元测试的部分为了能兼容现在的全局变量写法，做了一些workaround。在后面重构collector部分之后，单元测试部分也要重构。
*/
func TestParseMDTInfo(t *testing.T) {
	// copy MDTList manually
	var mdsMetricCopy MDSSpaceMetric
	mdsMetricCopy.MDTList = make(map[string]MDTMetric)
	for k, v := range mdsMetric.MDTList {
		mdsMetricCopy.MDTList[k] = v
	}
	clear(mdsMetric.MDTList)
	// test part
	err := parseMDTInfo(MDTTotalSpaceCommandOutput)
	if err != nil {
		t.Errorf("parseMDTInfo() throws error: %s", err.Error())
		return
	}
	tests := []struct {
		input    string
		expected int64
	}{
		{"sjtu-MDT0000", 11714622144},
		{"sjtu-MDT0001", 11714622144},
	}
	for _, tt := range tests {
		if got := mdsMetric.MDTList[tt.input].KBTotal; got != tt.expected {
			t.Errorf("ParseMDTInfo test fail, input=%s, want=%d", tt.input, tt.expected)
		}
	}
	// restore original metric
	for k, v := range mdsMetricCopy.MDTList {
		mdsMetric.MDTList[k] = v
	}
}
