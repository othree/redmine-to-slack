package main

import (
	"fmt"
	"time"
	"bytes"
	"net/http"
	"math/rand"
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"github.com/vaughan0/go-ini"
	"github.com/kardianos/osext"
)

func main() {
	zeromessage := "太好了，現在所有的議題都有人負責"
	normalmessage := "還有 %d 個議題沒人負責喔"

	messages := []string{
		"我身上太多議題要爆炸啦啊啊啊啊啊啊啊啊",
		"看啊，我體內的議題清單長的這麼大了",
		"曾經我身上是沒有議題的，直到我的膝蓋中了一箭",
		"那一天，人類終於回想起了，曾經一度支配他們的恐懼，還有 redmine 議題沒處理完的屈辱",
		"議題太多、時間太少",
	}

	premessages := []string{
		"以下精選議題，弟子速速接手",
		"帶個無主的議題回家吧",
		"沒人負責的議題，就像是放羊的孩子，講的話都沒人聽",
		"給看到的人：King Bob 指定你負責以下議題",
		"**s 看一下有沒有你喜歡的議題吧",
		"接手一個議題，勝造七級浮屠",
	}

	pwd, _ := osext.ExecutableFolder()

	config, err := ini.LoadFile(pwd + "/config.ini")
	if err != nil {
		panic("Config file not loaded.")
	}

	feed_url, ok := config.Get("feed", "url")

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if ok {
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL(feed_url)
		// fmt.Printf("%T\n", feed.Items)
		// fmt.Printf("%+v\n", feed)

		msg := make(map[string]interface{})
		var attaches []map[string]interface{}

		count := len(feed.Items)
		mainmessage := zeromessage

		if count > 0 {
			mainmessage = fmt.Sprintf(normalmessage, count)

			attach := make(map[string]interface{})
			attach["color"] = "#7CD197"
			attach["title"] = "無主事務清單"
			list, _ := config.Get("feed", "list")
			attach["title_link"] = list

			attaches = append(attaches, attach)
		}
		if count >= 10 {
			i := r.Intn(len(messages))
			mainmessage = messages[i]
		}

		// fmt.Printf("%s\n", mainmessage)
		msg["text"] = mainmessage

		if count > 0 {
			i := r.Intn(len(premessages))
			premessage := premessages[i]
			// fmt.Printf("%s\n", premessage)

			issue := feed.Items[0]
			// fmt.Printf("%+v\n", issue.Title)
			// fmt.Printf("%+v\n", issue.Link)

			attach := make(map[string]interface{})
			attach["color"] = "#7CD197"
			attach["title"] = issue.Title
			attach["title_link"] = issue.Link
			attach["pretext"] = premessage

			attaches = append(attaches, attach)

			for i := 1; i < count; i++ {
				pick := r.Intn(2)
				if pick == 1 {
					issue := feed.Items[i]
					// fmt.Printf("%T\n", issue)
					// fmt.Printf("%+v\n", issue.Title)
					// fmt.Printf("%+v\n", issue.Link)
					// fmt.Printf("%+v\n", issue.)
					attach := make(map[string]interface{})
					attach["color"] = "#7CD197"
					attach["title"] = issue.Title
					attach["title_link"] = issue.Link
					attaches = append(attaches, attach)
				}
			}
		}

		msg["attachments"] = attaches

		b, _ := json.Marshal(msg)
		// body := string(b)
		// fmt.Println(body)

		post_url, ok := config.Get("slack", "url")

		if ok {
			req, err := http.NewRequest("POST", post_url, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
					panic(err)
			}
			defer resp.Body.Close()

			fmt.Println("response Status:", resp.Status)
		}
	}

}