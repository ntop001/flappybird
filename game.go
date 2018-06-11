package main

import (
	"korok.io/korok/game"
	"korok.io/korok/gui"
	"korok.io/korok/gfx"
	"korok.io/korok/asset"
	"korok.io/korok/engi"
	"korok.io/korok"
	"korok.io/korok/math/f32"
	"korok.io/korok/hid/input"
)


type StateEnum int

const (
	Ready StateEnum = iota
	Running
	Over
)

type BirdStateEnum int

const (
	Flying BirdStateEnum = iota
	Dead
)


const (
	Gravity = 600
	TapImpulse = 280
	ScrollVelocity = 100

	RotTrigger = 200
	MaxAngle = 3.14/6
	MinAngle = -3.14/2
	AngleVelocity = 3.14 * 4
)

type GameScene struct {
	state StateEnum

	ready struct{
		gfx.Tex2D
		gui.Rect
	}
		tap struct{
		gfx.Tex2D
		gui.Rect
	}

	bg engi.Entity

	bird struct{
		state BirdStateEnum
		engi.Entity
		f32.Vec2
		vy float32
		w, h float32
		rotate float32
	}

	ground struct{
		engi.Entity
		f32.Vec2
		vx float32
	}

	PipeSystem
}

func (sn *GameScene) borrow(bird, bg, ground engi.Entity) {
	sn.bird.Entity, sn.bg, sn.ground.Entity = bird, bg, ground
}

func (sn *GameScene) OnEnter(g *game.Game) {
	at, _ := asset.Texture.Atlas("images/bird.png")

	// ready and tap image
	sn.ready.Tex2D, _ = at.GetByName("getready.png")
	sn.ready.Rect = gui.Rect{
		X: (320-233)/2,
		Y: 70,
		W: 233,
		H: 70,
	}
	sn.tap.Tex2D, _ = at.GetByName("tap.png")
	sn.tap.Rect = gui.Rect{
		X: (320-143)/2,
		Y: 200,
		W: 143, // 286
		H: 123, // 246
	}

	korok.Transform.Comp(sn.bird.Entity).SetPosition(f32.Vec2{80, 240})
	sn.bird.Vec2 = f32.Vec2{80, 240}

	sn.ground.Vec2 = f32.Vec2{0, 100}
	sn.ground.vx = ScrollVelocity

	// setup pipes (129, 801)
	top, _ := at.GetByName("top_pipe.png")
	bottom, _ := at.GetByName("bottom_pipe.png")

	ps := &sn.PipeSystem
	ps.initialize(top, bottom, 6)
	ps.setDelay(0) // 3 seconds
	ps.setRate(2.5)  // generate pipe every 2 seconds
	ps.setGap(100)
	ps.setLimit(300, 150)
	ps.StartScroll()

}

func (sn *GameScene) Update(dt float32) {
	if st := sn.state; st == Ready {
		sn.showReady(dt); return
	} else if st == Over {
		sn.showOver(dt)
		return
	}

	if input.PointerButton(0).JustPressed() {
		sn.bird.vy = TapImpulse
	}
	sn.bird.vy -= Gravity * dt
	sn.bird.Vec2[1] += sn.bird.vy * dt

	// rotate
	if sn.bird.vy > -RotTrigger && sn.bird.rotate < MaxAngle {
		sn.bird.rotate +=  AngleVelocity * dt
	} else if sn.bird.vy < -RotTrigger && sn.bird.rotate > MinAngle {
		sn.bird.rotate += -AngleVelocity * dt
	}

	// update bird position
	b := korok.Transform.Comp(sn.bird.Entity)
	b.SetPosition(sn.bird.Vec2)
	b.SetRotation(sn.bird.rotate)

	// scroll background
	// windows.width = 320, ground.width = 420, so if we move by 100(420-320),
	// then reset position, it looks like the ground scroll seamless.
	x := sn.ground.Vec2[0]
	if x < -100 {
		x = x + 90 // magic number (bridge start and end of the image)
	}
	x -= sn.ground.vx * dt
	sn.ground.Vec2[0] = x

	// update ground shift
	g := korok.Transform.Comp(sn.ground.Entity)
	g.SetPosition(sn.ground.Vec2)

	// update pipes
	sn.PipeSystem.Update(dt)
	// detect collision with pipes
	ps := &sn.PipeSystem
	if c, _ := ps.CheckCollision(sn.bird.Vec2, f32.Vec2{sn.bird.w, sn.bird.h}); c {
		if sn.bird.state != Dead {
			ps.StopScroll()
			sn.bird.state = Dead

			// stop bird animation
			korok.Flipbook.Comp(sn.bird.Entity).Stop()
		}
	}


	// detect collision with ground and sky
	if y := sn.bird.Vec2[1]; y > 480 {
		sn.bird.Vec2[1] = 480
	} else if y < 100 {
		y = 100; sn.state = Over

		if sn.bird.state != Dead {
			sn.bird.state = Dead
			korok.Flipbook.Comp(sn.bird.Entity).Stop()
		}
	}
}

func (sn *GameScene) showReady(dt float32) {
	// show ready
	gui.Image(1, sn.ready.Rect, sn.ready.Tex2D, nil)

	// show tap hint
	gui.Image(2, sn.tap.Rect, sn.tap.Tex2D, nil)

	// check any click
	if input.PointerButton(0).JustPressed() {
		sn.state = Running
	}
}

func (sn *GameScene) showOver(dt float32) {

}


func (sn *GameScene) OnExit() {
}
