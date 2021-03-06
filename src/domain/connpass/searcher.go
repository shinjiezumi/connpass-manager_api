package connpass

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"connpass-manager/domain/connpass/api"
)

// ISearcher .
type ISearcher interface {
	Search(condition api.EventSearchQuery, page, count int) ([]*Event, error)
}

// NewSearcher .
func NewSearcher() ISearcher {
	return &Searcher{
		URL: api.EventSearchURL,
	}
}

// Searcher .
type Searcher struct {
	URL string
}

// Search connpassイベントを検索する
func (s *Searcher) Search(condition api.EventSearchQuery, page, count int) ([]*Event, error) {
	u, err := s.makeURL(condition, page, count)
	if err != nil {
		return nil, err
	}

	// 検索実行
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	hc := new(http.Client)
	res, err := hc.Do(req)
	if err != nil {
		return nil, err
	}

	// イベントに変換
	ret, err := s.makeResponse(res)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// 検索URLを生成する
func (s *Searcher) makeURL(c api.EventSearchQuery, page, count int) (*url.URL, error) {
	// URL生成
	u, err := url.Parse(s.URL)
	if err != nil {
		return nil, err
	}

	// クエリ生成
	q := u.Query()
	if c.EventID != nil {
		q.Set("event_id", strconv.Itoa(*c.EventID))
	}
	if c.Keyword != nil {
		q.Set("keyword", *c.Keyword)
	}
	if c.KeywordOr != nil {
		q.Set("keyword_or", *c.Keyword)
	}
	if c.Ym != nil {
		q.Set("ym", strconv.Itoa(*c.Ym))
	}
	if c.Ymd != nil {
		q.Set("ymd", strconv.Itoa(*c.Ymd))
	}
	if c.NickName != nil {
		q.Set("nickname", *c.NickName)
	}
	if c.OwnerNickname != nil {
		q.Set("owner_nickname", *c.OwnerNickname)
	}
	if c.SeriesID != nil {
		q.Set("series_id", strconv.Itoa(*c.SeriesID))
	}

	// 検索結果の何件目から出力するかを指定
	q.Set("start", strconv.Itoa(1+((page-1)*count)))
	q.Set("count", strconv.Itoa(count))
	u.RawQuery = q.Encode()

	return u, nil
}

// APIレスポンスを生成する
func (s *Searcher) makeResponse(res *http.Response) ([]*Event, error) {
	defer func() {
		if err := res.Body.Close(); err != nil {
			panic(err)
		}
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response api.EventResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	ret := make([]*Event, 0, len(response.Events))
	for _, e := range response.Events {
		ret = append(ret, &Event{
			Title:            e.Title,
			Description:      e.Description,
			URL:              e.EventURL,
			StartedAt:        e.StartedAt,
			EndedAt:          e.EndedAt,
			Limit:            e.Limit,
			Series:           e.Series,
			Address:          e.Address,
			Place:            e.Place,
			OwnerDisplayName: e.OwnerDisplayName,
			Accepted:         e.Accepted,
			Waiting:          e.Waiting,
		})
	}

	return ret, nil
}
