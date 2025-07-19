package httpapi

import (
	"context"
	webappv1 "djl.com/DjlD1/api/v1"
	"djl.com/DjlD1/djluser"
	"djl.com/DjlD1/jwt"
	"djl.com/DjlD1/k8s"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type Apphttp struct {
	Client client.Client
	Router *gin.Engine
}
type nodejson struct {
	Name   string `json:"name"`
	Fuzai  int64  `json:"fuzai"`
	Status string `json:"status"`
}

func (a *Apphttp) Getnode(c *gin.Context, m *metrics.Clientset) {
	ctx := context.Background()
	nodes := make([]nodejson, 0)
	nodes1 := v1.NodeList{}
	err := a.Client.List(ctx, &nodes1, &client.ListOptions{})

	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取node失败",
		})
		return
	}
	// 3. 查询节点资源使用情况
	nodeMetricsList, err := m.MetricsV1beta1().NodeMetricses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取负载失败",
		})
	}
	for i, node := range nodes1.Items {
		nodeInfo := nodejson{
			Name:   node.Name,
			Status: string(node.Status.Conditions[4].Status),
			Fuzai:  nodeMetricsList.Items[i].Usage.Memory().Value(),
		}
		nodes = append(nodes, nodeInfo)
	}
	c.JSON(200, gin.H{
		"data": nodes,
	})
}

type namespacejson struct {
	Name string `json:"name"`
}

func (a *Apphttp) Getnamespace(c *gin.Context) {
	ctx := context.Background()
	namespaceslist := make([]namespacejson, 0)
	namespaceList1 := v1.NamespaceList{}
	err := a.Client.List(ctx, &namespaceList1, &client.ListOptions{})
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取名称空间失败",
		})
		return
	}
	for _, namespace := range namespaceList1.Items {
		namespaceInfo := namespacejson{
			Name: namespace.Name,
		}
		namespaceslist = append(namespaceslist, namespaceInfo)
	}
	c.JSON(200, gin.H{
		"data": namespaceslist,
	})
}

type podsjson struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Stytus    bool   `json:"stytus"`
}

func (a *Apphttp) Getpods(c *gin.Context) {
	ctx := context.Background()
	podList := make([]podsjson, 0)
	podList1 := v1.PodList{}
	err := a.Client.List(ctx, &podList1, &client.ListOptions{})
	if err != nil {
		c.JSON(400, gin.H{
			"message": "获取pod列表失败",
		})
	}
	for _, pod := range podList1.Items {
		podInfo := podsjson{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Stytus:    pod.Status.ContainerStatuses[0].Ready,
		}
		podList = append(podList, podInfo)
	}
	c.JSON(200, gin.H{
		"data": podList,
	})

}
func (a *Apphttp) GetdelDjlD1(c *gin.Context) {
	ctx := context.Background()
	app := webappv1.DjlD1{}
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := a.Client.Get(ctx, types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &app)
	klog.Infof("GetDjlD1 %s %s %v", namespace, name, err)
	if err != nil {
		if errors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "DjlD1 Not Found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "DjlD1 Error",
		})
		return
	}
	err = a.Client.Delete(ctx, &app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "DjlD1 Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "DjlD1 Deleted",
	})
	return
}

type DjlD1json struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Spec      webappv1.DjlD1Spec `json:"spec"`
}

func (a *Apphttp) GetDjlD1s(c *gin.Context) {
	ctx := context.Background()
	app := webappv1.DjlD1List{}
	app1 := make([]DjlD1json, 0)
	err := a.Client.List(ctx, &app, &client.ListOptions{})
	for _, app2 := range app.Items {
		appInfo := DjlD1json{
			Name:      app2.Name,
			Namespace: app2.Namespace,
			Spec:      app2.Spec,
		}
		app1 = append(app1, appInfo)
	}
	klog.Infof("GetDjlD1 %v", err)
	if err != nil {
		if errors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "DjlD1 Not Found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "DjlD1 Error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": app1,
	})
	return
}

type ziyuan struct {
	Name      string             `json:"name"`
	Namespace string             `json:"namespace"`
	Labels    map[string]string  `json:"labels"`
	Spce      webappv1.DjlD1Spec `json:"spce"`
}

func (a *Apphttp) CreateDjlD1(c *gin.Context) error {
	ctx := context.Background()
	app := webappv1.DjlD1{}
	shili := &ziyuan{}
	err := c.BindJSON(shili)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "序列化失败",
		})
		return err
	}
	app.Name = shili.Name
	app.Namespace = shili.Namespace
	app.Labels = shili.Labels
	app.Spec = shili.Spce
	err = a.Client.Create(ctx, &app)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建失败",
		})
		return err
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "创建成功",
		"namespace": app.Namespace,
		"name":      app.Name,
	})
	return nil
}

