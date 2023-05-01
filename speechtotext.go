package main
import (
  "fmt"
  "errors"
  "io/ioutil"
  "net/http"
  "bytes"
  "strings"
  "log"
)

//utilising microsoft azure

const (
  REGION = "uksouth"
  URI    = "https://" + REGION + ".stt.speech.microsoft.com/" +
           "speech/recognition/conversation/cognitiveservices/v1?" +
           "language=en-US"
  KEY    = "19c1cb3c0aa848608fed5a5a8a23d640"
)

func check( e error ) { if e != nil { panic( e ) } }

//speech to texxt function

func SpeechToText( speech []byte ) ( string, error ) {
  client   := &http.Client{}
  req, err := http.NewRequest( "POST", URI, bytes.NewReader( speech ) )
  check( err )

  req.Header.Set( "Content-Type",
                  "audio/wav;codecs=audio/pcm;samplerate=16000" )
  req.Header.Set( "Ocp-Apim-Subscription-Key", KEY )

  rsp, err2 := client.Do( req )
  check( err2 )

  defer rsp.Body.Close()

  if rsp.StatusCode == http.StatusOK {
    body, err3 := ioutil.ReadAll( rsp.Body )
    check( err3 )
    return string( body ),  nil
  } else {
    return "", errors.New( "cannot convert to speech to text" ) 
  }
}

func main() {

//declaring question text 
var question string
  speech, err1 := ioutil.ReadFile( "speech.wav" ) //example wav file is inputed here 
  check( err1 )
  text, err2 := SpeechToText( speech ) //retrieved the text from wav file
  check( err2 )
  fmt.Println( text )

//text is the question string for the wolfram alpha short answers api
question = text
//string being modified to fit inside the url used to access api
question = strings.ReplaceAll(question, " ", "+")

resp, getErr := http.Get("http://api.wolframalpha.com/v1/result?appid=H4RUYV-2WP4YG72XJ&i=" + question)
if getErr != nil {
	log.Fatal(getErr)
}
//body is the result string 
body, readErr := ioutil.ReadAll(resp.Body)
if readErr != nil {
	log.Fatal(readErr)
}
fmt.Println(string(body))

}



