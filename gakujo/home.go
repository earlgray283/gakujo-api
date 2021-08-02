package gakujo

import (
	"bytes"
	"io"
	"net/url"

	"github.com/szpp-dev-team/gakujo-api/model"
	"github.com/szpp-dev-team/gakujo-api/scrape"
)

func (c *Client) Home() (model.HomeInfo, error) {
	body, err := c.fetchHomeHtml()
	if err != nil {
		return model.HomeInfo{}, err
	}
	defer func() {
		body.Close()
		_, _ = io.Copy(io.Discard, body)
	}()
	b, _ := io.ReadAll(body)
	taskRows, err := scrape.TaskRows(io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return model.HomeInfo{}, err
	}
	noticeRows, err := scrape.NoticeRows(io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return model.HomeInfo{}, err
	}
	return model.HomeInfo{
		TaskRows:   taskRows,
		NoticeRows: noticeRows,
	}, nil
}

func (c *Client) fetchHomeHtml() (io.ReadCloser, error) {
	datas := make(url.Values)
	datas.Set("headTitle", "ホーム")
	datas.Set("menuCode", "Z07") // TODO: menucode を定数化(まとめてやる)
	datas.Set("nextPath", "/home/home/initialize")

	return c.getPage(GeneralPurposeUrl, datas)
}

func (c *Client) fetchNoiceDetailhtml() (io.ReadCloser, error) {
	reqURL := "https://gakujo.shizuoka.ac.jp/portal/portaltopcommon/newsForTop/deadLineForTop"

	params := make(url.Values)
	params.Set("org.apache.struts.taglib.html.TOKEN", c.token)
	params.Set("newsTargetIndexNo", "20")

	return c.getPage(reqURL, params)
}
