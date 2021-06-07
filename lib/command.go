package lib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/adrg/xdg"
)

type Command int

const (
	Kubectl            Command = iota
	Helm               Command = iota
	defaultKubectlPath         = "/usr/bin/kubectl"
	defaultHelmPath            = "/usr/bin/helm"
)

type Ctx struct {
	cfg Config
}

func NewCtx(cmd Command) (*Ctx, error) {
	configDir := "kubectl"
	command := defaultKubectlPath
	if cmd == Helm {
		configDir = "helm"
		command = defaultHelmPath
	}
	configPath := path.Join(xdg.ConfigHome, configDir, "wrapper.yml")
	cb, err := ioutil.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if err != nil && os.IsNotExist(err) {
		return &Ctx{
			cfg: Config{
				Command: command,
			},
		}, nil
	}
	cfg, err := parseConfig(cb)
	if err != nil {
		return nil, err
	}
	if cfg.Command == "" {
		cfg.Command = command
	}
	return &Ctx{
		cfg: cfg,
	}, nil
}

func (c *Ctx) Run() error {
	var args []string
	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}
	if err := c.guardDelete(args); err != nil {
		return err
	}
	cmd := exec.Command(c.cfg.Command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Ctx) guardDelete(args []string) error {
	if len(c.cfg.ProtectedContexts) == 0 {
		return nil
	}
	currentK8sCtx, err := getK8sContext()
	if err != nil {
		return err
	}
	var ctxFound bool
	for _, context := range c.cfg.ProtectedContexts {
		if context == currentK8sCtx {
			ctxFound = true
			break
		}
	}
	if !ctxFound {
		return nil
	}
	var containsDelete bool
	for _, arg := range args {
		if arg == "delete" || arg == "uninstall" {
			containsDelete = true
			break
		}
	}
	if !containsDelete {
		return nil
	}
	fmt.Printf("\n\nWarning!!\n\nYou are about to perform a destructive operation in the protected context %q\n\n", currentK8sCtx)
	confirmMsg := "yes-i-am-really-really-sure"
	if c.cfg.ConfirmString != "" {
		confirmMsg = c.cfg.ConfirmString
	}
	for {
		fmt.Printf("enter %q > ", confirmMsg)
		reader := bufio.NewReader(os.Stdin)
		c, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if strings.ReplaceAll(c, "\n", "") == confirmMsg {
			break
		}
	}
	return nil
}

func getK8sContext() (string, error) {
	ctx, err := NewCtx(Kubectl)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(ctx.cfg.Command, "config", "current-context")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(out), "\n", ""), nil
}
