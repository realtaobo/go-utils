package common

import (
	"testing"
)

// 单元测试结构体
type ComponentTest struct {
	User   string
	Action string
}

func TestJsonToInterface(t *testing.T) {
	varTest2 := make(map[string]string)
	type args struct {
		jsonStr string
		info    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestJsonToInterface1",
			args: args{
				jsonStr: `{"user":"tryturned", "action":"upsert"}`,
				info:    &ComponentTest{},
			},
			wantErr: false,
		},
		{
			name: "TestJsonToInterface2",
			args: args{
				jsonStr: `{"user":"tryturned", "action":"insert"}`,
				info:    &varTest2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := JsonToInterface(tt.args.jsonStr, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("JsonToInterface() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log(tt.args.info)
			}
		})
	}
}

func TestYamlToInterface(t *testing.T) {
	varTest2 := make(map[string]string)
	type args struct {
		yamlStr string
		info    interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestYamlToInterface1",
			args: args{
				yamlStr: "user: tryturned\naction: insert",
				info:    &ComponentTest{},
			},
			wantErr: false,
		},
		{
			name: "TestYamlToInterface2",
			args: args{
				yamlStr: "user: tryturned\naction: update",
				info:    &varTest2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := YamlToInterface(tt.args.yamlStr, tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("YamlToInterface() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log(tt.args.info)
			}
		})
	}
}

func TestCalcMonth(t *testing.T) {
	t.Log(CalcMonth(202101, -1, -2)) // 201911
	t.Log(CalcMonth(202101, -1, 0))  // 202001
	t.Log(CalcMonth(202101, 0, -2))  // 202011
}

func TestCalcSubMonth(t *testing.T) {
	t.Log(CalcSubMonth(202101, 202101)) // 1
	t.Log(CalcSubMonth(202101, 202104)) // 4
	t.Log(CalcSubMonth(202008, 202101)) // 6
}
