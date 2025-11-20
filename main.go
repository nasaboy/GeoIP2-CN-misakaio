package main

import (
	"bufio"
	"flag"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	log "github.com/sirupsen/logrus"
	"net" // ⭐️ 新增：用于 IP/CIDR 解析
	"os"
)

var (
	srcFile string
	dstFile string
	databaseType string
	// 精简后的 cnRecord，只保留 iso_code: "CN"
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

// 示例：实现 parseCIDRs 函数，将字符串列表解析为 net.IPNet 列表
func parseCIDRs(ipTxtList []string) []*net.IPNet {
	var ipList []*net.IPNet
	for _, ipTxt := range ipTxtList {
		if ipTxt == "" {
			continue // 跳过空行
		}
		// 尝试解析为 CIDR (如 1.1.1.0/24)
		_, ipNet, err := net.ParseCIDR(ipTxt)
		if err == nil {
			ipList = append(ipList, ipNet)
			continue
		}

		// 如果解析 CIDR 失败，尝试解析为单个 IP (如 8.8.8.8)
		ip := net.ParseIP(ipTxt)
		if ip != nil {
			// 对于单个 IP，我们创建一个 /32 (IPv4) 或 /128 (IPv6) 的 IPNet
			var mask net.IPMask
			if ip.To4() != nil {
				mask = net.CIDRMask(32, 32)
			} else {
				mask = net.CIDRMask(128, 128)
			}
			ipNet = &net.IPNet{IP: ip, Mask: mask}
			ipList = append(ipList, ipNet)
			continue
		}
		
		log.Warnf("Skipping invalid IP or CIDR: %s", ipTxt)
	}
	return ipList
}

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

	// ⭐️ 设置数据库描述 (Description) 和语言 (Languages)
	writer.Metadata.Description = map[string]string{
		"en":    "Customized GeoLite2 Country database",
	}
	writer.Metadata.Languages = []string{"en", "zh-CN"}

	var ipTxtList []string
	fh, err := os.Open(srcFile)
	if err != nil {
		// 修改日志级别以避免在找不到文件时直接退出
		log.Fatalf("fail to open source file %s: %v\n", srcFile, err)
	}
	defer fh.Close() // 确保文件句柄被关闭

	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ipTxtList = append(ipTxtList, scanner.Text())
	}

	ipList := parseCIDRs(ipTxtList)
	if len(ipList) == 0 {
		log.Fatalf("No valid IPs or CIDRs found in %s. Exiting.\n", srcFile)
	}
	
	log.Infof("Inserting %d IP/CIDR records into the database...", len(ipList))

	for _, ip := range ipList {
		// ip 是 *net.IPNet，可以直接插入
		err = writer.Insert(ip, cnRecord)
		if err != nil {
			log.Fatalf("fail to insert %s to writer: %v\n", ip.String(), err)
		}
	}

	outFh, err := os.Create(dstFile)
	if err != nil {
		log.Fatalf("fail to create output file %v\n", err)
	}
	defer outFh.Close() // 确保输出文件句柄被关闭

	log.Infof("Writing database to %s...", dstFile)
	_, err = writer.WriteTo(outFh)
	if err != nil {
		log.Fatalf("fail to write to file %v\n", err)
	}
	
	log.Infof("Successfully created %s!", dstFile)

}
