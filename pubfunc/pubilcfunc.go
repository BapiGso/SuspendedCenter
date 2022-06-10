package blogfunc

import (
	"github.com/russross/blackfriday/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//时间戳转时间
func unix2time(unix int64) string {
	mon := [...]string{"", "一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"}
	tm := time.Unix(unix, 0)
	tmp := tm.Format("04 02,2006")
	index, _ := strconv.Atoi(tmp[:2])
	tmp = strings.Replace(tmp, tmp[:2], mon[index], 1)
	return tmp
}

//markdown转html
func md2html(md string) string {
	md = md[16:] //去掉<!--markdown-->前缀
	tmp := []byte(md)
	output := string(blackfriday.Run(
		tmp,
		blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.HardLineBreak),
	))
	return output
}

//判断请求是否是ajax
func isAjax(r *http.Request, vars map[string]interface{}) map[string]interface{} {
	if r.Header["X-Requested-With"] == nil {
		vars["noheadfoot"] = true
		return vars
	} else {
		vars["noheadfoot"] = false
		return vars
	}
}
