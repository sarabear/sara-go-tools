package main

import (
	"fmt"

	dns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

var lastIPv6Map = make(map[string]string) // domain => ipv6

func newDNSClient(config *Config) (*dns.Client, error) {
	conf := &openapi.Config{
		AccessKeyId:     tea.String(config.AccessKeyId),
		AccessKeySecret: tea.String(config.AccessKeySecret),
		RegionId:        tea.String(config.RegionId),
	}
	client, err := dns.NewClient(conf)
	return client, err
}

func runOnce(client *dns.Client, config *Config, prefix string) {
	for _, d := range config.Domains {
		fullIPv6 := fmt.Sprintf("%s:%s", prefix, d.Suffix)
		lastIP := lastIPv6Map[d.RR+"."+d.Domain]

		if lastIP == fullIPv6 {
			logInfo("域名 %s.%s 对应的IPv6地址未变化，跳过更新", d.RR, d.Domain)
			continue
		}

		currentIP, err := getCurrentIPv6(client, d.Domain, d.RR)
		if err != nil {
			logError("获取域名 %s.%s 当前IPv6失败: %v", d.RR, d.Domain, err)
			continue
		}

		if currentIP == fullIPv6 {
			logInfo("域名 %s.%s 的IPv6地址已匹配，无需更新", d.RR, d.Domain)
			lastIPv6Map[d.RR+"."+d.Domain] = fullIPv6
			continue
		}

		err = updateIPv6Record(client, d.Domain, d.RR, fullIPv6)
		if err != nil {
			logError("更新域名 %s.%s 失败: %v", d.RR, d.Domain, err)
		} else {
			logInfo("✅ 域名 %s.%s 更新成功，新IPv6地址: %s", d.RR, d.Domain, fullIPv6)
			lastIPv6Map[d.RR+"."+d.Domain] = fullIPv6
		}
	}
}

func getCurrentIPv6(client *dns.Client, domain, rr string) (string, error) {
	request := &dns.DescribeDomainRecordsRequest{
		DomainName: tea.String(domain),
		RRKeyWord:  tea.String(rr),
		Type:       tea.String("AAAA"),
	}
	resp, err := client.DescribeDomainRecords(request)
	if err != nil {
		return "", err
	}
	if len(resp.Body.DomainRecords.Record) > 0 {
		return *resp.Body.DomainRecords.Record[0].Value, nil
	}
	return "", nil
}

func updateIPv6Record(client *dns.Client, domain, rr, ipv6 string) error {
	recordID, err := getRecordID(client, domain, rr)
	if err != nil {
		return err
	}
	if recordID == "" {
		return fmt.Errorf("未找到记录")
	}

	request := &dns.UpdateDomainRecordRequest{
		RecordId: tea.String(recordID),
		RR:       tea.String(rr),
		Type:     tea.String("AAAA"),
		Value:    tea.String(ipv6),
	}
	_, err = client.UpdateDomainRecord(request)
	return err
}

func getRecordID(client *dns.Client, domain, rr string) (string, error) {
	request := &dns.DescribeDomainRecordsRequest{
		DomainName: tea.String(domain),
		RRKeyWord:  tea.String(rr),
		Type:       tea.String("AAAA"),
	}
	resp, err := client.DescribeDomainRecords(request)
	if err != nil {
		return "", err
	}
	if len(resp.Body.DomainRecords.Record) > 0 {
		return *resp.Body.DomainRecords.Record[0].RecordId, nil
	}
	return "", nil
}
