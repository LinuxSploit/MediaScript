# MediaScript
MediaScript is an application that allows you to convert any media file, such as audio or video, into a subtitle file in the .srt format. The application uses OpenAI's Whisper model, combined with the whisper.cpp C++ library (which has been wrapped with a Go language binding), to perform speech recognition and generate the subtitle text.

The application provides a simple, user-friendly graphical user interface (GUI) built with the fyne.io toolkit, making it easy to use for both beginners and advanced users.

To use MediaScript, simply select the media file you want to convert using the file browser in the application's main window. The application will automatically start processing the file, generating a subtitle file in the .srt format.

Once the subtitle file is generated, you can save it to your local drive, and use it to display subtitles while playing the media file in any media player that supports the .srt subtitle format.

Overall, MediaScript provides an easy and convenient way to create subtitle files for any media file, without the need for manual transcription or editing.

## Demo:
![Screenshot](https://github.com/LinuxSploit/MediaScript/raw/main/screenshot/screenshot.gif)

# Build
```bash
git clone --recursive https://github.com/LinuxSploit/MediaScript.git
cd MediaScript/
```

Compile libwhisper.a (you can use make whisper in the bindings/go directory);
```bash
cd ./whisper.cpp/bindings/go/
make whisper
cd ../../..
```

 You can download and run the other models as follows:
```bash
# ./models/download_model.sh [Model]
./models/download_model.sh base
```
## Available models

| Model     | Disk   | Mem     | SHA                                        |
| ---       | ---    | ---     | ---                                        |
| tiny      |  75 MB | ~390 MB | `bd577a113a864445d4c299885e0cb97d4ba92b5f` |
| tiny.en   |  75 MB | ~390 MB | `c78c86eb1a8faa21b369bcd33207cc90d64ae9df` |
| base      | 142 MB | ~500 MB | `465707469ff3a37a2b9b8d8f89f2f99de7299dac` |
| base.en   | 142 MB | ~500 MB | `137c40403d78fd54d454da0f9bd998f78703390c` |
| small     | 466 MB | ~1.0 GB | `55356645c2b361a969dfd0ef2c5a50d530afd8d5` |
| small.en  | 466 MB | ~1.0 GB | `db8a495a91d927739e50b3fc1cc4c6b8f6c2d022` |
| medium    | 1.5 GB | ~2.6 GB | `fd9727b6e1217c2f614f9b698455c4ffd82463b4` |
| medium.en | 1.5 GB | ~2.6 GB | `8c30f0e44ce9560643ebd10bbe50cd20eafd3723` |
| large-v1  | 2.9 GB | ~4.7 GB | `b1caaf735c4cc1429223d5a74f0f4d0b9b59a299` |
| large     | 2.9 GB | ~4.7 GB | `0f4c8e34f21cf1a914c59d8b3ce882345ad349d6` |


Link MediaScript against whisper by setting the environment variables C_INCLUDE_PATH and LIBRARY_PATH to point to the whisper.h file directory and libwhisper.a file directory respectively.
```bash
. ./export.sh
```
compile binary
```bash
go build .
```

## Supported Lanuage:
- English
- Arabic
- Armenian
- Azerbaijani
- Basque
- Belarusian
- Bengali
- Bulgarian
- Catalan
- Chinese
- Croatian
- Czech
- Danish
- Dutch
- Estonian
- Filipino
- Finnish
- French
- Galician
- Georgian
- German
- Greek
- Gujarati
- Hebrew
- Hindi
- Hungarian
- Icelandic
- Indonesian
- Irish
- Italian
- Japanese
- Kannada
- Korean
- Latin
- Latvian
- Lithuanian
- Macedonian
- Malay
- Maltese
- Norwegian
- Persian
- Polish
- Portuguese
- Romanian
- Russian
- Serbian
- Slovak
- Slovenian
- Spanish
- Swahili
- Swedish
- Tamil
- Telugu
- Thai
- Turkish
- Ukrainian
- Urdu
- Vietnamese
- Welsh
- Yiddish