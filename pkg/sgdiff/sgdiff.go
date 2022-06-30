package sgdiff

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/urfave/cli/v2"
)

type SecurityGroup struct {
	Description   *string       `json:"description"`
	GroupId       *string       `json:"group_id"`
	GroupName     *string       `json:"group_name"`
	IpPermissions *IpPermission `json:"ip_permissions"`
}

type IpPermission struct {
	FromPort   *int64   `json:"from_port"`
	IpProtocol *string  `json:"ip_protocol"`
	IpRanges   *IpRange `json:"ip_range"`
}

type IpRange struct {
	CidrIp *string `json:"cidr_ip"`
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "show",
		Usage: "show sg list",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Security Group ID",
				Required: true,
			},
		},
		Action: action,
	}
}

func action(c *cli.Context) error {
	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config:            aws.Config{Region: aws.String("ap-northeast-1")},
			Profile:           "admin",
			SharedConfigState: session.SharedConfigEnable,
		},
	)
	if err != nil {
		return err
	}

	svc := ec2.New(sess)
	input := &ec2.DescribeSecurityGroupsInput{
		GroupIds: []*string{
			aws.String(c.String("id")),
		},
	}

	d, err := svc.DescribeSecurityGroups(input)
	if err != nil {
		return err
	}

	// ips := d.SecurityGroups[0].IpPermissions[0]
	ips := d.SecurityGroups[0]
	sg := SecurityGroup{
		Description: ips.Description,
		GroupId:     ips.GroupId,
		GroupName:   ips.GroupName,
	}
	ipp := IpPermission{
		FromPort:   ips.IpPermissions[0].FromPort,
		IpProtocol: ips.IpPermissions[0].IpProtocol,
	}
	ipr := IpRange{
		CidrIp: ips.IpPermissions[0].IpRanges[0].CidrIp,
	}
	sg.IpPermissions = &ipp
	ipp.IpRanges = &ipr

	byteArray, err := json.Marshal(sg)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(string(byteArray))

	hoge := "{\"description\":\"default VPC security group\",\"group_id\":\"sg-xxxx\",\"group_name\":\"default\",\"ip_permissions\":{\"from_port\":80,\"ip_protocol\":\"tcp\",\"ip_range\":{\"cidr_ip\":\"123.123.123.12/32\"}}}"
	fmt.Println(lineDiff(string(byteArray), hoge))
	// fmt.Println(d.SecurityGroups[0])
	// to := reflect.TypeOf(*ips)
	// vo := reflect.ValueOf(*ips)

	// remap := make(map[string]interface{})
	// for i := 0; i < to.NumField(); i++ {
	// 	f := to.Field(i)
	// 	key := f.Name
	// 	if key != "_" {
	// 		v := vo.FieldByName(key).Interface()
	// 		fmt.Println(reflect.TypeOf(v))
	// 		switch value := v.(type) {
	// 		case *string:
	// 			remap[key] = *value
	// 		case *int64:
	// 			remap[key] = *value
	// 			// case []*ec2.IpRange:
	// 			// 	remap[key] = value
	// 		}
	// 	}
	// }
	// fmt.Println(remap)

	return err
}

func lineDiff(src1, src2 string) string {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(src1, src2)
	diffs := dmp.DiffMain(a, b, true)
	result := dmp.DiffCharsToLines(diffs, c)

	return dmp.DiffPrettyText(result)
}
