// filter-ec2-ips prints a comma-separated list of IP addresses for running EC2
// hosts, filtered by the provided ServerGroupName.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const Version = "0.1"

func main() {
	version := flag.Bool("version", false, "Print the version string")
	flag.Parse()
	if *version {
		fmt.Printf("filter-ec2-ips version %s\n", Version)
		os.Exit(2)
	}
	if flag.NArg() < 1 {
		os.Stderr.WriteString("usage: filter-ec2-ips server-group-name\n")
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 31*time.Second)
	defer cancel()
	name := flag.Arg(0)
	sess, err := session.NewSession(&aws.Config{
		MaxRetries: aws.Int(3),
	})
	if err != nil {
		log.Fatal(err)
	}
	client := ec2.New(sess)
	inp := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}},
			&ec2.Filter{Name: aws.String("tag:ServerGroupName"), Values: []*string{aws.String(name)}},
		},
	}
	req, output := client.DescribeInstancesRequest(inp)
	req.SetContext(ctx)
	if err := req.Send(); err != nil {
		log.Fatal(err)
	}
	ips := make([]string, 0)
	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			if instance.PrivateIpAddress == nil {
				log.Fatal("nil ip address for instance: " + *instance.InstanceId)
			}
			ips = append(ips, *instance.PrivateIpAddress)
		}
	}
	for i := range ips {
		ips[i] = ips[i] + ":9092"
	}
	fmt.Println(strings.Join(ips, ","))
}
