package v1

import (
	"github.com/galaxy-future/BridgX/cmd/api/handler/gf-cluster/eci"
	"github.com/galaxy-future/BridgX/cmd/api/handler/gf-cluster/gf_cloud"
	"github.com/galaxy-future/BridgX/cmd/api/handler/gf-cluster/kubernetes"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/authorization"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(route *gin.RouterGroup) {

	route.Use(authorization.CheckTokenAuth())

	kubeRoute := route.Group("/kubernetes")
	{
		kubeRoute.POST("", kubernetes.HandleRegisterKubernetes)
		kubeRoute.GET("", kubernetes.HandleListKubernetes)

		kubeRoute.POST("/update", kubernetes.HandleUpdateKubernetes)
		kubeRoute.PATCH("/update", kubernetes.HandleUpdateKubernetes)

		kubeRoute.GET("/:cluster", kubernetes.HandleGetKubernetes)
	}


	instanceGroupRoute := route.Group("/instance_group")
	{
		instanceGroupRoute.POST("", eci.HandleCreateInstanceGroup)
		instanceGroupRoute.POST("/batch/create", eci.HandleBatchCreateInstanceGroup)

		instanceGroupRoute.GET("/delete/:instanceGroup", eci.HandleDeleteInstanceGroup)
		instanceGroupRoute.DELETE("/delete/:instanceGroup", eci.HandleDeleteInstanceGroup)
		instanceGroupRoute.POST("/delete/:instanceGroup", eci.HandleDeleteInstanceGroup)
		instanceGroupRoute.POST("/batch/delete", eci.HandleBatchDeleteInstanceGroup)

		instanceGroupRoute.POST("/update", eci.HandleUpdateInstanceGroup)
		instanceGroupRoute.PATCH("/update", eci.HandleUpdateInstanceGroup)

		instanceGroupRoute.GET("", eci.HandleListInstanceGroup)
		instanceGroupRoute.GET("/:instanceGroup", eci.HandleGetInstanceGroup)

		instanceGroupRoute.POST("/expand", eci.HandleExpandInstanceGroup)
		instanceGroupRoute.POST("/shrink", eci.HandleShrinkInstanceGroup)
		instanceGroupRoute.POST("/expand_shrink", eci.HandleExpandOrShrinkInstanceGroup)

		instance := route.Group("/instance")
		instance.POST("/restart", eci.HandleRestartInstance)
		instance.GET("/:instanceGroup", eci.HandleListInstance)
		instance.GET("/self", eci.HandleListMyInstance)
		instance.POST("/delete", eci.HandleDeleteInstance)
		instance.GET("/form", eci.HandleListInstanceForm)
	}


	clusterRoute := route.Group("/cluster")
	{
		clusterRoute.GET("/bridgx/available_clusters", gf_cloud.HandleListUnusedBridgxCluster)

		clusterRoute.DELETE("/:clusterId", gf_cloud.HandleDeleteKubernetes)
		clusterRoute.POST("", gf_cloud.HandleCreateCluster)

		clusterRoute.GET("/summary", gf_cloud.HandleListClusterSummary)
		clusterRoute.GET("/summary/:clusterId", gf_cloud.HandleGetClusterSummary)
		clusterRoute.GET("/nodes/:clusterId", gf_cloud.HandleListNodesSummary)
		clusterRoute.GET("/pods/:clusterId", gf_cloud.HandleListClusterPodsSummary)
	}


}
