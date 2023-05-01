// for the speech to specch to occur the speech needs to be converted into a text.
// the text is used to query the wolfram alpha api.
//the the answer is a text/string and is written to an xml file
//the xml file is read to produce a wav file for the user


package main
import (
  "fmt"
  "errors"
  "io/ioutil"
  "net/http"
  "bytes"
  "strings"
  "log"
  "encoding/xml"
  "os"
)
//used to build the xml file
type speak struct {
	Version  string `xml:"version,attr"`
	Language string `xml:"xml:lang,attr"`
	Voice    voice  `xml:"Voice"`
}
//used to build the xml file
type voice struct {
	Language string `xml:"xml:lang,attr"`
	Name     string `xml:"name,attr"`
	Words    string `xml:",chardata"`
}


const (
  REGION = "uksouth"
  URI    = "https://" + REGION + ".stt.speech.microsoft.com/" +
           "speech/recognition/conversation/cognitiveservices/v1?" +
           "language=en-US"
  KEY    = "19c1cb3c0aa848608fed5a5a8a23d640"
)

func check( e error ) { if e != nil { panic( e ) } }

//text to speech function
func TextToSpeech( text []byte ) ( []byte, error ) {
	client   := &http.Client{}
	req, err := http.NewRequest( "POST", URI, bytes.NewBuffer( text ) )
	check( err )
  
	req.Header.Set( "Content-Type", "application/ssml+xml" )
	req.Header.Set( "Ocp-Apim-Subscription-Key", KEY )
	req.Header.Set( "X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm" )
  
	rsp, err2 := client.Do( req )
	check( err2 )
  
	defer rsp.Body.Close()
  
	if rsp.StatusCode == http.StatusOK {
	  body, err3 := ioutil.ReadAll( rsp.Body )
	  check( err3 )
	  return body,  nil
	} else {
	  return nil, errors.New( "cannot convert text to speech" )
	}
  } 

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


  speech, err1 := ioutil.ReadFile( "testwav.wav" ) //example wav file is inputed here 
  check( err1 )
  text, err2 := SpeechToText( speech ) //retrieved the text from wav file
  check( err2 )
  fmt.Println( text )

//text is the question string for the wolfram alpha short answers api
question = text
//string being modified to fit inside the url used to access api
question = strings.ReplaceAll(question, " ", "+")

//wolframalpha appid is already in the link and not stored in a variable
resp, getErr := http.Get("http://api.wolframalpha.com/v1/result?appid=H4RUYV-2WP4YG72XJ&i=" + question)
if getErr != nil {
	log.Fatal(getErr)
}
//body is the result string 
body, readErr := ioutil.ReadAll(resp.Body)
if readErr != nil {
	log.Fatal(readErr)
}
fmt.Println(string(body)) //retrieved the resulting string

//the string is to be written to the xml file
var answer string
answer = string(body)

x := &speak{Version: "1.0", Language: "en-US",
		Voice: voice{Language: "en-US", Name: "en-US-JennyNeural", Words: answer},//inserts answer
	}

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(x); err != nil {
		fmt.Printf("error: %v\n", err)
	}

	data, err := xml.MarshalIndent(x, " ", "  ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("test.xml", data, 0666)
	if err != nil {
		log.Fatal(err)
	}
//after making an xml file, use this file to read the speech
//after making an xml file, use this file to read the speech

text2, err := ioutil.ReadFile( "test.xml" ) //reads the file
check( err )
speech2, err2 := TextToSpeech( text2 )//converts to wav file
check( err2 )
err3 := ioutil.WriteFile( "speech.wav", speech2, 0644 )//outputs wav file
check( err3 )
}

