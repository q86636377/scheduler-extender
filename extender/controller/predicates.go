package controller

import (
	"log"
	"math/rand"
	"strings"

	"k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

const (
	LuckyPred        = "Lucky"
	LuckyPredFailMsg = "Sorry, you're not lucky"
)

var predicatesFuncs = map[string]FitPredicate{
	LuckyPred: LuckyPredicate,
}

type FitPredicate func(pod *v1.Pod, node v1.Node) (bool, []string, error)

var predicatesSorted = []string{LuckyPred}

func filter(args schedulerapi.ExtenderArgs) *schedulerapi.ExtenderFilterResult {
	var filteredNodes []v1.Node
	failedNodes := make(schedulerapi.FailedNodesMap)
	pod := args.Pod
	for _, node := range args.Nodes.Items {
		fits, failReasons, _ := podFitsOnNode(pod, node)
		if fits {
			filteredNodes = append(filteredNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",")
		}
	}

	result := schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       "",
	}

	return &result
}

func podFitsOnNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	fits := true
	var failReasons []string
	for _, predicateKey := range predicatesSorted {
		fit, failures, err := predicatesFuncs[predicateKey](pod, node)
		if err != nil {
			return false, nil, err
		}
		fits = fits && fit
		failReasons = append(failReasons, failures...)
	}
	return fits, failReasons, nil
}

func LuckyPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	lucky := rand.Intn(2) == 1   //若为奇数则匹配成功
	if lucky {
		log.Printf("pod %v/%v 调度节点 %v 成功\n", pod.Name, pod.Namespace, node.Name)
		return true, nil, nil
	}
	log.Printf("pod %v/%v 匹配节点 %v 失败\n", pod.Name, pod.Namespace, node.Name)
	return false, []string{LuckyPredFailMsg}, nil
}