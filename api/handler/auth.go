package handler

import (
	"encoding/base64"
	"github.com/132yse/acgzone-server/api/db"
	"github.com/132yse/acgzone-server/api/def"
	"github.com/julienschmidt/httprouter"
	"github.com/132yse/acgzone-server/api/util"
	"net/http"
	"strconv"
)

//登陆校验，只负责校验登陆与否
func Auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	Cross(w, r)
	uid, err := r.Cookie("uid")
	if err != nil || uid == nil {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return
	}
	id, _ := strconv.Atoi(uid.Value)
	resp, err := db.GetUser("", id)
	if err != nil {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return
	}

	res := &def.User{Id: resp.Id, Name: resp.Name, Role: resp.Role, QQ: resp.QQ, Desc: resp.Desc}
	sendUserResponse(w, res, 201, "")
}

//鉴权校验，负责判断是否具有编辑和审核权限
func RightAuth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) string {
	uname, _ := r.Cookie("uname")
	name, _ := base64.StdEncoding.DecodeString(uname.Value)

	resp, err := db.GetUser(string(name), 0)
	if err != nil {
		sendErrorResponse(w, def.ErrorNotAuthUser)
		return ""
	}
	token, _ := r.Cookie("token")                      //从 cookie 里取 token
	newToken := util.CreateToken(resp.Name, resp.Role) //服务端生成新的 token

	if token.Value == newToken { //已经登陆
		if resp.Role == "admin" || resp.Role == "editor" {
			return resp.Role
		}
	} else {
		return ""
	}
	return ""
}

func Cross(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Max-Age", "3600")
	w.Header().Add("Access-Control-Allow-Headers", "x-requested-with")
}
