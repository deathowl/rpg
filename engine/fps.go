package engine

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

// FPSWatch measures the real-time frame rates and displays it on a target canvas.
type FPSWatch struct {
	txt   *text.Text     // shared variable
	atlas *text.Atlas    // borrowed atlas for txt
	imd   *imdraw.IMDraw // shared variable
	mutex sync.Mutex     // synchronize
	//
	fps    int              // The FPS evaluated every second.
	frames int              // Frames count before the FPS update.
	seccer <-chan time.Time // Ticks time every second.
	//
	desc     string
	colorBg  color.Color
	colorTxt color.Color
}

// NewFPSWatch is a constructor.
func NewFPSWatch(
	additionalCaption string,
	_colorBg, _colorTxt color.Color,
) (watch *FPSWatch) {
	return &FPSWatch{
		fps:      0,
		frames:   0,
		seccer:   nil,
		desc:     additionalCaption,
		colorBg:  _colorBg,
		colorTxt: _colorTxt,
	}
}

// NewFPSWatchSimple is a constructor.
func NewFPSWatchSimple() *FPSWatch {
	return NewFPSWatch("", colornames.Black, colornames.White)
}

// Start ticking every second.
func (watch *FPSWatch) Start() {
	watch.seccer = time.Tick(time.Second)
}

// Poll () should be called only once and in every single frame. (Obligatory)
// This is an extended behavior of Update() like funcs.
func (watch *FPSWatch) Poll() {
	watch.frames++
	select {
	case <-watch.seccer:
		watch.fps = watch.frames
		watch.frames = 0
		go watch._Update()
	default:
	}
}

// GetFPS returns the most recent FPS recorded.
// A non-ptr FPSWatch as a read only argument passes lock by value within itself but that seems totally fine.
func (watch FPSWatch) GetFPS() int {
	return watch.fps
}

// Draw FPSWatch.
func (watch *FPSWatch) Draw(t pixel.Target, pos pixel.Vec) {
	// lock before accessing txt & imdraw
	watch.mutex.Lock()
	defer watch.mutex.Unlock()
	if watch.imd == nil { // isInvisible set to true.
		return // An empty image is drawn.
	}
	str := fmt.Sprint("FPS: ", watch.fps, " ", watch.desc)
	txt := text.New(pos, watch.atlas)
	txt.Color = watch.colorTxt
	txt.Dot.X -= 1.0
	txt.Dot.Y += 5.0
	txt.WriteString(str)

	//watch.imd.Draw(t)
	txt.Draw(t, pixel.IM)
}

// unexported
func (watch *FPSWatch) _Update() {
	// lock before txt & imdraw update
	watch.mutex.Lock()
	defer watch.mutex.Unlock()
	if watch.atlas == nil {
		face := basicfont.Face7x13
		watch.atlas = text.NewAtlas(face, text.ASCII, nil)
	}

	// imdraw (a state machine)
	if watch.imd == nil { // lazy creation
		watch.imd = imdraw.New(nil)
	}
	imd := watch.imd
	imd.Clear()

	imd.Color = watch.colorBg
}
