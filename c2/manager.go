package c2

import (
	"sync"
	"time"
)

type Agent struct {
	Id       string
	Ip       string
	LastCall time.Time
	CmdQueue []string
}

type SafeAgentMap struct {
	mtx    sync.Mutex
	Agents map[string]*Agent // Key: Agent id, Value: pointer to Agent
}

// AgentMap is the global SafeAgentMap instance
var AgentMap SafeAgentMap = SafeAgentMap{Agents: make(map[string]*Agent)}

// SafeAgentMap.Add adds a new agent to the agent map
func (am *SafeAgentMap) Add(agent *Agent) {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	if _, exists := am.Agents[agent.Id]; !exists {
		am.Agents[agent.Id] = agent
	}
}

// SafeAgentMap.Get gets an agent from the agent map
func (am *SafeAgentMap) Get(agentId string) *Agent {
	am.mtx.Lock()
	defer am.mtx.Unlock()
	if agent, exists := am.Agents[agentId]; !exists {
		return agent
	}
	return nil
}

// SafeAgentMap.Queue queues a command for the agent
//func (am *SafeAgentMap) Queue()
