- 原始材料：

  - chainroll 包；
  - datasys 包；
  - dataflow 包；
  - dataflow 前端包；

- 操作：只能登录到当前服务器，进行单台机子部署（不能远程部署）；一条命令启动所有；

- 支持共用配置（例如 mongodb、mysql）（如果没有指定，则使用单独配置）；
- 支持单独配置；
- 支持单机“单套环境”配置（一个 chainroll、一个 datasys、一个 dataflow）；
- 支持单机“多套环境”配置；
- 支持的命令有：start、stop、restart、status、clear；当然最好还要支持一堆 update 操作；

- 和 docker 的启动命令类似，可以通过 bitxmesh-ops start stop restart 等命令来管理整个 BitXMesh 环境，也可以单独管理 chainroll datasys dataflow 等子系统 
- bitxmesh-ops start chainroll xxx // 启动名为 xxx 的 chainroll（如果不存在，则创建）（如果存在，则启动）（如果有库表，则创建）（） 
- bitxmesh-ops stop chainroll xxx // 停止名为 xxx 的 bitxmesh（递归停止所有组件，先停止最上层服务）（如果不存在，则报错） 
- bitxmesh-ops restart chainroll xxx // 重启名为 xxx 的 bitxmesh（递归停止、递归启动）（如果不存在，则报错） 
- bitxmesh-ops ps // 查看所有正在运行的子系统状态 
- bitxmesh-ops rm chainroll // 如果该进程正在运行，则停止并删除该进程（如果该进程正在被其他依赖，则提示）（同时也要删除所有的 mysql、mongo）


```json
{
  "work_dir": ".",
  // chainroll
  "mongo_ip": "172.22.67.81",
  "mongo_port": 27017,
  "mongo_username": "",
  "mongo_passwd": "",
  "chainroll_mongo_suffix": "172_22_67_81_bitxmesh_changshu",
  // dataflow 使用
  "mysql_ip": "172.22.67.81",
  "mysql_port": 3306,
  "mysql_username": "root",
  "mysql_passwd": "hyperchain@1n",

  "machines": [
    {
      "ip": "127.0.0.1",
      "orgs": [
        {
          // chainroll 配置
          "chainroll_grpc_port": 10100,
          "chainroll_http_port": 10200,
          // datasys
          "datasys_node_nums": 1,
          "datasys_grpc_port": 9400,
          "datasys_p2p_port": 7400,
          "dataflow_org_name": "常熟机构0",
          "dataflow_mysql_datasys_db": "172_22_67_81_bitxmesh_changshu_org0",
          // dataflow
          "dataflow_http_port": 9080,
          "dataflow_mysql_db": "172_22_67_81_bitxmesh_changshu_dataflow0"
        },
        {
          "chainroll_grpc_port": 10101,
          "chainroll_http_port": 10201,

          "datasys_node_nums": 1,
          "datasys_grpc_port": 9401,
          "datasys_p2p_port": 7401,
          "dataflow_org_name": "常熟机构1",
          "dataflow_mysql_datasys_db": "172_22_67_81_bitxmesh_changshu_org1",

          "dataflow_http_port": 9081,
          "dataflow_mysql_db": "172_22_67_81_bitxmesh_changshu_dataflow1"
        },
        {
          "chainroll_grpc_port": 10102,
          "chainroll_http_port": 10202,

          "datasys_node_nums": 1,
          "datasys_grpc_port": 9402,
          "datasys_p2p_port": 7402,
          "dataflow_org_name": "常熟机构2",
          "dataflow_mysql_datasys_db": "172_22_67_81_bitxmesh_changshu_org2",

          "dataflow_http_port": 9082,
          "dataflow_mysql_db": "172_22_67_81_bitxmesh_changshu_dataflow2"
        }
      ]
    }
  ],
  "hyperchain": {
    "nodes": "[\\\"172.22.67.127\\\", \\\"172.22.67.127\\\", \\\"172.22.67.127\\\", \\\"172.22.67.127\\\"]",
    "ports": "[\\\"8081\\\", \\\"8082\\\", \\\"8083\\\", \\\"8084\\\"]"
  }
}
```
