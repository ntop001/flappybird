package main

import (
	"korok.io/korok/game"
	"korok.io/korok"
	"korok.io/korok/asset"
	"korok.io/korok/math/f32"
	"korok.io/korok/gfx"
)

type StartScene struct {
}

func (sn *StartScene) Load() {
	asset.Texture.LoadAtlas("images/bird.png", "images/bird.json")
}

func (sn *StartScene) OnEnter(g *game.Game) {
	at, _ := asset.Texture.Atlas("images/bird.png")
	bg, _ := at.GetByName("background.png")
	ground, _ := at.GetByName("ground.png")

	// setup bg
	{
		entity := korok.Entity.New()
		spr := korok.Sprite.NewCompX(entity, bg)
		spr.SetSize(320, 480)
		xf := korok.Transform.NewComp(entity)
		xf.SetPosition( f32.Vec2{160, 240})
	}

	// setup ground {840 281}
	{
		entity := korok.Entity.New()
		spr := korok.Sprite.NewCompX(entity, ground)
		spr.SetSize(420, 140)
		spr.SetGravity(0, 1)
		spr.SetZOrder(1)
		xf := korok.Transform.NewComp(entity)
		xf.SetPosition(f32.Vec2{0, 100})
	}

	// flying animation
	bird1, _ := at.GetByName("bird1.png")
	bird2, _ := at.GetByName("bird2.png")
	bird3, _ := at.GetByName("bird3.png")

	frames := []gfx.Tex2D{bird1, bird2, bird3}
	g.AnimationSystem.SpriteEngine.NewAnimation("flying", frames, true)

	// setup bird
	bird := korok.Entity.New()
	spr := korok.Sprite.NewCompX(bird, bird1)
	spr.SetSize(48, 32)
	spr.SetZOrder(2)
	xf := korok.Transform.NewComp(bird)
	xf.SetPosition(f32.Vec2{160, 240})

	anim := korok.Flipbook.NewComp(bird)
	anim.SetRate(.1)
	anim.Play("flying")
}
func (sn *StartScene) Update(dt float32) {
	// draw something
}
func (sn *StartScene) OnExit() {
}

func main() {
	options := korok.Options{
		Title:"Flappy Bird",
		Width:320,
		Height:480,
	}
	korok.Run(&options, &StartScene{})
}
