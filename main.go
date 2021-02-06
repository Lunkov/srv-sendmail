package main

import (
  "fmt"
  "flag"
  "net"
  "net/smtp"
  "io/ioutil"
  "path/filepath"
  "gopkg.in/yaml.v2"
  "github.com/golang/glog"
  
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
  
  "crypto/tls"

  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/grpc-bpmn"
)

type SMTPInfo struct {
  EnableTLS       bool    `yaml:"enable_starttls"`
  ConnectStr      string  `yaml:"connect"`
  Address         string  `yaml:"address"`
  Port            int32   `yaml:"port"`
  Domain          string  `yaml:"domain"`
  Authentication  string  `yaml:"authentication"`
  UserName        string  `yaml:"user_name"`
  UserLogin       string  `yaml:"user_login"`
  Password        string  `yaml:"password"`
  Secret          string  `yaml:"secret"` // MD5
  Auth            smtp.Auth
  TLS             tls.Config
}

type BPMNInfo struct {
  ConnectStr      string  `yaml:"connect"`
}

type ConfigInfo struct {
  ConfigPath      string
  SMTP            map[string]SMTPInfo  `yaml:"smtp_settings"`
  BPMN            BPMNInfo
}

var globConf ConfigInfo

type BPMNJobService struct{}

func (s *BPMNJobService) CallFunction(ctx context.Context, in *srv_bpmn.RPCBPMNJob) (*srv_bpmn.RPCBPMNJobResponse, error) {
  ok := sendMail(&globConf.SMTP, &in.Parameters)
	return &srv_bpmn.RPCBPMNJobResponse{BpmnProcessId: in.BpmnProcessId, Ok: ok}, nil
}

func (c *SMTPInfo) expand() {
  if c.ConnectStr == "" {
    c.ConnectStr = fmt.Sprintf("%s:%d", c.Address, c.Port)
  }
  c.Auth = nil
  if c.Authentication == "plain" {
    c.Auth = smtp.PlainAuth("", c.UserLogin, c.Password, c.Address)
  }
  if c.Authentication == "md5" {
    c.Auth = smtp.CRAMMD5Auth(c.UserLogin, c.Secret)
  }
  if c.EnableTLS {
    c.TLS = tls.Config{ InsecureSkipVerify: false, ServerName: c.Address }
  }
}

func loadConfig(filename string) ConfigInfo {
  var err error
  cfg := ConfigInfo{}

  if !env.WaitFile(filename, 300) {
    return cfg
  }

  yamlFile, err := ioutil.ReadFile(filename)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s)  #%v ", filename, err)
    return cfg
  }
  err = yaml.Unmarshal(yamlFile, &cfg)
  if err != nil {
    glog.Errorf("ERR: yamlFile(%s): YAML: %v", filename, err)
  }
  if cfg.ConfigPath == "" {
    cfg.ConfigPath = filepath.Dir(filename)
  }
  for i, sm := range cfg.SMTP {
    sm.expand()
    cfg.SMTP[i] = sm
  }
  return cfg
}

func main() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	// flag.Set("v", "9")
  configPath := flag.String("config_path", "./etc/", "Config path")
  flag.Parse()
    
  globConf = loadConfig(*configPath + "config.yaml")

  glog.Infof("LOG: Start gRPC SendMail server")

	lis, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		glog.Fatalf("FAIL: Can not listen the port：%v", err)
	}

	s := grpc.NewServer()

  srv_bpmn.RegisterBPMNJobServer(s, &BPMNJobService{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		glog.Fatalf("ERR: Can not provide service：%v", err)
	}
}
