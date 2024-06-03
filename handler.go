package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/xiaoyi-byte/student-demo/dal/mysql"
	demo "github.com/xiaoyi-byte/student-demo/kitex_gen/demo"
	"github.com/xiaoyi-byte/student-demo/model"
	"strings"
	"sync"
)

var students sync.Map

// StudentServiceImpl implements the last service interface defined in the IDL.
type StudentServiceImpl struct{}

// Register implements the StudentServiceImpl interface.
func (s *StudentServiceImpl) Register(ctx context.Context, student *demo.Student) (resp *demo.RegisterResp, err error) {
	klog.Infof("Register req: %v", student)
	req := &demo.QueryReq{Id: student.Id}
	stu, _ := s.Query(ctx, req)
	if stu == nil {
		students.Store(student.Id, student)
		emails := ""
		if len(student.GetEmail()) > 0 {
			emails = student.GetEmail()[0]
			for i := 1; i < len(student.GetEmail()); i++ {
				emails = emails + "," + student.GetEmail()[1]
			}
		}
		user := &model.User{
			Id:             student.GetId(),
			Name:           student.GetName(),
			CollegeName:    student.GetCollege().GetName(),
			CollegeAddress: student.GetCollege().GetAddress(),
			Emails:         emails,
			Sex:            student.GetSex(),
		}
		// 写入数据库
		err := mysql.CreateUser(user)
		if err != nil {
			return nil, err
		}
	}
	resp = &demo.RegisterResp{
		Success: true,
		Message: "",
	}
	return
}

// Query implements the StudentServiceImpl interface.
func (s *StudentServiceImpl) Query(ctx context.Context, req *demo.QueryReq) (resp *demo.Student, err error) {
	klog.Infof("query req: %v", req)
	if stu, ok := students.Load(req.Id); ok {
		resp = stu.(*demo.Student)
		return
	}
	// query in db
	user, err := mysql.QueryUser(req.Id)
	if err != nil {
		return nil, err
	}
	resp = &demo.Student{
		Id:   user.Id,
		Name: user.Name,
		College: &demo.College{
			Name:    user.CollegeName,
			Address: user.CollegeAddress,
		},
		Email: strings.Split(user.Emails, ","),
		Sex:   user.Sex,
	}
	return
}