type respod struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    bool   `json:"status"`
	NodeName  string `json:"node_name"`
	Cname     string `json:"cname"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域请求（生产注意限制）
	},
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization") //获取请求头中的Authorization字段
		if authHeader == "" {
			c.JSON(401, gin.H{
				"massage": "token未提供",
			})
			c.Abort()
			return
		}
		parts := strings.Fields(authHeader) //以空格切分为字符串切片
		if parts[0] != "Bearer" || len(parts) != 2 {
			c.JSON(401, gin.H{
				"massage": "请求头格式不合适",
			})
			c.Abort()
			return
		}
		token := parts[1]
		username, err := jwt.Verifytoken(token)
		if err != nil {
			c.JSON(200, gin.H{
				"massage": "token校验失败",
			})
			c.Abort()
			return
		}
		// 将 username 存入上下文
		c.Set("username", username)
		c.Next()
	}
}
func (a *Apphttp) Start(ctx context.Context) error {
	var kubeconfig string
	kubeconfig = "httpapi/config"                                 //指定config文件位置
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig) //使用config文件
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	a.Router = gin.Default()
	a.Router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		// 如果前端需要凭证，这里设置为 true，同时不要用 '*'
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	// 2. 创建 metrics client
	metricsClient, err := metrics.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create metrics client failed: %v", err)
	}
	{
		d1app := a.Router.Group("/ddapp")
		//d1app.Use(Logger())
		d1app.DELETE("/:namespace/:name", func(c *gin.Context) {
			a.GetdelDjlD1(c)
			return
		})
		d1app.GET("", func(c *gin.Context) {
			a.GetDjlD1s(c)
			return
		})
		d1app.POST("", func(c *gin.Context) {
			err := a.CreateDjlD1(c)
			if err != nil {
				klog.Infof("错误信息%v", err)
				return
			}
			return
		})
	}
	{
		login := a.Router.Group("/login")
		login.POST("/login", func(c *gin.Context) {
			u := djluser.User1{}
			err := c.BindJSON(&u)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}
			token, err := u.Login()
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"token": token,
			})
			return
		})
		login.POST("/register", func(c *gin.Context) {})
	}
	podclien := k8s.NewK8spod(clientset)
	{
		k8s := a.Router.Group("/k8s")
		//k8s.Use(Logger())
		k8s.GET("delete/:namespace/:name", func(c *gin.Context) {
			namespace := c.Param("namespace")
			name := c.Param("name")
			err := podclien.Deletepod(context.Background(), namespace, name)
			var jieguo string
			if err != nil {
				jieguo = "删除失败"
			} else {
				jieguo = "删除成功"
			}
			c.String(http.StatusOK, jieguo)
		})
		k8s.GET("select/:namespace", func(c *gin.Context) {
			namespace := c.Param("namespace")
			err, pods := podclien.Selectpod(namespace, context.Background())
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"shulaing": "0",
				})
				return
			}
			podes := make([]respod, len(pods))
			for i, pod := range pods {
				podes[i].Name = pod.Name
				// fmt.Println(pod.Name)
				podes[i].Namespace = pod.Namespace
				podes[i].Status = pod.Status.ContainerStatuses[0].Ready
				podes[i].NodeName = pod.Spec.NodeName
				podes[i].Cname = pod.Spec.Containers[0].Name
			}
			c.JSON(200, gin.H{
				"pods": podes,
			})
		})
		k8s.GET("log/:namespace/:pname", func(c *gin.Context) {
			namespace := c.Param("namespace")
			pname := c.Param("pname")
			log := podclien.GetLogPod(context.Background(), namespace, pname)
			c.JSON(200, gin.H{
				"日志内容": log,
			})
		})
		k8s.GET("wslog/:namespace/:pname", func(c *gin.Context) {
			namespace := c.Param("namespace")
			pname := c.Param("pname")
			Conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			defer Conn.Close()
			podclien.WSGetLogPod(Conn, c.Request.Context(), namespace, pname)
		})
		k8s.GET("webshell/:namespace/:pname/:cname", func(c *gin.Context) {
			namespace := c.Param("namespace")
			pname := c.Param("pname")
			cname := c.Param("cname")
			Conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				return
			}
			defer Conn.Close()
			err = podclien.WebShell(config, Conn, namespace, pname, cname)
			if err != nil {
				Conn.WriteMessage(websocket.TextMessage, []byte("执行失败: "+err.Error()))
			}
		})
	}
	{
		testtoken := a.Router.Group(("/testtoken"))
		testtoken.Use(Logger())
		testtoken.GET("/", func(c *gin.Context) {
			username, exists := c.Get("username")
			if !exists {
				c.JSON(500, gin.H{
					"message": "未获取到用户信息",
				})
				return
			}
			c.JSON(200, gin.H{
				"message": "token校验成功",
				"欢迎":    username,
			})
		})
	}
	{
		node := a.Router.Group("/node")
		node.GET("/", func(c *gin.Context) {
			a.Getnode(c, metricsClient)
			return
		})
	}
	{
		node := a.Router.Group("/namespace")
		node.GET("", func(c *gin.Context) {
			a.Getnamespace(c)
			return
		})
	}
	{
		node := a.Router.Group("/podes")
		node.GET("", func(c *gin.Context) {
			a.Getpods(c)
			return
		})
	}
	err = a.Router.Run()
	if err != nil {
		return err
	}
	return nil
}
