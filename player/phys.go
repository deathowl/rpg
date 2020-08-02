package player
import 	"github.com/faiface/pixel"

type playerPhys struct {
	gravity   float64
	runSpeed  float64
	jumpSpeed float64

	rect   pixel.Rect
	vel    pixel.Vec
	ground bool
}

func (pp *playerPhys) update(dt float64, ctrl pixel.Vec) {
	// apply controls
	switch {
	case ctrl.X < 0:
		pp.vel.X = -pp.runSpeed
	case ctrl.X > 0:
		pp.vel.X = +pp.runSpeed
	default:
		pp.vel.X = 0
	}
	switch {
	case ctrl.Y < 0:
		pp.vel.Y = -pp.runSpeed
	case ctrl.Y > 0:
		pp.vel.Y = +pp.runSpeed
	default:
		pp.vel.X = 0
	}

	// apply gravity and velocity
	pp.rect = pp.rect.Moved(pp.vel.Scaled(dt))

	// check collisions against each platform
	// gp.ground = false
	// if gp.vel.Y <= 0 {
	// 	for _, p := range platforms {
	// 		if gp.rect.Max.X <= p.rect.Min.X || gp.rect.Min.X >= p.rect.Max.X {
	// 			continue
	// 		}
	// 		if gp.rect.Min.Y > p.rect.Max.Y || gp.rect.Min.Y < p.rect.Max.Y+gp.vel.Y*dt {
	// 			continue
	// 		}
	// 		gp.vel.Y = 0
	// 		gp.rect = gp.rect.Moved(pixel.V(0, p.rect.Max.Y-gp.rect.Min.Y))
	// 		gp.ground = true
	// 	}
	// }

	// // jump if on the ground and the player wants to jump
	// if gp.ground && ctrl.Y > 0 {
	// 	gp.vel.Y = gp.jumpSpeed
	// }
}
