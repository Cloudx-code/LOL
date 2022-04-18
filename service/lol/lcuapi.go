package lol

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// GetRoomId 获取当前对战聊天房间的ID
func GetRoomId() (string, error) {
	msg, err := cli.httpGet("/lol-chat/v1/conversations")
	if err != nil {
		return "", err
	}
	conversations := make([]Conversation, 10)
	json.Unmarshal(msg, &conversations)
	for _, conversation := range conversations {
		fmt.Printf("conversation: %+v\n", conversation)
		if conversation.Type == "championSelect" {
			return conversation.Id, nil
		}
	}
	return "", nil
}

// GetSummonerListByRoomId 通过聊天房间ID,查出聊天记录,从中找到5个召唤师的id值
func GetSummonerListByRoomId(roomId string) []int64 {
	msg, _ := cli.httpGet(fmt.Sprintf("/lol-chat/v1/conversations/%v/messages", roomId)) //得到这个房间内的所有消息
	lolMsgs := make([]LolMessage, 10)
	json.Unmarshal(msg, &lolMsgs)
	summonerIds := make([]int64, 0)
	for _, lolmsg := range lolMsgs {
		fmt.Printf("房间消息为%v\n", lolmsg)
		if lolmsg.Type == "system" { //系统发出的消息,形如xxx进入房间
			summonerIds = append(summonerIds, lolmsg.FromSummonerID)
		}
	}
	return summonerIds
}

// GetSummonerInfoById 根据召唤师id查找召唤师的完整信息
func GetSummonerInfoById(id int64) (SummonerInfo, error) {
	msg, err := cli.httpGet(fmt.Sprintf("/lol-summoner/v2/summoners?ids=[%v]", id))
	summoners := make([]SummonerInfo, 1)
	if err != nil {
		return SummonerInfo{}, err
	}

	json.Unmarshal(msg, &summoners)
	return summoners[0], nil
}

// ListGamesBySummonerID 根据召唤师id,查询最近[begin,begin+limit-1]的游戏战绩
func ListGamesBySummonerID(summonerId int64, begin, limit int) (*GameListResp, error) {
	bts, err := cli.httpGet(fmt.Sprintf("/lol-match-history/v3/matchlist/account/%d?begIndex=%d&endIndex=%d",
		summonerId, begin, begin+limit))
	if err != nil {
		return nil, err
	}
	data := &GameListResp{}
	json.Unmarshal(bts, data)
	return data, nil
}

// SendConversationMsg 根据房间id发送消息
func SendConversationMsg(msg interface{}, roomId string) error {
	TempByte, _ := json.Marshal(msg)
	data := struct { //发送消息时,服务端指定格式的数据
		Body string `json:"body"`
		Type string `json:"type"`
	}{
		Body: string(TempByte),
		Type: "chat",
	}
	fmt.Println("\n\n\n$$$$$$$$$$$$$$\n", data.Body)
	//json str 转map
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(data.Body), &dat); err == nil {
		fmt.Println("==============json str 转map=======================")
		fmt.Println(dat)
		fmt.Printf("type：%t\n", dat["currKDA"])
		//switch dat["currKDA"].(type) {
		//case string:
		//	fmt.Println(1111111)
		//case [][3]int:
		//	fmt.Println(2222222)
		//}
	} else {
		fmt.Println(err)
	}
	kda := strings.Split(data.Body, "currKDA\":")[1]
	kda = kda[:len(kda)-1]
	fmt.Println("kda:", kda)
	name := dat["summonerName"].(string)
	fmt.Println("name", name)
	score := dat["score"].(float64)
	fmt.Println("score", score)
	horseType := ""
	if score >= 10 {
		horseType = "上等马"
	} else if score >= 5 {
		horseType = "中等马"
	} else if score >= 3 {
		horseType = "下等马"
	} else {
		horseType = "牛马"
	}
	data.Body = horseType + ":" + name + ",近7场比赛KDA:" + strconv.FormatFloat(score, 'f', 3, 64) + "最近三把战绩为:" + kda
	if name == "皮皮九逗比小青年" || name == "Just随便一个名字" {
		data.Body = "恭喜你匹配到了传说中的半人马（偶然能当个人）:+" + name + "，希望本场比赛他能够发光发热，不要搞事！半人马的最近三把战绩为：" + kda
		data.Body = "恭喜你匹配到了传说中的:+" + name + "（具体为什么神秘咱也不知道）他最近三把的战绩为：" + kda
	}
	if name == "第一把位" {
		if horseType == "下等马" || horseType == "牛马" {
			data.Body = "恭喜你匹配到了神秘的:+" + name + "这个人很神秘" + kda
		}
	}
	if name == "好机油鳄鱼" || name == "寒带火熊" {
		if horseType == "下等马" || horseType == "牛马" {
			data.Body = "匹配到的是作者本人:" + name + ",战况保密！"
		}
	}
	if name == "盈盈一水間" {
		if horseType == "下等马" || horseType == "牛马" {
			data.Body = "匹配到的是作者开黑的小伙伴:" + name + ",战况保密！"
		}
	}
	mess, err := cli.httpPost(fmt.Sprintf("/lol-chat/v1/conversations/%s/messages", roomId), data)
	fmt.Println("响应为", string(mess))
	return err
}
