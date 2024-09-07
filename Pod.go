package kube

import (
	"context"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Pod struct {
	*core.Pod
	cmd    *exec.Cmd
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

func (p *Pod) Start(ctx *sync.Map) error {
	var container = p.Pod.Spec.Containers[0]
	p.SetCreationTimestamp(meta.Now())
	p.Status.Phase = core.PodPending
	for _, volume := range p.Pod.Spec.Volumes {
		os.Mkdir(volume.Name, 0755)
		if volume.HostPath != nil {
			syscall.Mount(volume.HostPath.Path, volume.Name, "", syscall.MS_BIND, "")
		} else if volume.ConfigMap != nil {
			var obj, _ = ctx.Load(volume.ConfigMap.Name)
			for k, v := range obj.(*core.ConfigMap).Data {
				os.WriteFile(filepath.Join(volume.Name, k), []byte(v), 0644)
			}
		}
	}
	switch p.Pod.Spec.RestartPolicy {
	default:
		p.cmd = exec.CommandContext(context.Background(), container.Image, container.Args...)
		p.Stdin, _ = p.cmd.StdinPipe()
		p.Stdout, _ = p.cmd.StdoutPipe()
		p.Stderr, _ = p.cmd.StderrPipe()
		p.cmd.Start()
		p.Status.Phase = core.PodRunning
		if probe := container.LivenessProbe; probe != nil {
			go func() {
				for p.cmd.ProcessState == nil && exec.Command(probe.Exec.Command[0], probe.Exec.Command[1:]...).Run() == nil {
					time.Sleep(time.Duration(probe.PeriodSeconds) * time.Second)
				}
				p.Stop(ctx)
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
func (p *Pod) Stop(ctx *sync.Map) error {
	return p.cmd.Cancel()
}
