package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/asmcos/requests"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)


type FofaFinger struct {
	RuleId         string `json:"rule_id"`
	Level          string `json:"level"`
	Softhard       string `json:"softhard"`
	Product        string `json:"product"`
	Company        string `json:"company"`
	Category       string `json:"category"`
	ParentCategory string `json:"parent_category"`
	Rules          [][]struct {
		Match   string `json:"match"`
		Content string `json:"content"`
	} `json:"rules"`
}

type FetchResult struct {
	Url           string
	Content       []byte
	Headers       http.Header
	HeadersString string
	Certs         []byte
}

//解析json指纹
func Parse(filename string) ([]FofaFinger, error) {
	Json, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var dataArray []FofaFinger
	err = json.Unmarshal(Json, &dataArray)
	if err != nil {
		return nil, err
	}
	return dataArray, nil
}

func Fetchbody(resp *FetchResult) {

	products := make([]string, 0)
	//获取网页返回数据并赋值
	web_Content := strings.ToLower(string(resp.Content)) //强制转换小写
	certString := string(resp.Certs)
	web_Certs := resp.Certs
	web_HeadersString := resp.HeadersString
	headerSeverString := fmt.Sprintf("Server: %v\n", resp.Headers["Server"]) //%v： 值得默认格式表示  根据Server参数生成格式化的字符串并返回该字符串

	fofajson, _ := Parse("fofa.json")

	for _, fp := range fofajson {
		//fofa指纹中的最后一项
		rules := fp.Rules
		matchFlag := false
		//对每个json的最后一项进行迭代
		for _, onerule := range rules {
			//控制继续器
			ruleMatchContinueFlag := true

			for _, rule := range onerule {
				if !ruleMatchContinueFlag {
					break
				}
				lowerRuleContent := strings.ToLower(rule.Content)
				switch strings.Split(rule.Match, "_")[0] {

				case "banner":
					reBanner := regexp.MustCompile(`(?im)<\s*banner.*>(.*?)<\s/\s*banner>`)
					matchResults := reBanner.FindAllString(web_Content, -1)
					if len(matchResults) == 0 {
						ruleMatchContinueFlag = false
						break
					}
					for _, matchResult := range matchResults {
						if !strings.Contains(strings.ToLower(matchResult), lowerRuleContent) {
							ruleMatchContinueFlag = false
							break
						}
					}
				case "title":
					reTitle := regexp.MustCompile(`(?im)<\s*title.*(.*?)<\s*/\s*title>`)
					matchResults := reTitle.FindAllString(web_Content, -1)
					if len(matchResults) == 0 {
						ruleMatchContinueFlag = false
						break
					}
					for _, matchResult := range reTitle.FindAllString(web_Content, -1) {
						if !strings.Contains(strings.ToLower(matchResult), lowerRuleContent) {
							ruleMatchContinueFlag = false
						}
					}
				case "body":
					if !strings.Contains(web_Content, lowerRuleContent) {
						ruleMatchContinueFlag = false
					}
				case "header":
					if !strings.Contains(web_HeadersString, rule.Content) {
						ruleMatchContinueFlag = false
					}
				case "Server":
					if !strings.Contains(headerSeverString, rule.Content) {
						ruleMatchContinueFlag = false
					}
				case "cert":
					if (web_Certs == nil) || (web_Certs != nil && !strings.Contains(certString, rule.Content)) {
						ruleMatchContinueFlag = false
					}
				default:
					ruleMatchContinueFlag = false
				}
				//单个rule之间是and关系
				if ruleMatchContinueFlag {
					matchFlag = true
					break
				}

			}
		}
		//多个rule之间是or关系
		if matchFlag {
			products = append(products, fp.Product)
		}
	}
	PrintResult(resp.Url, products)
}
func PrintResult(target string, products []string) {
	fmt.Print("[+] %s %s \n", target, strings.Join(products, ","))
}

func Reqdata(url string) (*FetchResult, error) {
	req := requests.Requests()
	req.SetTimeout(time.Duration(10))
	resp, err := req.Get(url)

	if err != nil {
		fmt.Println(err)
	}
	var headerString string
	req_data := FetchResult{
		Url:           url,
		Content:       resp.Content(),
		Headers:       resp.R.Header,
		HeadersString: headerString,
		Certs:         GetCerts(resp.R),
	}
	return &req_data, nil
}

func GetCerts(resp *http.Response) []byte {
	var certs []byte
	if resp.TLS != nil {
		cert := resp.TLS.PeerCertificates[0]
		var str string
		if js, err := json.Marshal(cert); err == nil {
			certs = js
		}
		str = string(certs) + cert.Subject.String()
		certs = []byte(str)
	}
	return certs
}

func Banner() {
	banner :=
		`
         _                   _                   _                  _                  _               _                   _
        /\ \     _          / /\                /\ \     _         / /\              /\ \             / /\                /\ \     _
       /  \ \   /\_\       / /  \              /  \ \   /\_\      / /  \            /  \ \           / /  \              /  \ \   /\_\
      / /\ \ \_/ / /      / / /\ \            / /\ \ \_/ / /     / / /\ \__        / /\ \ \         / / /\ \            / /\ \ \_/ / /
     / / /\ \___/ /      / / /\ \ \          / / /\ \___/ /     / / /\ \___\      / / /\ \ \       / / /\ \ \          / / /\ \___/ /
    / / /  \/____/      / / /  \ \ \        / / /  \/____/      \ \ \ \/___/     / / /  \ \_\     / / /  \ \ \        / / /  \/____/
   / / /    / / /      / / /___/ /\ \      / / /    / / /        \ \ \          / / /    \/_/    / / /___/ /\ \      / / /    / / /
  / / /    / / /      / / /_____/ /\ \    / / /    / / /     _    \ \ \        / / /            / / /_____/ /\ \    / / /    / / /
 / / /    / / /      / /_________/\ \ \  / / /    / / /     /_/\__/ / /       / / /________    / /_________/\ \ \  / / /    / / /
/ / /    / / /      / / /_       __\ \_\/ / /    / / /      \ \/___/ /       / / /_________\  / / /_       __\ \_\/ / /    / / /
\/_/     \/_/       \_\___\     /____/_/\/_/     \/_/        \_____\/        \/____________/  \_\___\     /____/_/\/_/     \/_/
                                                      Nanscan  version  : 1.0.0
                                                      Made     by       : PuPp1T.

`

	fmt.Println(banner)

}

func main() {
	Banner()

	var (
		Url     string
		File    string
		Threads int
	)

	fmt.Println("(-U=<targetUrl> | -F=<target File> | -T=<threads>)")

	flag.StringVar(&Url, "u", "", "输入url")
	flag.StringVar(&File, "f", "", "文件内为url")
	flag.IntVar(&Threads, "t", 10, "线程默认为10")

	flag.Parse()
	//执行的函数
	if Url != "" && File == "" {
		req, _ := Reqdata(Url)
		Fetchbody(req)

	}
}
