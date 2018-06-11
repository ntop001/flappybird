package main

import (
	"korok.io/korok/engi"
	"korok.io/korok"
	"korok.io/korok/gfx"
	"korok.io/korok/math"
	"korok.io/korok/math/f32"
)

// Pipe manager system.

type Pipe struct {
	top struct{
		engi.Entity
		f32.Vec2
	}
	bottom struct{
		engi.Entity
		f32.Vec2
	}
	high float32
	active bool
	x, vx float32
}

func (p *Pipe) initialize(texTop, texBottom gfx.Tex2D) {
	top := korok.Entity.New()
	spr := korok.Sprite.NewComp(top)
	spr.SetSprite(texTop)
	spr.SetSize(65, 400)
	spr.SetGravity(.5, 0)

	bottom := korok.Entity.New()
	spr = korok.Sprite.NewComp(bottom)
	spr.SetSize(65, 400)
	spr.SetSprite(texBottom)
	spr.SetGravity(.5, 1)

	// out of screen
	korok.Transform.NewComp(top).SetPosition(f32.Vec2{-100, 210})
	korok.Transform.NewComp(bottom).SetPosition(f32.Vec2{-100, 160})

	p.top.Entity = top
	p.bottom.Entity = bottom
	p.vx = ScrollVelocity
}

func (p *Pipe) reset(x, high, gap float32) {
	p.active = true
	p.x = x
	p.top.Vec2 = f32.Vec2{x, high + gap}
	p.bottom.Vec2 = f32.Vec2{x, high}
}

func (p *Pipe) update(dt float32) {
	p.x -= p.vx * dt
	if p.x < -50 {
		p.active = false
	}

	p.top.Vec2[0] = p.x
	p.bottom.Vec2[0] = p.x

	korok.Transform.Comp(p.top.Entity).SetPosition(p.top.Vec2)
	korok.Transform.Comp(p.bottom.Entity).SetPosition(p.bottom.Vec2)
}

type PipeSystem struct {
	gap, top, bottom float32 // gap, top, bottom limit
	respawn float32          // respawn location
	scroll bool

	delay struct{
		clock float32
		limit float32
	}
	generate struct{
		clock float32
		limit float32
	}

	pipes []*Pipe
	frees []*Pipe

	_pool []Pipe
}

func (ps *PipeSystem) initialize(texTop, texBottom gfx.Tex2D, size int) {
	ps._pool = make([]Pipe, size)
	ps.frees = make([]*Pipe, size) // add to freelist
	for i := range ps._pool {
		ps.frees[i] = &ps._pool[i]
		ps.frees[i].initialize(texTop, texBottom)
	}
	ps.respawn = 320 + 20
}

func (ps *PipeSystem) setDelay(d float32) {
	ps.delay.limit = d
}

func (ps *PipeSystem) setRate(r float32) {
	ps.generate.limit = r
}

func (ps *PipeSystem) setGap(gap float32) {
	ps.gap = gap
}

func (ps *PipeSystem) setLimit(top, b float32) {
	ps.top, ps.bottom = top, b
}

func (ps *PipeSystem) Update(dt float32) {
	if !ps.scroll {
		return
	}

	// delay some time
	if d := &ps.delay; d.clock < d.limit {
		d.clock += dt; return
	}

	// generate new pipe
	if g := &ps.generate; g.clock < g.limit {
		g.clock += dt
	} else {
		g.clock = 0
		ps.newPipe()
	}

	// update pipe
	for _, p := range ps.pipes {
		p.update(dt)
	}

	// recycle
	ps.recycle()
}

func (ps *PipeSystem) StopScroll() {
	ps.scroll = false
}

func (ps *PipeSystem) StartScroll() {
	ps.scroll = true
}

func (ps *PipeSystem) Reset() {
	for _, p := range ps.pipes {
		p.x = -100
		// out of screen
		korok.Transform.NewComp(p.top.Entity).SetPosition(f32.Vec2{-100, 210})
		korok.Transform.NewComp(p.bottom.Entity).SetPosition(f32.Vec2{-100, 160})
	}
}

func (ps *PipeSystem) newPipe() {
	if sz := len(ps.frees); sz > 0 {
		p := ps.frees[sz-1]; ps.frees = ps.frees[:sz-1]
		ps.pipes = append(ps.pipes, p)
		p.reset(ps.respawn, math.Random(ps.bottom, ps.top), ps.gap)
	}
}

// inactive pipes come first
func (ps *PipeSystem) recycle() {
	pipes, inactive := ps.pipes, -1
	for i, p := range pipes {
		if p.active {
			break
		}
		inactive = i
	}
	if inactive >= 0 {
		ps.pipes = pipes[inactive+1:]
		ps.frees = append(ps.frees, pipes[:inactive+1]...)
	}
}

//  左下坐标
type AABB struct {
	x, y float32
	width, height float32
}

func OverlapAB(a, b *AABB) bool {
	if a.x < b.x+b.width && a.x+a.width>b.x && a.y < b.y+b.height && a.y+a.height > b.y {
		return true
	}
	return false
}

// check collision
func (ps *PipeSystem) CheckCollision(p f32.Vec2, sz f32.Vec2) (bool, float32) {
	tolerance := float32(8)
	sz[0], sz[1] = sz[0]-tolerance, sz[1]-tolerance
	bird := &AABB{p[0]-sz[0]/2, p[1]-sz[1]/2, sz[0], sz[1]}
	for _, p := range ps.pipes {
		top := &AABB{
			p.top.Vec2[0] - 32,
			p.top.Vec2[1],
			65,
			400,
		}
		if OverlapAB(bird, top) {
			return true, bird.x - top.x
		}

		bottom := &AABB{
			p.bottom.Vec2[0] - 32,
			p.bottom.Vec2[1] - 400,
			65,
			400,
		}
		if OverlapAB(bird, bottom) {
			return true, bird.x - top.x
		}
	}
	return false, 0
}


