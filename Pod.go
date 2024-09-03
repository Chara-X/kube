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

func (p *Pod) Run(ctx context.Context) error {
	var container = p.Pod.Spec.Containers[0]
	if probe := container.LivenessProbe; probe != nil {
		go func() {
			for p.cmd.ProcessState == nil {
				if err := exec.CommandContext(ctx, probe.Exec.Command[0], probe.Exec.Command[0:]...).Run(); err != nil {
					p.cmd.Cancel()
					break
				}
				time.Sleep(time.Duration(probe.PeriodSeconds) * time.Second)
			}
		}()
	}
	for _, volume := range p.Pod.Spec.Volumes {
		os.Mkdir(volume.Name, 0755)
		if volume.HostPath != nil {
			syscall.Mount(volume.HostPath.Path, volume.Name, "", syscall.MS_BIND, "")
		} else if volume.NFS != nil {
			syscall.Mount(volume.NFS.Server, volume.Name, "nfs", 0, "")
		} else if volume.ConfigMap != nil {
			for k, v := range ctx.Value("resources").(Store)[volume.ConfigMap.Name].(*core.ConfigMap).Data {
				os.WriteFile(filepath.Join(volume.Name, k), []byte(v), 0644)
			}
		}
	}
	switch p.Pod.Spec.RestartPolicy {
	default:
		p.cmd = exec.CommandContext(ctx, container.Image, container.Command...)
		p.Status.Phase = core.PodRunning
		var err = p.cmd.Run()
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
