package mjpeg

import (
  "io"
  "log"
  "bufio"
)

type FormatError string

func (e FormatError) Error() string { return string(e) }

const DEF_BUF_SZ = 4096
const DEF_BLK_SZ = 64

var SOI = []byte { 0xff, 0xd8 }
var EOI = []byte { 0xff, 0xd9 }

type Reader interface {
  io.Reader
  ReadByte() (c byte, err error)
}

type Parser struct {
  r Reader
  buf []byte
  read int
  cnt int64
}

func NewParser(r io.Reader) (*Parser, error) {
  p := new(Parser)
  if rr, ok := r.(Reader); ok {
    p.r = rr
  } else {
    p.r = bufio.NewReader(r)
  }
  p.buf = make([]byte, DEF_BUF_SZ)
  p.read = 0
  p.cnt = 0
  return p, nil
}

// func (this *Parser) parse()(c chan interface{}){

//   c := make(chan interface{})

//   var err error

//   for {
//     log.Printf("start parsing")

//     frame, err := this.parseFrame()

//     if err == io.EOF {
//       c <- err
//     }

//     if frame != nil {
//       c <- frame
//     }

//     select t := c.(type) {
//       case errc := <- this.closing:
//         errc <- err
//         return
//     }

//     if err == io.EOF {
//       c <- err
//     }

//     // if err != nil {
//     //   log.Print(err)
//     // }

//     if data != nil {
//       c <- data
//     }

//     //time.Sleep(100 * time.Millisecond)
//   }
// }

func (this *Parser) parse(handle func([]byte)){

  for {
    frame, err := this.parseFrame()
    if err == io.EOF {
      break
    }
    // if err != nil {
    //   t.Logf("can not parse frame %s", err)
    // }
    if frame != nil {
      this.cnt++
      handle(frame)
    }
  }
}


func (this *Parser) parseFrame() ([]byte, error){
  bound := this.read + DEF_BLK_SZ
  if bound > cap(this.buf) {
    // allocate double what's needed, for future growth.
    dubSize := (cap(this.buf))*2
    newBuf := make([]byte, dubSize, dubSize)
    copy(newBuf, this.buf)
    this.buf = newBuf
  }
  chunk, err := io.ReadFull(this.r, this.buf[this.read:bound])
  if err != nil {
    log.Printf("goes here %s\n", err)
    return nil, err
  }
  this.read += chunk
  imgStart := IndexOfBytes(this.buf, SOI, 0)
  imgEnd := IndexOfBytes(this.buf, EOI, 0)
  if imgStart < 0 {
    return nil, FormatError("SOI not located")
  }
  if imgEnd < 0 {
    return nil, FormatError("EOI not located")
  }
  imgEnd = imgEnd + 2 //add last two bytes of EOI
  imgSize := imgEnd - imgStart
  imgBuf := make([]byte, imgSize)
  copy(imgBuf, this.buf[imgStart:imgEnd])
  newBuf := make([]byte, cap(this.buf))
  copy(newBuf, this.buf[imgEnd:])
  this.buf = newBuf
  this.read = 0
  return imgBuf, nil
}


func IndexOfBytes(buf []byte, search []byte, start int) (int){
  if start < 0 { start = 0 }
  m := len(search)
  n := len(buf) - m
  for ; start < n; start++ {
    if buf[start] == search[0] {
      next := 1
      for ; next < m ; next++ {
        if buf[start+next] != search[next] {
          break
        }
      }

      if next == m {
        return start
      }
    }
  }
  return -1
}

