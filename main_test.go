package main

import (
  "flag"
  "testing"
  "net"
  "context"
  
  "github.com/golang/glog"
  
  "google.golang.org/grpc"
  "google.golang.org/grpc/test/bufconn"
  
  "github.com/Lunkov/lib-env"
  "github.com/Lunkov/grpc-bpmn"
)

/////////////////////////
// TESTS
/////////////////////////

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
    lis = bufconn.Listen(bufSize)
    s := grpc.NewServer()
    srv_bpmn.RegisterBPMNJobServer(s, &BPMNJobService{})
    go func() {
        if err := s.Serve(lis); err != nil {
          glog.Errorf("ERR: Server exited with error: %v", err)
        }
    }()
}

func bufDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial()
}

/////
func TestGRPCSendMail(t *testing.T) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", ".")
	flag.Set("v", "9")
	flag.Parse()
  
  globConf = loadConfig("./etc4test/config.yaml")
  
  ctx := context.Background()
  conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
  if err != nil {
      t.Fatalf("Failed to dial bufnet: %v", err)
  }
  defer conn.Close()
  client := srv_bpmn.NewBPMNJobClient(conn)
  
  prop := map[string]string{"SEND_MAIL_BODY_HTML": "Invoice in attachment", "SEND_MAIL_ATTACMENTS": "./storage/invoice.ru.1.html",
              "ACCOUNT_TO_NAME": "ООО \"Получатель\"",
              "ACCOUNT_TO_INDEX": "127282",
              "ACCOUNT_TO_CITY": "Москва",
              "ACCOUNT_FROM_CITY": "Москва",
              "PAYMENT_SUM_WITH_VAT": "105.23",
              "SEND_MAIL_FROM": "notify_mail",
              "SEND_MAIL_TO": "s_lunkov@mail.ru",
              "SEND_MAIL_SUBJECT": "Счёт №1",
             }
  
  resp, err := client.CallFunction(ctx, &srv_bpmn.RPCBPMNJob{Parameters: prop})
  if err != nil {
    t.Fatalf("SayHello failed: %v", err)
  }
  if resp.Ok != true {
    t.Error(
      "For", "ERR: gRPC template",
      "expected", true,
      "got", resp.Ok,
    )
  }
  glog.Infof("LOG: Response: %+v", resp)

}
