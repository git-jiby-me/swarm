package strategy

import (
	"github.com/git-jiby-me/swarm/cluster"
	"github.com/git-jiby-me/swarm/scheduler/node"
)

// WeightedNode represents a node in the cluster with a given weight, typically used for sorting
// purposes.
type weightedNode struct {
	Node *node.Node
	// Weight is the inherent value of this node.
	Weight int64
}

type weightedNodeList []*weightedNode

func (n weightedNodeList) Len() int {
	return len(n)
}

func (n weightedNodeList) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n weightedNodeList) Less(i, j int) bool {
	var (
		ip = n[i]
		jp = n[j]
	)

	// If the nodes have the same weight sort them out by number of containers.
	if ip.Weight == jp.Weight {
		return len(ip.Node.Containers) < len(jp.Node.Containers)
	}
	return ip.Weight < jp.Weight
}

func weighNodes(config *cluster.ContainerConfig, nodes []*node.Node, healthinessFactor int64) (weightedNodeList, error) {
	weightedNodes := weightedNodeList{}

	for _, node := range nodes {
		nodeMemory := node.TotalMemory
		nodeCpus := node.TotalCpus

		// Skip nodes that are smaller than the requested resources.
		if nodeMemory < int64(config.HostConfig.Memory) || nodeCpus < config.HostConfig.CPUShares {
			continue
		}

		var (
			cpuScore    int64 = 100
			memoryScore int64 = 100
		)

		if config.HostConfig.CPUShares > 0 {
			cpuScore = (node.UsedCpus + config.HostConfig.CPUShares) * 100 / nodeCpus
		}
		if config.HostConfig.Memory > 0 {
			memoryScore = (node.UsedMemory + config.HostConfig.Memory) * 100 / nodeMemory
		}

		if cpuScore <= 100 && memoryScore <= 100 {
			weightedNodes = append(weightedNodes, &weightedNode{Node: node, Weight: cpuScore + memoryScore + healthinessFactor*node.HealthIndicator})
		}
	}

	if len(weightedNodes) == 0 {
		return nil, ErrNoResourcesAvailable
	}

	return weightedNodes, nil
}
