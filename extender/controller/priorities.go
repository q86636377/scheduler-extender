package controller

import (
	"log"
	"math/rand"

	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

const (  //显示取得了多少分
	luckyPrioMsg = "pod %v/%v 取得 %v 分\n"
)

func prioritize(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	pod := args.Pod
	nodes := args.Nodes.Items

	hostPriorityList := make(schedulerapi.HostPriorityList, len(nodes))
	for i, node := range nodes {
		score := rand.Intn(schedulerapi.MaxPriority + 1)
		log.Printf(luckyPrioMsg, pod.Name, pod.Namespace, score)
		hostPriorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: score,
		}
	}

	return &hostPriorityList
}