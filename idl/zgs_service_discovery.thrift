namespace go zgs_service_discovery

// 注册

struct RegistQuery {
    1: string uuid;
    2: string ip;
    3: string port;
}

struct RegistResponse {
    1: i32 code;
    2: string message;
}

// 查询可用服务列表

struct AgentInfo {
    1: string uuid;
    2: string ip;
    3: string port;
    4: string status;
    5: string group;
    6: string ext;
}

struct ListAgentsInfoQuery {
    1: list<string> group;
    2: list<string> status;
}

struct ListAgentsInfoResponse {
    1: i32 total;
    2: list<AgentInfo> agents;
}

service ZgsServiceDiscovery {
    RegistResponse Regist(1:RegistQuery request) (api.post="/regist");
    ListAgentsInfoResponse ListAgents(1:ListAgentsInfoQuery reqeust) (api.post="/list_agents");
}
