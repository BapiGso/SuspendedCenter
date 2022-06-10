package blogfunc

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"
	"text/template"
	"unicode/utf8"
)

//查询数据库首页信息
func indexdbquery() (string, string) {
	db, err := sql.Open("sqlite", "assets/123.db")
	checkErr(err)
	//查询标题、创建日期、内容、观看数、喜欢数
	var title, slug, text, views, likes, str_value string
	var created int64
	var containerfinal, pyimaryfinal bytes.Buffer
	//查文章的表然后从自定义字段表跟着查背景图片
	rows, err := db.Query("SELECT  title,slug,created,text,views,likes,str_value FROM typecho_contents  left join typecho_fields on typecho_contents.slug=typecho_fields.cid WHERE typecho_contents.type='post' and typecho_fields.name='coverList' ORDER BY typecho_contents.rowid DESC LIMIT 5\n")
	checkErr(err)

	i := 1
	container, _ := template.ParseFiles("gtpl/index-container.html")
	primary, _ := template.ParseFiles("gtpl/index-primary.html")
	for rows.Next() {
		err = rows.Scan(&title, &slug, &created, &text, &views, &likes, &str_value)
		vars := make(map[string]interface{}, 9)
		vars["title"] = title
		vars["slug"] = slug
		vars["created"] = unix2time(created)
		vars["text"] = text
		vars["views"] = views
		vars["likes"] = likes
		vars["count"] = utf8.RuneCountInString(text)
		vars["subcount"] = string([]rune(text)[:80])
		vars["coverList"] = str_value
		if i == 1 {
			container.Execute(&containerfinal, vars)
		} else {
			primary.Execute(&pyimaryfinal, vars)
		}
		i += 1
	}
	db.Close()
	return containerfinal.String(), pyimaryfinal.String()
}

func Indexpage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cid := ps.ByName("name")
	//判断是否为数字
	_, err := strconv.ParseFloat(cid, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	fmt.Fprint(w, cid)
}

func Index(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	containerfinal, pyimaryfinal := indexdbquery()
	vars := make(map[string]interface{}, 3)
	vars["container"] = containerfinal
	vars["pyimary"] = pyimaryfinal
	vars = isAjax(r, vars)
	index, _ := template.ParseFiles("gtpl/index.html", "gtpl/header.html", "gtpl/footer.html")
	index.Execute(w, vars)
}
