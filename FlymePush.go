package flyme_push

import (
	"github.com/cocotyty/httpclient"
	"encoding/json"
	"sort"
	"fmt"
	"crypto/md5"
	"encoding/hex"
	"github.com/heycayc/flyme_push/consts"

)

type FlymePush struct {
	AppId  string
	AppKey string
}

type Response struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Value    Value `json:"value,omitempty"`
	Redirect string `json:"redirect"`
	MsgID    string `json:"msgId"`
}

type Value struct {
	TaskId int `json:"taskid"`
	PushType int `json:"pushType"`
	AppId int `json:"appId"`
}




/**
 * 通过PushId推送透传消息
 */
func (f FlymePush) SendThroughByPushIds(pushIds, messageJson string) error {
	sendByPushIds := map[string]string{
		"appId":       f.AppId,
		"pushIds":     pushIds,
		"messageJson": messageJson,
	}
	sign := GenerateSign(sendByPushIds, f.AppKey)
	_, err := httpclient.
	Post(consts.PushThroughMessageByPushId).
		Head("charset", "UTF-8").
		Param("appId", f.AppId).
		Param("pushIds", pushIds). //多个逗号隔开
		Param("sign", sign).
		Param("messageJson", messageJson).
		Send().String()
	if err != nil {
		return err
	}

	return nil
}

//taskId推送接口（通知栏消息）
func (f FlymePush) SendGetTaskId(pushType string, messageJson string) (*Response, error) {
	pushNotificationMessageMap := map[string]string{
		"appId":       f.AppId,
		"pushType":     pushType,
		"messageJson": messageJson,
	}
	sign := GenerateSign(pushNotificationMessageMap, f.AppKey)
	respStr, err := httpclient.
	Post(consts.PushGetTaskId).
		Head("charset", "UTF-8").
		Param("appId", f.AppId).
		Param("pushType", pushType). //多个逗号隔开
		Param("sign", sign).
		Param("messageJson", messageJson).
		Send().String()
	fmt.Println(respStr)
	res := &Response{}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(respStr), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (f FlymePush) SendMessageByTaskId(taskId string,pushIds string ) (*Response, error){
	pushNotificationMessageMap := map[string]string{
		"taskId":       taskId,
		"appId":       f.AppId,
		"pushIds":     pushIds,
	}
	sign := GenerateSign(pushNotificationMessageMap, f.AppKey)
	respStr, err := httpclient.
	Post(consts.PushNotificationMessageTaskId).
		Head("charset", "UTF-8").
		Param("taskId", taskId).
		Param("appId", f.AppId).
		Param("pushIds", pushIds). //多个逗号隔开
		Param("sign", sign).
		Send().String()
	res := &Response{}
	fmt.Println(respStr)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(respStr), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//pushId推送接口（通知栏消息）
func (f FlymePush) SendNotificationMessageByPushId(pushIds string, messageJson string) (*Response, error) {
	pushNotificationMessageMap := map[string]string{
		"appId":       f.AppId,
		"pushIds":     pushIds,
		"messageJson": messageJson,
	}
	sign := GenerateSign(pushNotificationMessageMap, f.AppKey)
	respStr, err := httpclient.
	Post(consts.PushNotificationMessageByPushId).
		Head("charset", "UTF-8").
		Param("appId", f.AppId).
		Param("pushIds", pushIds). //多个逗号隔开
		Param("sign", sign).
		Param("messageJson", messageJson).
		Send().String()
	res := &Response{}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(respStr), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//别名推送接口（通知栏消息）
func (f FlymePush) SendNotificationMessageByAlias(alias string, messageJson string) error {
	maps := map[string]string{
		"appId":       f.AppId,
		"alias":       alias,
		"messageJson": messageJson,
	}

	sign := GenerateSign(maps, f.AppKey)
	_, err := httpclient.
	Post(consts.PushNotificationMessageByAlias).
		Head("charset", "UTF-8").
		Param("appId", f.AppId).
		Param("alias", alias). //推送别名，一批最多不能超过1000个多个英文逗号分割（必填）
		Param("sign", sign).
		Param("messageJson", messageJson).
		Send().String()
	if err != nil {
		return err
	}
	return nil
}

/*******************************************标签推送****************************************/
//全部推送
func (f FlymePush) SendAllMessage(pushType string, messageJson string) (*Response, error)  {
	maps := map[string]string{
		"appId":       f.AppId,
		"pushType":    pushType,
		"messageJson": messageJson,
	}

	sign := GenerateSign(maps, f.AppKey)
	respStr , err := httpclient.
	Post(consts.PushAllMessage).
		Head("charset", "UTF-8").
		Param("appId", f.AppId).
		Param("pushType", pushType).
		Param("sign", sign).
		Param("messageJson", messageJson).
		Send().String()
	if err != nil {
		return nil, err
	}
	res := &Response{}
	err = json.Unmarshal([]byte(respStr), res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// 标签推送
func (f FlymePush) SendMessageByTipic(pushType string, tagNames string, scope string, messageJson string) error {
	maps := map[string]string{
		"appId":       f.AppId,
		"pushType":    pushType,
		"tagNames":    tagNames,
		"scope":       scope,
		"messageJson": messageJson,
	}

	sign := GenerateSign(maps, f.AppKey)
	_, err := httpclient.
	Post(consts.PushMessageByTipic).
		Head("charset", "UTF-8").
		Param("appId", f.AppId).
		Param("pushType", pushType).
		Param("tagNames", tagNames).
		Param("scope", scope). //标签集合（必填）0 并集1 交集
		Param("sign", sign).
		Param("messageJson", messageJson).
		Send().String()
	if err != nil {
		return err
	}
	return nil
}

/********************************推送统计*******************************/
//PushStatistics get 请求
func (f FlymePush) SendStatistics(taskId string) (string, error) {
	maps := map[string]string{
		"appId":  f.AppId,
		"taskId": taskId,
	}
	sign := GenerateSign(maps, f.AppKey)
	str , err := httpclient.
	Post(consts.PushStatistics).
		Head("charset", "UTF-8").
		Param("appId", f.AppId).
		Param("taskId", taskId).
		Param("sign", sign).
		Send().String()
	if err != nil {
		return "", err
	}
	return str, nil
}
func GenerateSign(params map[string]string, appKey string) string {
	var signStr string
	if params != nil {
		keys := make([]string, len(params))
		i := 0
		for key, _ := range params {
			keys[i] = key
			i++
		}
		sort.Strings(keys)
		for _, k := range keys {
			signStr += k + "=" + params[k]
		}
		signStr += appKey
	}
	return PushParamMD5(signStr)
}

func PushParamMD5(encodeStr string) string {
	hasher := md5.New()
	hasher.Write([]byte(encodeStr))
	return hex.EncodeToString(hasher.Sum(nil))
}

/*json方法*/
func JSON(i interface{}) string {
	outi, err := json.Marshal(i)
	if err != nil {
		panic("json出错了～～")
	}
	return string(outi)
}
