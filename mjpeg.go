package mjpeg

import (
  "fmt"
  "log"
  //"time"
  //"io"
  //"os"
  "os/exec"
  //"image"
  //"image/jpeg"
  "encoding/base64"
  //"bufio"
  //"lev/mjpeg"
)

type MjpegProcess struct {
  name string
  cmd string
  args []string
  onFrame func([]byte)

  closing chan chan error
}

func NewMjpegProcess(name string, cmd string, args []string, onFrame func([]byte))(*MjpegProcess){
  p := new(MjpegProcess)
  p.name = name
  p.cmd = cmd
  p.args = args
  p.onFrame = onFrame
  return p
}

func (this *MjpegProcess) Close() {
  errc := make(chan error)
  this.closing <- errc
}

func (this *MjpegProcess) RunProc() <-chan interface{} {

  //fmt.Printf("%s\n", arg)

  c := make(chan interface{})

  proc := exec.Command(this.cmd, this.args...)


  fmt.Printf("%s\n", proc.Args)
  // Create stdout, stderr streams of type io.Reader
  // stdout, err := proc.StdoutPipe()
  // checkError(err)
  // stderr, err := proc.StderrPipe()
  // checkError(err)

  // err = proc.Start()
  // checkError(err)

  // file, err := os.Create("mjpeg.txt")
  // checkError(err)

  // defer func() {
  //   if err := file.Close(); err != nil {
  //       panic(err)
  //   }
  // }()



  // w := bufio.NewWriter(file)

  // _, err := NewParser(stdout)
  // checkError(err)


  // Non-blockingly echo command output to terminal
  // go io.Copy(os.Stdout, stdout)
  // go io.Copy(os.Stderr, stderr)

  // Don't let main() exit before our command has finished running
  //defer proc.Wait()  // Doesn't block

  // I love Go's trivial concurrency :-D
  //fmt.Printf("Do other stuff here! No need to wait.\n\n")

  // r := bufio.NewReader(stdout)

  // handle := func(frame []byte){
  //   log.Printf("found image with size %d\n", len(frame))
  // }


  //go p.parse(handle)


  // for {
  //   img, err := jpeg.Decode(r)

  //   if err != nil && err != io.EOF{
  //     panic(err)
  //   }

  //   if img != nil{

      

  //     b64 := base64.NewEncoder(base64.StdEncoding, img)

  //     fmt.Printf("%s\n", b64)
  //   }
  // }

  // err = proc.Wait()
  // checkError(err)
  // if err = proc.Wait(); err != nil {
  //   log.Fatal(err)
  // }

  

  // fmt.Printf("%s\n", img)

  return c
}



func encode(bin []byte) []byte {
  e64 := base64.StdEncoding

  maxEncLen := e64.EncodedLen(len(bin))
  encBuf := make([]byte, maxEncLen)

  e64.Encode(encBuf, bin)
  return encBuf
}

func checkError(err error){
  if err != nil {
        log.Fatalf("Error: %s", err)
    }
}

func forgiveError(err error){
  if err != nil {
    log.Printf("Forgive: %s", err)
  }
}




// func main() {

//   cmd := "ffmpeg"

//   param := []string{"-loglevel", "debug", "-analyzeduration", "0", "-i", "rtmp://www.cr298.com:1954/?rtmp://172.16.11.104/alvs/ba_1/videoStream1 live=true subscribe=videoStream1 app=/?rtmp://172.16.11.104/alvs/ba_1/ playpath=videoStream1 buffer=300k conn=S:acs-inter-conn-mjpeg conn=S:asdasdokmkoihkklasdsdfbsdfvase", "-r", "8", "-s", "160x120", "-an", "-vcodec", "mjpeg", "-q:v", "4", "-f", "mjpeg", "-" } 

//   handle := func(frame []byte){
//     log.Printf("found image with size %d\n", len(frame))
//   }

//   p := NewMjpegProcess("ba_1", cmd, param, handle)

//   p.RunProc()

// }

