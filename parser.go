package mjpeg

import (
  "io"
  "log"
  "bufio"
  "fmt"
)

type FormatError string

func (e FormatError) Error() string { return string(e) }

const DEF_BUF_SZ = 4096
const DEF_BLK_SZ = 256

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
    newBuf := make([]byte, len(this.buf), dubSize)
    copy(newBuf, this.buf)
    this.buf = newBuf
  }

  //fmt.Printf("buf len %d, cap %d\n", len(this.buf), cap(this.buf))

  chunk, err := io.ReadFull(this.r, this.buf[this.read:bound])
  if err != nil {
    log.Printf("goes here %s\n", err)
    return nil, err
  }
  this.read += chunk
  //grow the slice by read bytes
  this.buf = this.buf[:this.read]

  //fmt.Printf("buf len %d, cap %d\n", len(this.buf), cap(this.buf))

  imgStart := IndexOfBytes(this.buf, SOI)
  imgEnd := IndexOfBytes(this.buf, EOI)
  if imgStart < 0 {
    return nil, FormatError("SOI not located")
  }
  if imgEnd < 0 {
    return nil, FormatError("EOI not located")
  }
  imgEnd = imgEnd + 2 //add last two bytes of EOI
  fmt.Printf("start %d, end %d\n", imgStart, imgEnd)

  imgBuf := append([]byte{}, this.buf[imgStart:imgEnd]...)
  
  rest := len(this.buf) - imgEnd
  copy(this.buf[0:], this.buf[imgEnd:])
  this.buf = this.buf[:rest]
  this.read = rest
  return imgBuf, nil
}


func IndexOfBytes(buf []byte, search []byte) (int){
  start := 0
  m := len(search)
  n := len(buf) - m
  matched := 0
  for ; start < n; start++ {
    matched = 0
    for next := 0; next < m ; next++ {
      if buf[start+next] == search[next] {
        matched++
      } else {
        break
      }
    }
    if matched == m {
      return start
    }
  }
  return -1
}



