package ffmpeg

import (
	"bytes"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go-server-base/global"
	"os"
)

type FfmpegUtils struct{}

func (f *FfmpegUtils) VideoCut(inFilePath string, outputFilePath string, startTime uint32, endTime uint32) {
	err := ffmpeg.Input("rtmp://192.168.127.12:1935/livetv/hunantv").Output("output.mp4").Run()
	if err != nil {
		return
	}

}

func (f *FfmpegUtils) ReadTimePositionAsJpeg(inFilePath string, seconds int, outputFilePath string) error {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFilePath, ffmpeg.KwArgs{"ss": seconds}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		global.LOG.Error(err.Error())
	}
	err = imaging.Save(img, outputFilePath)
	if err != nil {
		global.LOG.Error(err.Error())
	}
	return err
}
