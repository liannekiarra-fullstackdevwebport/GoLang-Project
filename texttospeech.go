//this file contains functionalities that are able to turn text to speech. 
//string type to xml 

package main
import (
  "bytes"
  "errors"
  "bufio"
  "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"encoding/xml"

)
//used for creating xml file
type speak struct {
	Version  string `xml:"version,attr"`
	Language string `xml:"xml:lang,attr"`
	Voice    voice  `xml:"Voice"`
}

type voice struct {
	Language string `xml:"xml:lang,attr"`
	Name     string `xml:"name,attr"`
	Words    string `xml:",chardata"`
}


const (
  REGION = "uksouth"
  URI    = "https://" + REGION + ".tts.speech.microsoft.com/" +
           "cognitiveservices/v1"
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
 //speech to text function
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
  //for the wolfram

  //user enters the question
  fmt.Printf("Enter your question : ")
	inputReader := bufio.NewReader(os.Stdin)
	input, _ := inputReader.ReadString('\n')
//modifiying the string to attach it to the url 
	input = strings.ReplaceAll(input, " ", "+")
	input = strings.ReplaceAll(input, "\n", "")

	resp, getErr := http.Get("http://api.wolframalpha.com/v1/result?appid=H4RUYV-2WP4YG72XJ&i=" + input)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
//retrieved response
//add to the xml file as requirement for microsoft AZURE
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
text, err := ioutil.ReadFile( "test.xml" ) //reads the file
check( err )
speech, err2 := TextToSpeech( text )//converts to wav file
check( err2 )
err3 := ioutil.WriteFile( "speech.wav", speech, 0644 )//outputs wav file
check( err3 )
}







