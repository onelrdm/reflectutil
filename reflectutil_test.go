package reflectutil

import (
	"fmt"
	"testing"
	"time"
)

type BaseModel struct {
	ID        uint
	CreatedAt *time.Time `col:"Application Time"`
}

type OrderInfo struct {
	BaseModel
	OrderNo               string     `col:"Loan ID"`
	CustomerFullName      string     `col:"Name"`
	CustomerPhone         string     `col:"Mobile number"`
	IDType                string     `col:"ID Type"`
	IDNumber              string     `col:"ID No"`
	ProductName           string     `col:"Product"`
	OperatorName          string     `col:"Reviewer"`
	ApplicationStatus     string     `col:"Review status"`
	ApplicationFinishTime *time.Time `col:"Review Time"`
}

type BaseHttpResp struct {
	Success bool
	Data    interface{}
}

type OrderListRespData struct {
	Orders     []*OrderInfo
	TotalCount uint
	TotalPage  uint
}

func TestReflect(t *testing.T) {
	now := time.Now()
	orders := []*OrderInfo{
		{
			BaseModel:             BaseModel{CreatedAt: &now},
			OrderNo:               "OrderNo1",
			CustomerFullName:      "CustomerFullName1",
			CustomerPhone:         "CustomerPhone1",
			IDType:                "IDType1",
			IDNumber:              "IDNumber1",
			ProductName:           "ProductName1",
			OperatorName:          "OperatorName1",
			ApplicationStatus:     "ApplicationStatus1",
			ApplicationFinishTime: &now,
		},
		{
			BaseModel:             BaseModel{CreatedAt: &now},
			OrderNo:               "OrderNo2",
			CustomerFullName:      "CustomerFullName2",
			CustomerPhone:         "CustomerPhone2",
			IDType:                "IDType2",
			IDNumber:              "IDNumber2",
			ProductName:           "ProductName2",
			OperatorName:          "OperatorName2",
			ApplicationStatus:     "ApplicationStatus2",
			ApplicationFinishTime: &now,
		},
	}
	data := &OrderListRespData{
		Orders:     orders,
		TotalCount: 0,
		TotalPage:  0,
	}
	resp := &BaseHttpResp{
		Success: true,
		Data:    data,
	}
	fmt.Printf("%+v", resp)
}

func BenchmarkReflect(b *testing.B) {
	for i := 0; i < b.N; i++ {
	}
}
