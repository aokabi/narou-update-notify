package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Allcount       int    `json:"allcount,omitempty"`
	Title          string `json:"title,omitempty"`
	Ncode          string `json:"ncode,omitempty"`
	Userid         int    `json:"userid,omitempty"`
	Writer         string `json:"writer,omitempty"`
	Story          string `json:"story,omitempty"`
	Biggenre       int    `json:"biggenre,omitempty"`
	Genre          int    `json:"genre,omitempty"`
	Gensaku        string `json:"gensaku,omitempty"`
	Keyword        string `json:"keyword,omitempty"`
	GeneralFirstup string `json:"general_firstup,omitempty"`
	GeneralLastup  string `json:"general_lastup,omitempty"`
	NovelType      int    `json:"novel_type,omitempty"`
	End            int    `json:"end,omitempty"`
	GeneralAllNo   int    `json:"general_all_no,omitempty"`
	Length         int    `json:"length,omitempty"`
	Time           int    `json:"time,omitempty"`
	Isstop         int    `json:"isstop,omitempty"`
	Isr15          int    `json:"isr15,omitempty"`
	Isbl           int    `json:"isbl,omitempty"`
	Isgl           int    `json:"isgl,omitempty"`
	Iszankoku      int    `json:"iszankoku,omitempty"`
	Istensei       int    `json:"istensei,omitempty"`
	Istenni        int    `json:"istenni,omitempty"`
	PcOrK          int    `json:"pc_or_k,omitempty"`
	GlobalPoint    int    `json:"global_point,omitempty"`
	DailyPoint     int    `json:"daily_point,omitempty"`
	WeeklyPoint    int    `json:"weekly_point,omitempty"`
	MonthlyPoint   int    `json:"monthly_point,omitempty"`
	QuarterPoint   int    `json:"quarter_point,omitempty"`
	YearlyPoint    int    `json:"yearly_point,omitempty"`
	FavNovelCnt    int    `json:"fav_novel_cnt,omitempty"`
	ImpressionCnt  int    `json:"impression_cnt,omitempty"`
	ReviewCnt      int    `json:"review_cnt,omitempty"`
	AllPoint       int    `json:"all_point,omitempty"`
	AllHyokaCnt    int    `json:"all_hyoka_cnt,omitempty"`
	SasieCnt       int    `json:"sasie_cnt,omitempty"`
	Kaiwaritu      int    `json:"kaiwaritu,omitempty"`
	NovelupdatedAt string `json:"novelupdated_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// 一つの作品の情報を返す
func GetNovelInfo(ctx context.Context) ([]Response, error) {
	// create a new http client
	client := &http.Client{}

	// create a new request
	url := "https://api.syosetu.com/novelapi/api/"
	// 小説ごとにふられるID
	ncode := "n2267be" // Ｒｅ：ゼロから始める異世界生活
	respFormat := "json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("ncode", ncode)
	q.Add("out", respFormat)
	req.URL.RawQuery = q.Encode()

	// send the request
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read the response body
	response := make([]Response, 0)
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	if len(response) != 2 {
		return nil, fmt.Errorf("failed to get novel info: len=%d", len(response))
	}

	return response, nil

}
