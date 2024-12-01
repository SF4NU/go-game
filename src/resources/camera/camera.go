package camera

import "github.com/setanarut/kamera/v2"

type Camera struct {
	Cam *kamera.Camera
}

func New(w, h, x, y int) *Camera {
	newCam := &Camera{
		Cam: kamera.NewCamera(float64(x), float64(y), float64(w), float64(h)),
	}
	newCam.Cam.ZoomFactor = 3
	newCam.Cam.LerpEnabled = true
	newCam.Cam.LerpSpeed = 0.7
	newCam.Cam.ShakeEnabled = true
	return newCam
}
