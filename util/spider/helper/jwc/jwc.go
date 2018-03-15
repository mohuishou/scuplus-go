// Package jwc 教务处helper, 利用chromedp返回教务处cookie
package jwc

import (
	"context"
	"log"

	"github.com/chromedp/chromedp/client"

	"github.com/chromedp/cdproto/network"

	"github.com/chromedp/cdproto/cdp"

	"github.com/chromedp/chromedp"
)

// GetCookies 获取cookie字符串
func GetCookies() string {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)))
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	// run task list
	var cookieStr string
	err = c.Run(ctxt, chromedp.Tasks{
		chromedp.Navigate(`http://jwc.scu.edu.cn/jwc/moreNotice.action`),
		chromedp.WaitVisible(`table`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context, h cdp.Executor) error {
			cookies, err := network.GetAllCookies().Do(ctx, h)
			for _, v := range cookies {
				cookieStr = cookieStr + v.Name + "=" + v.Value + ";"
			}
			if err != nil {
				return err
			}
			return nil
		}),
	})
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	return cookieStr
}
