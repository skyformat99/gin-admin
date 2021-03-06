package web

import (
	"fmt"

	"github.com/LyricTian/gin-admin/src/inject"
	"github.com/LyricTian/gin-admin/src/web/context"
	"github.com/LyricTian/gin-admin/src/web/middleware"
	"github.com/LyricTian/gin-admin/src/web/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Init 初始化所有服务
func Init(obj *inject.Object) *gin.Engine {
	gin.SetMode(viper.GetString("run_mode"))
	app := gin.New()

	// 注册中间件
	apiPrefixes := []string{"/api/"}

	if dir := viper.GetString("web_dir"); dir != "" {
		app.Use(middleware.WWWMiddleware(dir, apiPrefixes...))
	}

	app.Use(middleware.TraceMiddleware(apiPrefixes...))
	app.Use(middleware.LoggerMiddleware(apiPrefixes))
	app.Use(middleware.RecoveryMiddleware())
	app.Use(middleware.SessionMiddleware(obj, apiPrefixes...))

	app.NoMethod(func(c *gin.Context) {
		context.New(c).ResError(fmt.Errorf("方法不允许"), 405)
	})

	app.NoRoute(func(c *gin.Context) {
		context.New(c).ResError(fmt.Errorf("资源不存在"), 404)
	})

	// 注册/api/v1路由
	router.APIV1Handler(app, obj)

	// 加载casbin策略数据
	err := loadCasbinPolicyData(obj)
	if err != nil {
		panic("加载casbin策略数据发生错误：" + err.Error())
	}

	return app
}

// 加载casbin策略数据，包括角色权限数据、用户角色数据
func loadCasbinPolicyData(obj *inject.Object) error {
	c := obj.CtlCommon

	err := c.RoleAPI.RoleBll.LoadAllPolicy()
	if err != nil {
		return err
	}

	err = c.UserAPI.UserBll.LoadAllPolicy()
	if err != nil {
		return err
	}
	return nil
}
