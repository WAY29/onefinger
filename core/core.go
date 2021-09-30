package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/WAY29/onefinger/structs"
	"github.com/WAY29/onefinger/utils"

	"github.com/asmcos/requests"
	"github.com/panjf2000/ants"
)

var (
	pool           *ants.PoolWithFunc
	wg             sync.WaitGroup
	requestOptions structs.RequestOptions
)

// 初始化协程池和请求选项
func Initialize(targets []string, threads int, options structs.RequestOptions) {
	pool, _ = ants.NewPoolWithFunc(threads, Run)
	requestOptions = options
}

// 开始任务
func Start(targetsInterface interface{}) {
	targets, ok := targetsInterface.([]string)
	if !ok {
		utils.OptionsError("Can't load targets as string slice", 3)
	}

	for _, target := range targets {
		wg.Add(1)
		pool.Invoke(target)
	}
}

// 运行任务
func Run(targetInterface interface{}) {
	defer wg.Done()

	target := targetInterface.(string)

	// 请求并返回结果
	err, result := fetch(target)
	if err != nil {
		fmt.Printf("[-] can't request target: " + target)
		return
	}

	// 指纹识别
	detect(result)
}

// 等待任务结束
func Wait() {
	wg.Wait()
}

// 任务结束
func End() {
	pool.Release()
}

// 发送请求并返回结果
func fetch(target string) (error, *structs.FetchResult) {
	// 设置请求选项
	req := requests.Requests()
	req.SetTimeout(time.Duration(requestOptions.Timeout))

	// 请求
	resp, err := req.Get(target)
	if err != nil {
		return err, nil
	}

	// 返回结构体

	var headerString string
	for k, v := range resp.R.Header {
		headerString += fmt.Sprintf("%v: %v\n", k, v[0])
	}

	fetchResult := structs.FetchResult{
		Url:          target,
		Content:      resp.Content(),
		Headers:      resp.R.Header,
		HeaderString: headerString,
		Certs:        getCerts(resp.R),
	}

	return nil, &fetchResult
}

// 获取证书内容，参考byro07/fwhatweb
func getCerts(resp *http.Response) []byte {
	var certs []byte
	if resp.TLS != nil {
		cert := resp.TLS.PeerCertificates[0]
		var str string
		if js, err := json.Marshal(cert); err == nil {
			certs = js
		}
		str = string(certs) + cert.Issuer.String() + cert.Subject.String()
		certs = []byte(str)
	}
	return certs
}

// 指纹识别
func detect(result *structs.FetchResult) {
	products := make([]string, 0)
	// 获取响应内容
	responseContent := strings.ToLower(string(result.Content))
	headerString := result.HeaderString
	headerServerString := fmt.Sprintf("Server : %v\n", result.Headers["Server"])
	certString := string(result.Certs)

	for _, fp := range FofaFingerPrints {
		rules := fp.Rules
		matchFlag := false

		for _, oneRule := range rules {
			ruleMatchContinueFlag := true
			// 单个rule之间是AND关系
			for _, rule := range oneRule {
				if !ruleMatchContinueFlag {
					break
				}
				lowerRuleContent := strings.ToLower(rule.Content)

				switch strings.Split(rule.Match, "_")[0] {
				case "banner":
					reBanner := regexp.MustCompile(`(?im)<\s*banner.*>(.*?)<\s*/\s*banner>`)
					matchResults := reBanner.FindAllString(responseContent, -1)
					if len(matchResults) == 0 {
						ruleMatchContinueFlag = false
						break
					}

					for _, matchResult := range matchResults {
						if !strings.Contains(strings.ToLower(matchResult), lowerRuleContent) {
							ruleMatchContinueFlag = false
						}
					}
				case "title":
					reTitle := regexp.MustCompile(`(?im)<\s*title.*>(.*?)<\s*/\s*title>`)
					matchResults := reTitle.FindAllString(responseContent, -1)
					if len(matchResults) == 0 {
						ruleMatchContinueFlag = false
						break
					}

					for _, matchResult := range reTitle.FindAllString(responseContent, -1) {
						if !strings.Contains(strings.ToLower(matchResult), lowerRuleContent) {
							ruleMatchContinueFlag = false
						}
					}
				case "body":
					if !strings.Contains(responseContent, lowerRuleContent) {
						ruleMatchContinueFlag = false
					}
				case "header":
					if !strings.Contains(headerString, rule.Content) {
						ruleMatchContinueFlag = false
					}
				case "server":
					if !strings.Contains(headerServerString, rule.Content) {
						ruleMatchContinueFlag = false
					}
				case "cert":
					if (result.Certs == nil) || (result.Certs != nil && !strings.Contains(certString, rule.Content)) {
						ruleMatchContinueFlag = false
					}
				// case "protocol":
				default:
					ruleMatchContinueFlag = false
				}

			}
			// 单个rule之间是AND关系，匹配成功
			if ruleMatchContinueFlag {
				matchFlag = true
				break
			}
		}

		// 多个rule之间是OR关系，匹配成功
		if matchFlag {
			products = append(products, fp.Product)
		}
	}
	printResult(result.Url, products)
}

func printResult(target string, products []string) {
	fmt.Printf("[+] %s %s\n", target, strings.Join(products, ", "))
}
