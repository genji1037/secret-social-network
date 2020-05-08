package client

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"secret-social-network/app/config"
	"sort"
	"strconv"
	"time"
)

// OpenResult represents open platform respond.
type OpenResult struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg"`
}

type ApplyPaymentArgs struct {
	AppID   string
	OpenID  string
	OrderID string
	Token   string
	Amount  decimal.Decimal
	Remark  string
}

// ApplyPayment apply trade_no from open platform
func ApplyPayment(args ApplyPaymentArgs) (string, error) {
	cfg := config.GetServe().Open
	m := make(map[string]interface{})
	m["open_id"] = args.OpenID
	m["app_id"] = args.AppID
	m["amount"] = args.Amount
	m["token"] = args.Token
	m["order_id"] = args.OrderID
	m["pay_type"] = 20
	m["remark"] = args.Remark
	m["t"] = time.Now().Unix()
	m["s"] = generateSignCode(m, cfg.SecretKey)
	url := cfg.BaseURL + "/payment/create"

	rsp, err := post(url, m)
	if err != nil {
		return "", err
	}

	tradeNo, ok := rsp.Data["trade_no"].(string)
	if !ok {
		return "", fmt.Errorf("convert trade_no to string failed")
	}

	return tradeNo, nil
}

func GetUID(appID, openID1, openID2 string) (string, string, error) {
	UIDs, err := BatchGetUID(appID, []string{openID1, openID2})
	if err != nil {
		return "", "", err
	}
	if len(UIDs) < 2 {
		return "", "", fmt.Errorf("unexpected uids returned by openplatfrom")
	}
	return UIDs[0], UIDs[1], nil
}

func BatchGetUID(appID string, openID []string) ([]string, error) {
	// todo
	return openID, nil
}

func post(url string, m map[string]interface{}) (rsp *OpenResult, err error) {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	for k, v := range m {
		vStr := getStringFromGivenType(v)
		writer.WriteField(k, vStr)
	}
	writer.Close()
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return
	}
	req.Header.Add("content-type", writer.FormDataContentType())

	defer req.Body.Close()

	client := &http.Client{}
	client.Timeout = time.Minute
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	rsp = new(OpenResult)
	body, _ := ioutil.ReadAll(resp.Body)
	if err = json.Unmarshal(body, rsp); err != nil {
		return rsp, fmt.Errorf("response: body:%s", string(body))
	}
	return
}

func generateSignCode(m map[string]interface{}, secretKey string) (signCode string) {
	s := make([]string, 0)
	for k := range m {
		s = append(s, k)
	}
	sort.Strings(s)

	str := ""
	for _, v := range s {
		str += getStringFromGivenType(m[v])
	}
	str += secretKey
	signCode = fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return
}

// GetStringFromType could get string from given type
func getStringFromGivenType(v interface{}) string {
	var str string
	switch v.(type) {
	case int:
		str = strconv.Itoa(v.(int))
	case int64:
		str = strconv.FormatInt(v.(int64), 10)
	case string:
		str, _ = v.(string)
	case decimal.Decimal:
		str = v.(decimal.Decimal).String()
	default:
		str = ""
	}
	return str
}
