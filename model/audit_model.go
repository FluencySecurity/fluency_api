package model

type AuditEvent struct {
	Account  string `json:"account"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	ClientIP string `json:"clientIP"`
	Country  string `json:"country"`
	City     string `json:"city"`

	CreatedOn  int64  `json:"createdOn"`
	Dur        int64  `json:"dur"`
	Category   string `json:"category"`
	Error      string `json:"error"`
	DayIndex   string `json:"dayIndex"`
	Partition  string `json:"partition"`
	Action     string `json:"action"`
	ActionType string `json:"action_type"`
	ActionInfo string `json:"action_info"`
}

type DbStatusResponse struct {
	Status     string               `json:"status"`
	EtcdStatus *ClusterStatus  `json:"etcdStatus"`
	Master     *MasterStatus `json:"master"`
	Indexes    []*IndexStatus       `json:"indexes"`
	Nodes      []*NodeStatus        `json:"nodes,omitempty"`
}

type NodeStatus struct {
	NodeName  string `json:"nodeName"`
	IP        string `json:"ip"`
	Ec2ID     string `json:"ec2ID,omitempty"`
	NodeMgrUp bool   `json:"nodeMgrUp"`
	WorkerUp  bool   `json:"workerUp"`
}

type IndexStatus struct {
	IndexName string       `json:"indexName,omitempty"`
	Queue     *QueueStatus `json:"queue"`
}

type QueueStatus struct {
	Length     uint   `json:"length"`
	LastUpdate int64  `json:"lastUpdate"`
	Status     string `json:"status"`
}

type EtcdNodeStatus struct {
	ID uint64 `json:"id"`

	Name     string `json:"name"`
	IsLeader bool   `json:"isLeader"`
	DBSize   int64  `json:"dbSize"`
	Endpoint string `json:"endpoint"`
	Ip       string `json:"ip"`
	Active   bool   `json:"active"`
}
type ClusterStatus struct {
	Nodes  []*EtcdNodeStatus `json:"nodes"`
	Status string            `json:"status"` // "green","orange", "red"
}
type MasterStatus struct {
	Leader  string   `json:"leader"`
	Standby []string `json:"standby"`
	Status  string   `json:"status"`
}
