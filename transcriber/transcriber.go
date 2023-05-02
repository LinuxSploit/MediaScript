package transcriber

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	wav "github.com/go-audio/wav"
)

var (
	languages = map[string]string{"English": "en", "Arabic": "ar", "Armenian": "hy", "Azerbaijani": "az", "Basque": "eu", "Belarusian": "be", "Bengali": "bn", "Bulgarian": "bg", "Catalan": "ca", "Chinese": "zh", "Croatian": "hr", "Czech": "cs", "Danish": "da", "Dutch": "nl", "Estonian": "et", "Filipino": "tl", "Finnish": "fi", "French": "fr", "Galician": "gl", "Georgian": "ka", "German": "de", "Greek": "el", "Gujarati": "gu", "Hebrew": "iw", "Hindi": "hi", "Hungarian": "hu", "Icelandic": "is", "Indonesian": "id", "Irish": "ga", "Italian": "it", "Japanese": "ja", "Kannada": "kn", "Korean": "ko", "Latin": "la", "Latvian": "lv", "Lithuanian": "lt", "Macedonian": "mk", "Malay": "ms", "Maltese": "mt", "Norwegian": "no", "Persian": "fa", "Polish": "pl", "Portuguese": "pt", "Romanian": "ro", "Russian": "ru", "Serbian": "sr", "Slovak": "sk", "Slovenian": "sl", "Spanish": "es", "Swahili": "sw", "Swedish": "sv", "Tamil": "ta", "Telugu": "te", "Thai": "th", "Turkish": "tr", "Ukrainian": "uk", "Urdu": "ur", "Vietnamese": "vi", "Welsh": "cy", "Yiddish": "yi"}
)

type ScriptLine struct {
	Start time.Duration
	End   time.Duration
	Num   int
	Text  string
}

type Transcriber struct {
	tmpresamplewavfile string
	Duration           time.Duration
	Sample             []float32
	ScriptLines        []ScriptLine
}

func NewTranScriber() *Transcriber {
	return &Transcriber{}
}

// Reads Wav file resample it for whisper
func (x *Transcriber) ReadWav(wavpath string) error {
	// check if wav file exist
	_, err := os.Stat(wavpath)
	if err != nil {
		return err
	}
	// create a temp file
	whisperwavfile, err := ioutil.TempFile("", "sample*.wav")
	if err != nil {
		return err
	}
	defer whisperwavfile.Close()
	x.tmpresamplewavfile = whisperwavfile.Name()
	// conversion to 1 channel and whisper sample rate
	cmd := exec.Command("ffmpeg", "-i", wavpath, "-ac", "1", "-ar", "16000", x.tmpresamplewavfile, "-y")
	err = cmd.Run()
	if err != nil {
		return err
	}

	// decode wav file
	wavfile, err := os.Open(whisperwavfile.Name())
	if err != nil {
		return err
	}
	defer wavfile.Close()
	wavdec := wav.NewDecoder(wavfile)
	// reading PCM buffer
	if wavbuffer, err := wavdec.FullPCMBuffer(); err != nil {
		return err
	} else if wavdec.SampleRate != whisper.SampleRate {
		return errors.New("UNSUPPORTED SAMPLE RATE")
	} else if wavdec.NumChans != 1 {
		return errors.New("UNSUPPORTED NUMBER OF CHANNELS")
	} else {
		x.Sample = wavbuffer.AsFloat32Buffer().Data
	}

	// getting total duration
	dur, err := wavdec.Duration()
	if err != nil {
		return err
	}
	x.Duration = dur

	return nil

}

func (x *Transcriber) Transcribe(modelpath string, lang string) error {
	// check if model file exist
	modelpath, err := filepath.Abs(modelpath)
	if err != nil {
		return err
	}
	// Load the model
	model, err := whisper.New(modelpath)
	if err != nil {
		return err
	}
	defer model.Close()

	// Process samples
	context, err := model.NewContext()
	if err != nil {
		return err
	}

	log.Println(lang, languages[lang])
	context.SetLanguage(languages[lang])
	if err := context.Process(x.Sample, nil); err != nil {
		return err
	}

	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		x.ScriptLines = append(x.ScriptLines, ScriptLine{
			Start: segment.Start,
			End:   segment.End,
			Num:   segment.Num,
			Text:  segment.Text,
		})
	}

	os.Remove(x.tmpresamplewavfile)
	return nil
}
