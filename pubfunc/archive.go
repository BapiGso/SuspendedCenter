package blogfunc

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	_ "modernc.org/sqlite"
	"net/http"
	"strconv"
	"text/template"
	"unicode/utf8"
)

//查询数据库文章信息返回一个map
func dbquery(cid string) map[string]interface{} {
	db, _ := sql.Open("sqlite", "assets/123.db")

	//查询标题、创建日期、内容、观看数、喜欢数
	var title, text, views, likes string
	var created int64
	_ = db.QueryRow("SELECT title,created,text,views,likes FROM typecho_contents WHERE slug ='"+cid+"'").Scan(&title, &created, &text, &views, &likes)
	db.Close()
	customtime := unix2time(created)
	text = md2html(text)
	vars := make(map[string]interface{})
	vars["title"] = title + " - "
	vars["created"] = customtime
	vars["text"] = text
	vars["views"] = views
	vars["likes"] = likes
	vars["count"] = utf8.RuneCountInString(text)
	vars["subcount"] = string([]rune(text)[:80])
	return vars

}

func Content(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cid := ps.ByName("name")
	post, _ := template.ParseFiles("gtpl/archive-post.html", "gtpl/header.html", "gtpl/footer.html")
	//判断是否为数字
	_, err := strconv.ParseFloat(cid, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	fill := dbquery(cid)
	fill = isAjax(r, fill)

	post.Execute(w, fill)

}
