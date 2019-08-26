package wx

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

const (
	Scope_Base     = "snsapi_base"     // 静默登录
	Scope_UserInfo = "snsapi_userinfo" // 获取用户信息

	Uri_Authorize = "https://open.weixin.qq.com/connect/oauth2/authorize"

	CookieName_Redirect = "_wx_redirect"
)

type HTTPServer interface {
	NewHandler() http.Handler
}

type wxProxyHTTPServer struct {
	appId      string
	webroot    string
	fs         http.Handler
	allowHosts []string
}

// NewWXProxyHTTPServer
func NewWXProxyHTTPServer() HTTPServer {
	log.Println(viper.AllSettings())
	appId := viper.GetString("app_id")
	if len(appId) == 0 {
		log.Panic("缺少微信公众号 AppID配置")
	}
	webroot := viper.GetString("web_root_dir")
	if len(webroot) == 0 {
		log.Panic("缺少WEB根目录路径配置")
	}
	allowHosts := viper.GetStringSlice("allow_hosts")

	return &wxProxyHTTPServer{
		appId:      appId,
		webroot:    webroot,
		fs:         http.FileServer(http.Dir(webroot)),
		allowHosts: allowHosts,
	}
}

func (s *wxProxyHTTPServer) NewHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/connect/oauth2/authorize", s.newAuthorizeHandleFunc())
	mux.HandleFunc("/redirect", s.newRedirectHandleFunc())
	mux.HandleFunc("/", s.newWebRootHandleFunc())

	return mux
}

func (s *wxProxyHTTPServer) newRedirectHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie(CookieName_Redirect)
		if err != nil {
			http.Error(w, "缺少重定向地址", http.StatusBadRequest)
			return
		}
		redirectTo, err := url.Parse(cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		redirectTo.RawQuery = request.URL.Query().Encode()
		http.Redirect(w, request, redirectTo.String(), http.StatusTemporaryRedirect)
	}
}

func (s *wxProxyHTTPServer) newAuthorizeHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var (
			errcode = http.StatusBadRequest
			err     error

			redirectTo    string
			redirectToUri *url.URL
			hostAllowed   bool
			state         string
			scope         string

			cookie http.Cookie

			proxyRedirect = request.URL.Scheme + "://" + request.URL.Host + "/redirect"
			authorizeUri  *url.URL // 重定向微信链接
			query         = make(url.Values)
		)

		redirectTo = request.URL.Query().Get("redirect_uri")
		if len(redirectTo) == 0 {
			err = errors.New("缺少 redirect_uri 参数")
			goto errret
		}

		if redirectToUri, err = url.Parse(redirectTo); err != nil {
			goto errret
		}

		if len(s.allowHosts) > 0 {
			for _, h := range s.allowHosts {
				if h == redirectToUri.Host {
					hostAllowed = true
					break
				}
			}
		} else {
			hostAllowed = true
		}

		if !hostAllowed {
			err = errors.New("重定向域名不在白名单中")
			goto errret
		}

		scope = request.URL.Query().Get("scope")
		if len(scope) == 0 {
			scope = Scope_Base
		}
		state = request.URL.Query().Get("state")

		if authorizeUri, err = url.Parse(Uri_Authorize); err != nil {
			goto errret
		}

		query.Set("appid", s.appId)
		query.Set("redirect_uri", proxyRedirect)
		query.Set("response_type", "code")
		query.Set("scope", scope)
		if len(state) > 0 {
			query.Set("state", state)
		}
		authorizeUri.RawQuery = query.Encode()
		authorizeUri.Fragment = "wechat_redirect"

		cookie = http.Cookie{
			Name:     CookieName_Redirect,
			Value:    redirectTo,
			HttpOnly: true,
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, request, authorizeUri.String(), http.StatusTemporaryRedirect)

		return
	errret:
		w.WriteHeader(errcode)
		w.Write([]byte(err.Error()))
	}
}

func (s *wxProxyHTTPServer) newWebRootHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		s.fs.ServeHTTP(w, request)
	}
}
