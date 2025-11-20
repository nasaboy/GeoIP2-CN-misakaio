package main

import (
	"bufio"
	"flag"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	log "github.com/sirupsen/logrus"
	"os"
)

// 移除所有注释掉的字段
// "continent": mmdbtype.Map{...}
// "registered_country": mmdbtype.Map{...}
// "traits": mmdbtype.Map{...}

var (
	srcFile string
	dstFile string
	databaseType string
	// **修改 cnRecord**：
	// 仅保留 country -> iso_code: "CN"
	cnRecord = mmdbtype.Map{
		"country": mmdbtype.Map{
			"iso_code":             mmdbtype.String("CN"),
		},
	}
)

func init()  {
	flag.StringVar(&srcFile, "s", "ipip_cn.txt", "specify source ip list file")
	flag.StringVar(&dstFile, "d", "Country.mmdb", "specify destination mmdb file")
	flag.StringVar(&databaseType,"t", "GeoLite2-Country", "specify MaxMind database type")
	flag.Parse()
}

// 注意：你的代码中缺少 parseCIDRs 函数的定义，
// 如果在同一目录下有其他文件定义了它，请确保它在编译时可用。
// 否则，你需要在这里或另一个文件中定义它，例如：
// func parseCIDRs(ipTxtList []string) []net.IPNet { ... }

func main()  {
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: databaseType,
			RecordSize:   24,
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

	// 假设 parseCIDRs 函数已定义并能将字符串列表转换为 *net.IPNet 列表
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
