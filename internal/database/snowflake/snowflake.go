package snowflake

import (
	"restdoc/third/sonyflake"
	"restdoc/third/sonyflake/awsutil"
)

var Sf *sonyflake.Sonyflake

func Init() {
	var st sonyflake.Settings
	st.MachineID = awsutil.AmazonEC2MachineID
	Sf = sonyflake.NewSonyflake(st)
	if Sf == nil {
		panic("sonyflake not created")
	}
}
