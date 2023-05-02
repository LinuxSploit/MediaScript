# MediaScript
MediaScript is an application that allows you to convert any media file, such as audio or video, into a subtitle file in the .srt format. The application uses OpenAI's Whisper model, combined with the whisper.cpp C++ library (which has been wrapped with a Go language binding), to perform speech recognition and generate the subtitle text.

The application provides a simple, user-friendly graphical user interface (GUI) built with the fyne.io toolkit, making it easy to use for both beginners and advanced users.

To use MediaScript, simply select the media file you want to convert using the file browser in the application's main window. The application will automatically start processing the file, generating a subtitle file in the .srt format.

Once the subtitle file is generated, you can save it to your local drive, and use it to display subtitles while playing the media file in any media player that supports the .srt subtitle format.

Overall, MediaScript provides an easy and convenient way to create subtitle files for any media file, without the need for manual transcription or editing.

![Screenshot](https://github.com/LinuxSploit/MediaScript/raw/main/screenshot/screenshot.gif)
