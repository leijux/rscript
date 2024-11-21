package parser

import (
	"errors"
	"fmt"
	"log/slog"
	"net/netip"
	"os"
	"path/filepath"

	"github.com/valyala/fasttemplate"
	"gopkg.in/yaml.v3"

	"github.com/leijux/rscript/internal/pkg/version"
)

type RScript struct {
	SchemaVersion uint              `yaml:"schema_version"`
	Variables     map[string]string `yaml:"variables,omitempty"`
	Commands      []string          `yaml:"commands"`
	DefaultConfig defaultConfig     `yaml:"default,omitempty"`
	Remotes       []Remote          `yaml:"remotes"`
}

type Remote struct {
	AddrPort netip.AddrPort `yaml:"-"`
	AddrStr  string         `yaml:"ip"`
	Username string         `yaml:"username,omitempty"`
	Password string         `yaml:"password,omitempty"`
}

type defaultConfig struct {
	Port     uint16 `yaml:"port,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

func rscriptFilePath(path string) (rscriptPath string) {
	rscriptPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		rscriptPath = path
		logFunc := slog.Error
		if errors.Is(err, os.ErrNotExist) {
			logFunc = slog.Debug
		}

		logFunc("evaluating config path: %s; using %q", err, rscriptPath)
		return ""
	}

	return rscriptPath
}

func validateConfig(script *RScript) (err error) {
	if script.SchemaVersion != version.LastSchemaVersion {
		return errors.New("invalid schema version")
	}

	if len(script.Remotes) == 0 {
		return errors.New("remotes can't be empty")
	}

	if len(script.Commands) == 0 {
		return errors.New("commands can't be empty")
	}

	dPort := script.DefaultConfig.Port
	dUser := script.DefaultConfig.Username
	dPass := script.DefaultConfig.Password

	m := make(map[netip.Addr]struct{}, len(script.Remotes))

	for i, r := range script.Remotes {
		if r.Username == "" {
			if dUser == "" {
				return fmt.Errorf("remote %s username can't be empty", r.AddrStr)
			}
			r.Username = dUser
		}

		if r.Password == "" {
			if dPass == "" {
				return fmt.Errorf("remote %s password can't be empty", r.AddrStr)
			}
			r.Password = dPass
		}

		addrPort, err := netip.ParseAddrPort(r.AddrStr)
		if err != nil { //如果解析失败 尝试解析地址加端口
			if dPort == 0 { //端口不能为空
				return fmt.Errorf("remote %s port can't be empty", r.AddrStr)
			}

			addr, err := netip.ParseAddr(r.AddrStr)
			if err != nil { // 解析地址失败
				return fmt.Errorf("remote %s parse addr err: %w", r.AddrStr, err)
			}
			// 地址加端口
			addrPort = netip.AddrPortFrom(addr, dPort)
		}

		addr := addrPort.Addr()

		if _, ok := m[addr]; ok { //判断地址是否重复
			return fmt.Errorf("remote ip repeat %v", addr)
		}

		r.AddrPort = addrPort
		script.Remotes[i] = r

		m[addr] = struct{}{}
	}

	return nil
}

func ParseWithBytes(fileData []byte) (*RScript, error) {
	diskConf := yobj{}

	err := yaml.Unmarshal(fileData, &diskConf)
	if err != nil {
		return nil, err
	}

	es, ok, err := fieldVal[map[string]any](diskConf, "variables")

	if ok && len(es) != 0 && err == nil {
		t := fasttemplate.New(string(fileData), "{{", "}}")

		fileData = []byte(t.ExecuteString(es))
	}

	config := &RScript{
		SchemaVersion: 0,
	}

	err = yaml.Unmarshal(fileData, &config)
	if err != nil {
		return nil, err
	}

	err = validateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func ParseWithPath(path string) (*RScript, error) {
	rscriptPath := rscriptFilePath(path)
	fileData, err := os.ReadFile(rscriptPath)
	if err != nil {
		return nil, err
	}

	return ParseWithBytes(fileData)
}
