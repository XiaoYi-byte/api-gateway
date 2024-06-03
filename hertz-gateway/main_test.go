package main

import (
	"bytes"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"github.com/xiaoyi-byte/student-demo/kitex_gen/demo"
	"io"
	"net/http"
	"testing"
	"time"
)

const gatewayURL = "http://127.0.0.1:8888/gateway/student"

var httpCli = &http.Client{Timeout: 3 * time.Second}

type reqParam struct {
	Method    string `json:"method"`
	BizParams string `json:"biz_params"`
}

func TestStudentService(t *testing.T) {
	for i := 1; i <= 100; i++ {
		newStu := genStudent(i)
		resp, err := request("register", newStu)
		Assert(t, err == nil, err)
		Assert(t, resp["message"] == "ok")
		jsonData, err := json.Marshal(resp["data"])
		Assert(t, err == nil, err)

		// 将JSON转换为demo.RegisterResp类型
		var registerResp demo.RegisterResp
		err = json.Unmarshal(jsonData, &registerResp)
		Assert(t, err == nil, err)
		Assert(t, registerResp.Success)

		resp, err = request("query", newStu)
		Assert(t, err == nil, err)
		Assert(t, resp["message"] == "ok")

		jsonData, err = json.Marshal(resp["data"])
		Assert(t, err == nil, err)
		var stu demo.Student
		err = json.Unmarshal(jsonData, &stu)
		Assert(t, stu.Id == newStu.Id)
		Assert(t, stu.Name == newStu.Name)
		Assert(t, stu.Email[0] == newStu.Email[0])
		Assert(t, stu.College.Name == newStu.College.Name)
	}
}

func BenchmarkStudentService(b *testing.B) {
	for i := 1; i < b.N; i++ {
		newStu := genStudent(i)
		resp, err := request("register", newStu)
		Assert(b, err == nil, err)
		Assert(b, resp["message"] == "ok")
		// 将resp["data"]转成RegisterResp类型
		jsonData, err := json.Marshal(resp["data"])
		Assert(b, err == nil, err)
		var registerResp demo.RegisterResp
		err = json.Unmarshal(jsonData, &registerResp)
		Assert(b, err == nil, err)
		Assert(b, registerResp.Success)

		resp, err = request("query", newStu)
		Assert(b, err == nil, err)
		Assert(b, resp["message"] == "ok")

		jsonData, err = json.Marshal(resp["data"])
		Assert(b, err == nil, err)
		var stu demo.Student
		err = json.Unmarshal(jsonData, &stu)
		Assert(b, stu.Id == newStu.Id)
		Assert(b, stu.Name == newStu.Name)
		Assert(b, stu.Email[0] == newStu.Email[0])
		Assert(b, stu.College.Name == newStu.College.Name)
	}
}

func request(method string, bizParam any) (rResp map[string]interface{}, err error) {
	bizParamBody, err := json.Marshal(bizParam)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: err=%v", err)
	}
	reqBody, err := json.Marshal(&reqParam{
		Method:    method,
		BizParams: string(bizParamBody),
	})
	if err != nil {
		return nil, fmt.Errorf("marshal bizParam failed: err=%v", err)
	}
	var resp *http.Response
	req, err := http.NewRequest(http.MethodPost, gatewayURL, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = httpCli.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}
	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return
	}

	if err = json.Unmarshal(body, &rResp); err != nil {
		return
	}
	return
}

func genStudent(id int) *demo.Student {
	return &demo.Student{
		Id:   int32(id),
		Name: fmt.Sprintf("student-%d", id),
		College: &demo.College{
			Name:    "",
			Address: "",
		},
		Email: []string{fmt.Sprintf("student-%d@nju.com", id)},
	}
}

// Assert asserts cond is true, otherwise fails the test.
func Assert(t testingTB, cond bool, val ...interface{}) {
	t.Helper()
	if !cond {
		if len(val) > 0 {
			val = append([]interface{}{"assertion failed:"}, val...)
			t.Fatal(val...)
		} else {
			t.Fatal("assertion failed")
		}
	}
}

// testingTB is a subset of common methods between *testing.T and *testing.B.
type testingTB interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
}
