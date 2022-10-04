namespace go zgs_service_discovery

// 注册

struct RegistRequest {
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
    4: string status;// online offline
    5: string group;
    6: string ext;
}

struct ListAgentsByGroupAndStatusRequest {
    1: list<string> group;
    2: list<string> status;
}

struct ListAgentsByGroupAndStatusResponse {
    1: i32 total;
    2: list<AgentInfo> agents;
}

service ZgsServiceDiscovery {
    RegistResponse Regist(1:RegistRequest request) (api.post="/regist");
    ListAgentsByGroupAndStatusResponse ListAgents(1:ListAgentsByGroupAndStatusRequest reqeust) (api.post="/list_agents");
}
