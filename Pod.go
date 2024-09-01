package kube

import (
	"context"
	"os/exec"
	"time"

	core "k8s.io/api/core/v1"
)

type Pod struct {
	*core.Pod
	cmd *exec.Cmd
}

func (p *Pod) Start() error {
	var container = p.Pod.Spec.Containers[0]
	if probe := container.LivenessProbe; probe != nil {
		go func() {
			for p.cmd.ProcessState == nil {
				if err := exec.Command(probe.Exec.Command[0], probe.Exec.Command[0:]...).Run(); err != nil {
					p.Stop()
					break
				}
				time.Sleep(time.Duration(probe.PeriodSeconds) * time.Second)
			}
		}()
	}
	switch p.Pod.Spec.RestartPolicy {
	default:
		p.cmd = exec.CommandContext(context.Background(), container.Image, container.Command...)
		p.cmd.Start()
		p.Status.Phase = core.PodRunning
		var err = p.cmd.Wait()
		if err == nil {
			p.Status.Phase = core.PodSucceeded
		} else {
			p.Status.Phase = core.PodFailed
		}
		return err
	}
}
func (p *Pod) Stop() error {
	p.Pod.Status.Phase = core.PodFailed
	return p.cmd.Cancel()
}
