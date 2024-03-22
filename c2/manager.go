package c2

import (
	"fmt"
	"sync"
	"time"
)

type Agent struct {
	Id       string
	Ip       string
	LastCall time.Time
	CmdQueue [][]string
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
	if agent, exists := am.Agents[agentId]; exists {
		return agent
	}
	return nil
}

// SafeAgentMap.Enqueue queues a command for the agent
func (am *SafeAgentMap) Enqueue(agentId string, cmd []string) error {
	agent := am.Get(agentId)
	if agent == nil {
		return fmt.Errorf("agent '%s' doesnt exist", agentId)
	}
	am.mtx.Lock()
	defer am.mtx.Unlock()
	agent.CmdQueue = append(agent.CmdQueue, cmd)
	return nil
}

// SafeAgentMap.Dequeue dequeues a command from the command queue
func (am *SafeAgentMap) Dequeue(agentId string) ([]string, error) {
	agent := am.Get(agentId)
	if agent == nil {
		return nil, fmt.Errorf("agent '%s' doesnt exist", agentId)
	}
	am.mtx.Lock()
	defer am.mtx.Unlock()
	if len(agent.CmdQueue) < 1 {
		return nil, fmt.Errorf("agent '%s' has no queued  commands", agentId)
	}
	cmd := agent.CmdQueue[0]
	agent.CmdQueue = agent.CmdQueue[1:]
	return cmd, nil
}
