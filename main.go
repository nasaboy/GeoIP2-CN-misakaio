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
		//"country": mmdbtype.Map{
		//"continent": mmdbtype.Map{
		//"code": mmdbtype.String("AS"),
		//"geoname_id": mmdbtype.Uint32(6255147),
		//"names":mmdbtype.Map{"de":mmdbtype.String("Asien"),
		//"en":mmdbtype.String("Asia"),
		//"es":mmdbtype.String("Asia"),
		//"fr":mmdbtype.String("Asie"),
		//"ja":mmdbtype.String("アジア"),
		//"pt-BR":mmdbtype.String("Ásia"),
		//"ru":mmdbtype.String("Азия"),
		//"zh-CN":mmdbtype.String("亚洲")},
		//},
		//"country": mmdbtype.Map{
		//"geoname_id":mmdbtype.Uint32(1814991),
		//"is_in_european_union":mmdbtype.Bool(false),
		//"iso_code":mmdbtype.String("CN"),
		//"names":mmdbtype.Map{
		//"de":mmdbtype.String("China"),
		//"en":mmdbtype.String("China"),
		//"es":mmdbtype.String("China"),
		//"fr":mmdbtype.String("Chine"),
		//"ja":mmdbtype.String("中国"),
		//"pt-BR":mmdbtype.String("China"),
		//"ru":mmdbtype.String("Китай"),
		//"zh-CN":mmdbtype.String("中国"),
		//},
		//},
		//"registered_country": mmdbtype.Map{
		//"geoname_id":mmdbtype.Uint32(1814991),
		//"is_in_european_union":mmdbtype.Bool(false),
		//"iso_code":mmdbtype.String("CN"),
		//"names":mmdbtype.Map{
		//"de":mmdbtype.String("China"),
		//"en":mmdbtype.String("China"),
		//"es":mmdbtype.String("China"),
		//"fr":mmdbtype.String("Chine"),
		//"ja":mmdbtype.String("中国"),
		//"pt-BR":mmdbtype.String("China"),
		//"ru":mmdbtype.String("Китай"),
		//"zh-CN":mmdbtype.String("中国"),
		//},
		//},
		//"traits": mmdbtype.Map{
		//"is_anonymous_proxy": mmdbtype.Bool(false),
		//"is_satellite_provider":mmdbtype.Bool(false),
		//},
		//},

var (
	srcFile string
	dstFile string
	databaseType string
	// 已修改：只保留 iso_code
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

// ⭐️ 新增：parseCIDRs 函数的实现，解决编译错误
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

	// 已新增：设置 Description 和 Languages
	writer.Metadata.Description = map[string]string{
		"en":    "IP-to-Country Database (CN only)",
		"zh-CN": "IP到国家/地区数据库 (仅中国)",
	}
	writer.Metadata.Languages = []string{"en", "zh-CN"}

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



}
