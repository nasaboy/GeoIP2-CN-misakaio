package main

import (
	"bufio"
	"flag"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	srcFile string
	dstFile string
	databaseType string
	// 只保留了 iso_code
	cnRecord = mmdbtype.Map{
		"country": mmdbtype.Map{
			"iso_code":             mmdbtype.String("CN"),
		},
	}
)

func init()  {
	flag.StringVar(&srcFile, "s", "ipip_cn.txt", "specify source ip list file")
	flag.StringVar(&dstFile, "d", "Country.mmdb", "specify destination mmdb file")
	flag.StringVar(&databaseType,"t", "GeoIP2-Country", "specify MaxMind database type")
	flag.Parse()
}

func main()  {
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: databaseType,
			RecordSize:   24,
			// 在 Metadata 中增加 Description
			Metadata: mmdbtype.Map{
				"description": mmdbtype.Map{
					"en": mmdbtype.String("Custom MMDB for China ISO Code"),
					"zh-CN": mmdbtype.String("自定义中国ISO代码MMDB"),
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("fail to new writer %v\n", err)
	}

	var ipTxtList []string
	fh, err := os.Open(srcFile)
	if err != nil {
		log.Fatalf("fail to open %s\n", err)
	}
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ipTxtList = append(ipTxtList, scanner.Text())
	}

	// 注意：`parseCIDRs` 函数未在提供的代码片段中定义。
	// 假设它是一个存在的函数，用于将文本行转换为可插入的 IP/CIDR 列表。
	ipList := parseCIDRs(ipTxtList) 
	for _, ip := range ipList {
		err = writer.Insert(ip, cnRecord)
		if err != nil {
			log.Fatalf("fail to insert to writer %v\n", err)
		}
	}

	outFh, err := os.Create(dstFile)
	if err != nil {
		log.Fatalf("fail to create output file %v\n", err)
	}

	_, err = writer.WriteTo(outFh)
	if err != nil {
		log.Fatalf("fail to write to file %v\n", err)
	}

}
