package main

import (
	"bufio"
	"flag"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	log "github.com/sirupsen/logrus"
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
	// ⭐️ 修改点 1: 精简 cnRecord，只保留 iso_code
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
		},
	)
	if err != nil {
		log.Fatalf("fail to new writer %v\n", err)
	}

	// ⭐️ 修改点 2: 增加 Description 和 Languages
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

	for scanner.Scan
