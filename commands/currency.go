package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	aw "github.com/deanishe/awgo"
	"github.com/sxyazi/alfred-workflows/utils"
	collect "github.com/sxyazi/go-collection"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var currencyIcons = map[string]string{
	"AUD": "australia.png",
	"CAD": "canada.png",
	"CHF": "switzerland.png",
	"CNY": "china.png",
	"EUR": "europe.png",
	"GBP": "uk.png",
	"HKD": "hk.png",
	"INR": "india.png",
	"JPY": "japan.png",
	"USD": "usa.png",
}

type currency struct {
	wf    *aw.Workflow
	rates map[string]float64
}

func (c *currency) prepare() {
	if len(c.rates) > 0 {
		return
	}
	reload := func() (any, error) {
		return c.update()
	}
	c.rates = make(map[string]float64)
	_ = c.wf.Cache.LoadOrStoreJSON("currency-rates", 30*time.Minute, reload, &c.rates)
}

func (c *currency) update() (map[string]float64, error) {
	data := url.Values{}
	data.Set("f.req", `[[["Ba1tad","[[[\"/g/11bvv_1vt1\"],[\"/g/11bvvzdz__\"],[\"/g/11bvvzmlkc\"],[\"/g/11bvvzq4m1\"],[\"/g/11bvvzlhck\"],[\"/g/11bvvy0g2l\"],[\"/g/11bvvzv0_k\"],[\"/g/11bvvzjvx_\"],[\"/g/11bvvznqzd\"],[\"/g/11bvvzh029\"],[\"/g/11bvv_1vxq\"],[\"/g/11bvvzrnfv\"],[\"/g/11bvvzh84q\"],[\"/g/12m8wtr9h\"],[\"/m/01nj9h\"],[\"/g/1q58fgpkq\"],[\"/m/0cqyw\"],[\"/m/02xl7xj\"],[\"/g/11c1wl07cj\"],[\"/g/1218f7nx\"],[\"/m/04xk2h\"],[\"/g/11bc6t2dgw\"],[\"/m/016yss\"],[\"/m/0ckcw1p\"],[\"/g/11c20sl1zv\"],[\"/m/0ckzfmt\"],[\"/m/0ckvbnk\"],[\"/m/0ckthcg\"],[\"/m/0ckyjwv\"],[\"/g/11bc6md9fs\"],[\"/m/0ckhqlx\"],[\"/m/03qn8cx\"],[\"/g/11bvvzkdm5\"],[\"/g/11bvvztstd\"],[\"/g/11bvvxtc_c\"],[\"/g/11bvvzp9n7\"],[\"/g/11bvvzpb58\"],[\"/g/11bvvykq82\"]]]",null,"generic"]]]`)
	body, err := utils.HttpPost("https://www.google.com/finance/_/GoogleFinanceUi/data/batchexecute", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// match the json from response
	main := regexp.MustCompile(`(?m)^\[.*]$`).FindSubmatch(body)
	if len(main) < 1 {
		return nil, errors.New("invalid response")
	}

	// parse the json
	var parsed []any
	if err := json.Unmarshal(main[0], &parsed); err != nil {
		return nil, err
	}

	// parse the sub json
	sub, _ := utils.DeepValue[string](parsed, "0.2")
	if err := json.Unmarshal([]byte(sub), &parsed); err != nil {
		return nil, err
	}

	// parse the data of currency
	parsed, err = collect.AnyGet[[]any](parsed, "0")
	result := map[string]float64{"USD": 1}
	for i := 0; i < len(parsed); i++ {
		from, _ := utils.DeepValue[string](parsed, fmt.Sprintf("%d.8.0", i))
		to, _ := utils.DeepValue[string](parsed, fmt.Sprintf("%d.8.1", i))
		if from == "USD" && to != "" {
			result[to], _ = utils.DeepValue[float64](parsed, fmt.Sprintf("%d.1.0", i))
		}
	}

	return result, nil
}

func (c *currency) toward(s, from, to string) error {
	amount, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	for k, v := range c.rates {
		if k == to {
			continue
		}
		if from != "" && !strings.Contains(k, from) {
			continue
		}

		round := math.Round(amount/v*c.rates[to]*100) / 100
		c.wf.
			NewItem(fmt.Sprintf("%s %s = %.2f %s", s, k, round, to)).
			Arg(fmt.Sprintf("%.2f", round)).
			Valid(true).
			Icon(&aw.Icon{Value: c.wf.Dir() + "/static/flags/" + currencyIcons[k]})
	}
	c.wf.SendFeedback()
	return nil
}

func (c *currency) exec() error {
	c.prepare()
	if len(c.rates) < 1 {
		return errors.New("can not get currency rates")
	}

	arg := strings.TrimSpace(c.wf.Args()[1])
	if arg == "" {
		return errors.New("invalid argument")
	}

	if arg[len(arg)-1] == '$' {
		return c.toward(arg[:len(arg)-1], "USD", "CNY")
	}
	if m := regexp.MustCompile(`([\d.]+)([a-zA-Z]{1,3})`).FindStringSubmatch(arg); len(m) > 2 {
		return c.toward(m[1], strings.ToUpper(m[2]), "CNY")
	}
	return c.toward(arg, "", "CNY")
}
