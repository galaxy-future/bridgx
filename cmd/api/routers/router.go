package routers

import (
	"net/http"

	"github.com/galaxy-future/BridgX/cmd/api/handler"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/authorization"
	"github.com/galaxy-future/BridgX/config"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	var router *gin.Engine
	if config.GlobalConfig.DebugMode {
		gin.SetMode(gin.DebugMode)
		router = gin.Default()
		//visit http://0.0.0.0:9090/debug/pprof/
		pprof.Register(router)
	} else {
		gin.SetMode(gin.ReleaseMode)
		router = gin.Default()
	}

	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello world!")
	})

	user := router.Group("/user")
	{
		user.POST("/login", handler.Login)
		user.Use(authorization.RefreshTokenConditionCheck()).POST("/refresh_token", handler.RefreshToken)
	}

	v1Api := router.Group("/api/v1/")
	v1Api.Use(authorization.CheckTokenAuth())

	{
		cloudAccountPath := v1Api.Group("cloud_account/")
		{
			cloudAccountPath.POST("create", handler.CreateCloudAccount)
			cloudAccountPath.GET("list", handler.ListCloudAccounts)
			cloudAccountPath.GET("info", handler.GetAccountInfo)
			cloudAccountPath.POST("edit", handler.EditCloudAccount)
			cloudAccountPath.DELETE("delete/:ids", handler.DeleteCloudAccount)
		}

		clusterPath := v1Api.Group("cluster/")
		{
			clusterPath.GET("id/:id", handler.GetClusterById)
			clusterPath.GET("name/:name", handler.GetClusterByName)
			clusterPath.POST("create", handler.CreateCluster)
			clusterPath.POST("edit", handler.EditCluster)
			clusterPath.DELETE("delete/:ids", handler.DeleteClusters)
			clusterPath.GET("num", handler.GetClusterCount)
			clusterPath.GET("instance_stat", handler.GetInstanceStat)
			clusterPath.GET("describe_all", handler.ListClusters)

			clusterPath.POST("list_by_tags", handler.ListClustersByTags)
			clusterPath.GET("get_tags", handler.GetClusterTags)
			clusterPath.POST("add_tags", handler.AddClusterTags)
			clusterPath.POST("edit_tags", handler.EditClusterTags)
			clusterPath.DELETE("delete_tags", handler.DeleteClusterTags)

			clusterPath.POST("expand", handler.ExpandCluster)
			clusterPath.POST("shrink", handler.ShrinkCluster)
			clusterPath.POST("shrink_all", handler.ShrinkAllInstances)
		}
		vpcPath := v1Api.Group("vpc/")
		{
			vpcPath.POST("create", handler.CreateVpc)
			vpcPath.GET("describe", handler.DescribeVpc)
		}
		subnetPath := v1Api.Group("subnet/")
		{
			subnetPath.POST("create", handler.CreateSwitch)
			subnetPath.GET("describe", handler.DescribeSwitch)

		}
		groupPath := v1Api.Group("security_group/")
		{
			groupPath.POST("create", handler.CreateSecurityGroup)
			groupPath.GET("describe", handler.DescribeSecurityGroup)
			groupPath.POST("rule/add", handler.AddSecurityGroupRule)
			groupPath.POST("create_with_rule", handler.CreateSecurityGroupWithRules)
		}
		networkPath := v1Api.Group("network_config/")
		{
			networkPath.POST("create", handler.CreateNetworkConfig)
		}
		regionPath := v1Api.Group("region/")
		{
			regionPath.GET("list", handler.ListRegions)
		}
		zonePath := v1Api.Group("zone/")
		{
			zonePath.GET("list", handler.ListZones)
		}
		instanceTypePath := v1Api.Group("instance_type/")
		{
			instanceTypePath.GET("list", handler.ListInstanceType)
		}
		instancePath := v1Api.Group("instance/")
		{
			instancePath.GET("num", handler.GetInstanceCount)
			instancePath.GET("id/describe", handler.GetInstance)
			instancePath.GET("describe_all", handler.GetInstanceList)
			instancePath.GET("usage_total", handler.GetInstanceUsageTotal)
			instancePath.GET("usage_statistics", handler.GetInstanceUsageStatistics)
		}
		taskPath := v1Api.Group("task/")
		{
			taskPath.GET("num", handler.GetTaskCount)
			taskPath.GET("list", handler.GetTaskList)
			taskPath.GET("describe", handler.GetTaskDescribe)
			taskPath.GET("describe_all", handler.GetTaskDescribeAll)
			taskPath.GET("instances", handler.GetTaskInstances)
		}
		userPath := v1Api.Group("user/")
		{
			userPath.GET("info", handler.GetUserInfo)
			userPath.POST("create_ram_user", handler.CreateUser)
			userPath.POST("modify_password", handler.ModifyAdminPassword)
			userPath.POST("modify_username", handler.ModifyUsername)
			userPath.POST("enable_ram_user", handler.EnableUser)
			userPath.GET("list", handler.ListUsers)
		}
		orgPath := v1Api.Group("org/")
		{
			orgPath.POST("create", handler.CreateOrg)
			orgPath.POST("edit", handler.EditOrg)
			orgPath.GET("list", handler.ListOrgs)
			orgPath.GET("id/:id", handler.GetOrgById)
		}
		imagePath := v1Api.Group("image/")
		{
			imagePath.GET("list", handler.GetImageList)
		}
	}
	return router
}
