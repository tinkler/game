package char

import "sfs.ink/liang/game/pkg/attr"

type Char struct {
	p     attr.Position
	speed attr.Speed
}

func (c *Char) MoveUp() {
	if c.speed.IsZero() {
		return
	}
	c.p.X -= float64(c.speed)
}

func (c *Char) MoveRight() {
	if c.speed.IsZero() {
		return
	}
	c.p.Y += float64(c.speed)
}

func (c *Char) MoveDown() {
	if c.speed.IsZero() {
		return
	}
	c.p.X += float64(c.speed)
}

func (c *Char) MoveLeft() {
	if c.speed.IsZero() {
		return
	}
	c.p.Y -= float64(c.speed)
}
