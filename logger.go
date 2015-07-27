package main

import (
  "log"
  "os"
  "io"
  // "bufio"
)

const (
  DEBUG = 1 << iota
  INFO
  WARNING
  ERROR
)

type Logger struct {
  level int
  path *string
  descriptor *os.File
  debugLog   *log.Logger //1
  infoLog    *log.Logger //2
  warningLog *log.Logger //4
  errorLog   *log.Logger //5
}

func (self *Logger) Init() {
  var (
    writer io.Writer
    err error
  )

  if self.path == nil || *self.path == "" {
    writer = os.Stdout
  }else{
    self.descriptor, err = os.OpenFile(*self.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
    if err != nil {
      panic(err)
    }

    writer = io.MultiWriter(os.Stdout, self.descriptor)
  }

  self.debugLog = log.New(writer, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
  self.infoLog = log.New(writer, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
  self.warningLog = log.New(writer, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
  self.errorLog = log.New(writer, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (self *Logger) Debug (a ...interface{}) {
    if self.level & DEBUG > 0 {
      self.debugLog.Println(a...)
    }
}

func (self *Logger) Info (a ...interface{}) {
    if self.level & INFO > 0 {
      self.infoLog.Println(a...)
    }
}

func (self *Logger) Warning (a ...interface{}) {
    if self.level & WARNING > 0 {
      self.warningLog.Println(a...)
    }
}

func (self *Logger) Error (a ...interface{}) {
    if self.level & ERROR > 0 {
      self.errorLog.Println(a...)
    }
}

func (self *Logger) Close() {
  if self.descriptor != nil {
    self.descriptor.Close()
  }
}
