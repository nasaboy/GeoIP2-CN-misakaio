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
	srcFile      string
	dstFile      string
	databaseType string
	description  string
	// 仅保留 country -> iso_code: "CN"
	cnRecord = mmdbtype.Map{
		"country": mmdbtype.Map{
			"iso_code": mmdbtype.String("CN"),
		},
	}
)

func init() {
	flag.StringVar(&srcFile, "s", "ipip_cn.txt", "specify source ip list file")
	flag.StringVar(&dstFile, "d", "Country.mmdb", "specify destination mmdb file")
	flag.StringVar(&databaseType, "t", "GeoLite2-Country", "specify MaxMind database type")
	flag.StringVar(&description, "description", "Customized GeoLite2 Country database", "specify a description for mmdb metadata")
	flag.Parse()
}

// parseCIDRs 定义见 ip2cidr.go

func main() {
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: databaseType,
			RecordSize:   24,
			Description: map[string]string{
				"en": description,
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
