package v1

import (
	"github.com/galaxy-future/BridgX/cmd/api/handler/gf-cluster/cluster"
	"github.com/galaxy-future/BridgX/cmd/api/handler/gf-cluster/instance"
	"github.com/galaxy-future/BridgX/cmd/api/handler/gf-cluster/kubernetes"
	"github.com/galaxy-future/BridgX/cmd/api/middleware/authorization"
	"github.com/galaxy-future/BridgX/internal/gf-cluster/calibrator"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(route *gin.RouterGroup) {
	route.Use(authorization.CheckTokenAuth())
	calibrator.Init()

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
		instanceGroupRoute.POST("", instance.HandleCreateInstanceGroup)
		instanceGroupRoute.POST("/batch/create", instance.HandleBatchCreateInstanceGroup)

		instanceGroupRoute.GET("/delete/:instanceGroup", instance.HandleDeleteInstanceGroup)
		instanceGroupRoute.DELETE("/delete/:instanceGroup", instance.HandleDeleteInstanceGroup)
		instanceGroupRoute.POST("/delete/:instanceGroup", instance.HandleDeleteInstanceGroup)
		instanceGroupRoute.POST("/batch/delete", instance.HandleBatchDeleteInstanceGroup)

		instanceGroupRoute.POST("/update", instance.HandleUpdateInstanceGroup)
		instanceGroupRoute.PATCH("/update", instance.HandleUpdateInstanceGroup)

		instanceGroupRoute.GET("", instance.HandleListInstanceGroup)
		instanceGroupRoute.GET("/:instanceGroup", instance.HandleGetInstanceGroup)

		instanceGroupRoute.POST("/expand_shrink", instance.HandleExpandOrShrinkInstanceGroup)

		instanceRoute := route.Group("/instance")
		instanceRoute.POST("/restart", instance.HandleRestartInstance)
		instanceRoute.GET("/:instanceGroup", instance.HandleListInstance)
		instanceRoute.GET("/self", instance.HandleListMyInstance)
		instanceRoute.POST("/delete", instance.HandleDeleteInstance)
		instanceRoute.GET("/form", instance.HandleListInstanceForm)
	}

	clusterRoute := route.Group("/cluster")
	{
		clusterRoute.GET("/bridgx/available_clusters", cluster.HandleListUnusedBridgxCluster)

		clusterRoute.DELETE("/:clusterId", cluster.HandleDeleteKubernetes)
		clusterRoute.POST("", cluster.HandleCreateCluster)

		clusterRoute.GET("/summary", cluster.HandleListClusterSummary)
		clusterRoute.GET("/summary/:clusterId", cluster.HandleGetClusterSummary)
		clusterRoute.GET("/nodes/:clusterId", cluster.HandleListNodesSummary)
		clusterRoute.GET("/pods/:clusterId", cluster.HandleListClusterPodsSummary)
	}

}
