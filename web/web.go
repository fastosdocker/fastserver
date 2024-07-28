package web

import (
	"fast/config"
	"fast/pkg/engine"
	"fast/web/app"
	"fast/web/base"
	"fast/web/cmd"
	"fast/web/compose"
	"fast/web/container"
	"fast/web/daemon"
	"fast/web/home"
	"fast/web/hub"
	"fast/web/image"
	"fast/web/mw"
	"fast/web/repo"
	"fast/web/store"
	"fast/web/user"
	"fast/web/volume"
	"log"
	"net/http"
	"time"

	"github.com/fvbock/endless"
)

// Start 开启web服务
func Start(s time.Time) {
	r := engine.Get()
	r.GET("/time", base.Handler(func(c *base.Ctx) {
		c.Success(s)
	}))

	r.Use(mw.TlsHandler(), base.Handler(mw.CrossDomain))
	r.GET("/getversion", base.Handler(app.GetVersion)) //不需要前端 后端 => 线上后端   docker我们  => 服务端
	r.GET("/update", base.Handler(app.Update))
	r.Use(base.Handler(mw.Log))
	r.StaticFile("/", "./static/index.html")
	r.StaticFS("/pc", http.Dir("./static/pc"))
	r.POST("/login", base.Handler(user.Login))
	r.POST("/register", base.Handler(user.Register))
	r.GET("/ws", base.Handler(container.Ws))
	// r.Use(base.Handler(mw.JWTAuth))
	r.POST("/docker", base.Handler(cmd.Shell))
	r.POST("/password", base.Handler(user.Password))
	r.POST("/imagesave", base.Handler(image.Save))
	r.POST("/imageload", base.Handler(image.Load))
	r.POST("/imagesprune", base.Handler(image.Prune))
	r.POST("/hubadd", base.Handler(hub.Add))
	r.POST("/hubdel", base.Handler(hub.Del))
	r.POST("/hubsave", base.Handler(hub.Save))
	r.GET("/hubfind", base.Handler(hub.Find))
	r.POST("/daemonsave", base.Handler(daemon.Save))
	r.GET("/daemonload", base.Handler(daemon.Load))
	// tls证书
	r.POST("/tls/uploader", base.Handler(app.UploadTls))
	r.POST("/tls/signers", base.Handler(app.StartTls))
	r.GET("/tls/status", base.Handler(app.GetTlsStatus))
	// 重启服务
	r.POST("/restart", base.Handler(app.Reload))
	// 容器
	r.GET("/container/list", base.Handler(container.List))
	r.POST("/container/create", base.Handler(container.Create))
	r.POST("/container/stats", base.Handler(container.Stats))
	r.POST("/container/historystats", base.Handler(container.HistoryStats))
	r.POST("/container/delete", base.Handler(container.Delete))
	r.POST("/container/repo", base.Handler(repo.GetDockerRepo))
	// compose
	r.POST("/compose/create", base.Handler(compose.Create))
	r.GET("/compose/del", base.Handler(compose.Del))
	r.POST("/compose/update", base.Handler(compose.Update))
	r.GET("/compose/log", base.Handler(compose.Log))
	r.GET("/compose/stop", base.Handler(compose.Stop))
	r.GET("/compose/info", base.Handler(compose.Get))
	r.GET("/compose/list", base.Handler(compose.List))
	r.GET("/compose/page", base.Handler(store.Page))
	r.GET("/compose/composeUpdate", base.Handler(store.Update))
	r.GET("/compose/class", base.Handler(store.Class))
	r.POST("/compose/pageId", base.Handler(store.PageId))
	r.GET("/compose/run", base.Handler(compose.Run))
	// volume
	r.POST("/volume/create", base.Handler(volume.Create))
	r.GET("/volume/list", base.Handler(volume.List))
	r.GET("/volume/remove", base.Handler(volume.Remove))
	r.GET("/home", base.Handler(home.Overview))

	if config.Conf.Https.Flag {
		go func() {
			err := endless.ListenAndServeTLS(":"+config.Conf.Https.Port, config.Conf.TLs.Pem, config.Conf.TLs.Key, r)
			if err != nil {
				log.Println("err:", err)
			}
		}()
	}
	err := endless.ListenAndServe(":"+config.Conf.Port, r)
	if err != nil {
		log.Println("err:", err)
	}
}
