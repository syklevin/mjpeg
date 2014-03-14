package mjpeg

import (
  "testing"
  "bufio"
  "os"
)

func TestParse(t *testing.T) {
  f, err := os.OpenFile("mjpeg_bin.txt", os.O_RDONLY, 0644) 
  if err != nil { 
    t.Errorf("can not open file %s", err)
    return
  }
  defer f.Close()

  p, err := NewParser(bufio.NewReader(f));
  if err != nil {
    t.Errorf("can not new parser %s", err)
  }

  handle := func(frame []byte){
    t.Logf("%X\n", frame[:20])
    t.Logf("-----------------------------------\n")
  }

  p.parse(handle)

  t.Logf("found images %d\n", p.cnt)

}
