package kube

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	core "k8s.io/api/core/v1"
)

type Pod struct {
	*core.Pod
	cmd *exec.Cmd
}

func (p *Pod) Start(ctx *Context) error {
	var container = p.Pod.Spec.Containers[0]
	for _, volume := range p.Pod.Spec.Volumes {
		os.Mkdir(volume.Name, 0755)
		if volume.HostPath != nil {
			syscall.Mount(volume.HostPath.Path, volume.Name, "", syscall.MS_BIND, "")
		} else if volume.NFS != nil {
			syscall.Mount(volume.NFS.Server, volume.Name, "nfs", 0, "")
		} else if volume.ConfigMap != nil {
			for k, v := range ctx.Get(volume.ConfigMap.Name).(*core.ConfigMap).Data {
				os.WriteFile(filepath.Join(volume.Name, k), []byte(v), 0644)
			}
		}
	}
	switch p.Pod.Spec.RestartPolicy {
	default:
		p.cmd = exec.CommandContext(context.Background(), container.Image, container.Command...)
		p.cmd.Start()
		p.Status.Phase = core.PodRunning
		if probe := container.LivenessProbe; probe != nil {
			go func() {
				for p.cmd.ProcessState == nil && exec.Command(probe.Exec.Command[0], probe.Exec.Command[0:]...).Run() == nil {
					time.Sleep(time.Duration(probe.PeriodSeconds) * time.Second)
				}
				p.Stop()
			}()
		}
		var err = p.cmd.Wait()
		for _, volume := range p.Pod.Spec.Volumes {
			os.RemoveAll(volume.Name)
		}
		if err == nil {
			p.Status.Phase = core.PodSucceeded
		} else {
			p.Status.Phase = core.PodFailed
		}
		return err
	}
}
func (p *Pod) Stop() error {
	return p.cmd.Cancel()
}
