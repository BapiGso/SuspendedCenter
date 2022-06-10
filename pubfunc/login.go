package blogfunc

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/gorilla/sessions"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"text/template"
)

// 初始化存储器（基于 Cookie）
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY)")))

//计算密码的sha1值
func hash(passwd string) string {
	h := sha1.New() // md5加密类似md5.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(passwd))
	passwdhash := h.Sum(nil)
	h.Reset()
	passwdhash16 := hex.EncodeToString(passwdhash) //转16进制
	return passwdhash16
}

func Loginget(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	admin, _ := template.ParseFiles("gtpl/admin/admin.html")
	welcome, _ := template.ParseFiles("gtpl/admin/login.html")
	session, _ := store.Get(r, "GOSESSID")
	if session.Values["count"] == true {
		admin.Execute(w, nil)
	} else {
		welcome.Execute(w, nil)
	}
}

func Loginpost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session, _ := store.Get(r, "GOSESSID")
	passwd := r.PostFormValue("password")
	passwdhash := hash(passwd)
	//密码对了发true
	if passwdhash == "40bd001563085fc35165329ea1ff5c5ecbdbbeef" {
		session.Values["count"] = true
	} else {
		session.Values["count"] = false
	}
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
	}
	//密码对错都跳转get处理
	http.Redirect(w, r, "/login", 302)

}
