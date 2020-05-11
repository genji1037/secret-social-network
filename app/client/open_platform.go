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
	"strings"
	"time"
)

// OpenResult represents open platform respond.
type OpenResult struct {
	Code int                    `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg"`
}

// ApplyPaymentArgs represent args for apply payment.
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
	m["s"] = generateSignCode(m, cfg.AppKeys.GetByAppID(args.AppID))
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

// GetUID is convenience func get uid by open id and app id.
func GetUID(appID, openID1, openID2 string) (string, string, error) {
	UIDs, err := BatchGetUID(appID, []string{openID1, openID2})
	if err != nil {
		return "", "", err
	}
	var uid1, uid2 string
	for _, uidInfo := range UIDs {
		if uidInfo.OpenID == openID1 {
			uid1 = uidInfo.UID
		}
		if uidInfo.OpenID == openID2 {
			uid2 = uidInfo.UID
		}
	}
	return uid1, uid2, nil
}

// BatchGetUID get uid from open platform
func BatchGetUID(appID string, openID []string) ([]UIDInfo, error) {
	cfg := config.GetServe().Open
	url := fmt.Sprintf("%s/manager/user/uid/batch?app_id=%s&open_ids=%s", cfg.BaseURL, appID, strings.Join(openID, ","))
	rsp, err := get(url)
	if err != nil {
		return nil, err
	}

	jsb, err := json.Marshal(rsp.Data["uids"])
	if err != nil {
		return nil, err
	}
	var rs []UIDInfo
	err = json.Unmarshal(jsb, &rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func get(url string) (rsp *OpenResult, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
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
