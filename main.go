package main

import (
	"bufio"
	"flag"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	log "github.com/sirupsen/logrus"
	"net" // <-- 新增：处理 IP 地址和 CIDR 所需
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

// <-- 新增：实现 parseCIDRs 函数，将字符串列表解析为 net.IPNet 列表
func parseCIDRs(ipTxtList []string) []*net.IPNet {
	var ipList []*net.IPNet
	for _, ipTxt := range ipTxtList {
		// 使用 net.ParseCIDR 解析 IPv4 或 IPv6 地址/CIDR
		_, ipNet, err := net.ParseCIDR(ipTxt)
		if err != nil {
			// 如果解析失败（可能是单独的IP地址而不是CIDR），尝试解析为单个IP
			ip := net.ParseIP(ipTxt)
			if ip != nil {
				// 对于单个 IP，创建一个 /32 (IPv4) 或 /128 (IPv6) 的网络
				maskLen := net.IPv4len * 8 // 32
				if ip.To4() == nil { // 是 IPv6 地址
					maskLen = net.IPv6len * 8 // 128
				}
				
				// net.CIDRMask 创建子网掩码
				ipNet = &net.IPNet{
					IP: ip,
					Mask: net.CIDRMask(maskLen, maskLen),
				}
			} else {
				// 真正的错误，跳过该行
				log.Warnf("Skipping invalid IP/CIDR entry: %s, error: %v", ipTxt, err)
				continue
			}
		}
		ipList = append(ipList, ipNet)
	}
	return ipList
}
// --> end parseCIDRs

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
