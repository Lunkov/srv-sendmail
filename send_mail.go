package main

import (
  "path/filepath"
  "io/ioutil"
  "bytes"
  "strings"
  "github.com/golang/glog"
)

func sendMail(smtp_settings *map[string]SMTPInfo, prop *map[string]string) bool {
  if glog.V(2) {
    glog.Infof("LOG: SEND MAIL: run prepare...\n")
  }
  
  var ok           bool
  var mailFromName string
  var mailFrom     SMTPInfo
  var mailTo       string
  var mailSubject  string
  var mailBody     string
  var attachments  string
  
  mailFromName, ok = (*prop)["SEND_MAIL_FROM"]
  if !ok {
    glog.Errorf("ERR: SEND MAIL: 'SEND_MAIL_FROM' don`t set\n")
    return false
  }
  
  mailFrom, ok = (*smtp_settings)[mailFromName]
  if !ok {
    glog.Errorf("ERR: SEND MAIL: 'SEND_FROM(%s)' not found", mailFromName)
    return false
  }
  mailTo, ok = (*prop)["SEND_MAIL_TO"]
  if !ok {
    glog.Errorf("ERR: SEND MAIL: 'SEND_MAIL_TO' don`t set\n")
    return false
  }
  mailSubject, ok = (*prop)["SEND_MAIL_SUBJECT"]
  if !ok {
    glog.Errorf("ERR: SEND MAIL: 'SEND_MAIL_SUBJECT' don`t set\n")
    return false
  }
  // INIT MAIL
  mail := NewMail(&mailFrom) // mailFrom.ConnectStr, mailFrom.Auth, TLS)
  if mail == nil {
    glog.Errorf("ERR: CONNECT SMTP")
    return false
  }

  mail.From(mailFrom.UserLogin)
  mail.FromName(mailFrom.UserName)
  
  mail.Subject(mailSubject)

  mailBody, ok = (*prop)["SEND_MAIL_BODY_PLAIN"]
  if ok {
    mail.Plain().Set(mailBody)
  }
  mailBody, ok = (*prop)["SEND_MAIL_BODY_HTML"]
  if ok {
    mail.HTML().Set(mailBody)
  }
  
  attachments, ok = (*prop)["SEND_MAIL_ATTACMENTS"]
  if ok {
    arAttachments := strings.Split(attachments + ";", ";")
    if glog.V(9) {
      glog.Infof("DBG: SEND MAIL: READ FILES: %v", arAttachments)
    }
    for _, filename := range arAttachments {
      if filename != "" {
        if glog.V(9) {
          glog.Infof("DBG: SEND MAIL: READ FILE: %v", filename)
        }
        rawBytes, err := ioutil.ReadFile(filename)
        if err != nil {
          glog.Errorf("ERR: SEND MAIL: CANN`T READ FILE(%s): %v", filename, err)
          return false
        }

        mail.Attach(filepath.Base(filename), bytes.NewBuffer(rawBytes))
      }
    }
  }
  
  arMailTo := strings.Split(mailTo + ";", ";")
  result := true  
  for _, mail2 := range arMailTo {
    if mail2 != "" {
      if glog.V(2) {
        glog.Infof("LOG: SENDING MAIL(%s: %s) ...", mail2, mailSubject)
      }
      mail.To(mail2)
      
      glog.Infof("DBG: SENDING MAIL## %s", mail.String())
      
      
      if err := mail.Send(); err != nil {
        glog.Errorf("ERR: SEND MAIL(%s: %s): %v", mail2, mailSubject, err)
        result = false
      } else {
        if glog.V(2) {
          glog.Infof("LOG: SENDED MAIL(%s: %s)", mail2, mailSubject)
        }
      }
    }
  }
  return result
}

