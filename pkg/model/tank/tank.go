package tank

import (
	"math"

	"sfs.ink/liang/game/pkg/attr"
)

type Tank struct {
	name string
	// position value offset to the origin of world
	position                       attr.Offset
	bodyAngle, targetBodyAngle     float64
	turretAngle, targetTurretAngle float64
}

func NewTank(name string, position attr.Offset, bodyAngle, targetBodyAngle, turretAngle, targetTurretAngle float64) Tank {
	return Tank{name, position, bodyAngle, targetBodyAngle, turretAngle, targetTurretAngle}
}

func (t Tank) Name() string {
	return t.name
}

func (t Tank) Attr() (dx, dy, bodyAngle, turretAngle float64) {
	return t.position.Dx, t.position.Dy, t.bodyAngle, t.turretAngle
}

func (t *Tank) UpdateTo(tk Tank) {
	t.position = tk.position
	t.bodyAngle, t.targetBodyAngle = tk.bodyAngle, tk.targetBodyAngle
	t.turretAngle, t.targetTurretAngle = tk.turretAngle, tk.targetTurretAngle
}

func (t *Tank) Update(dt float64) {
	rotationRate := math.Pi * dt
	if t.targetBodyAngle != 0 {
		if t.bodyAngle < t.targetBodyAngle {
			if math.Abs(t.targetBodyAngle-t.bodyAngle) > math.Pi {
				t.bodyAngle = t.bodyAngle - rotationRate
				if t.bodyAngle < -math.Pi {
					t.bodyAngle += math.Pi * 2
				}
			} else {
				t.bodyAngle = t.bodyAngle + rotationRate
				if t.bodyAngle > t.targetBodyAngle {
					t.bodyAngle = t.targetBodyAngle
				}
			}
		}
		if t.bodyAngle > t.targetBodyAngle {
			if math.Abs(t.targetBodyAngle-t.bodyAngle) > math.Pi {
				t.bodyAngle = t.bodyAngle + rotationRate
				if t.bodyAngle > math.Pi {
					t.bodyAngle -= math.Pi * 2
				}
			} else {
				t.bodyAngle = t.bodyAngle - rotationRate
				if t.bodyAngle < t.targetBodyAngle {
					t.bodyAngle = t.targetBodyAngle
				}
			}
		}
		if t.bodyAngle == t.targetBodyAngle {
			of := attr.FromDirection(t.bodyAngle, 100*dt)
			t.position = t.position.Add(of)
		} else {
			of := attr.FromDirection(t.bodyAngle, 50*dt)
			t.position = t.position.Add(of)
		}
	}
	if t.targetTurretAngle != 0 {
		localTargetTurretAngle := t.targetTurretAngle - t.bodyAngle
		if t.turretAngle < localTargetTurretAngle {
			if math.Abs(localTargetTurretAngle-t.turretAngle) > math.Pi {
				t.turretAngle = t.turretAngle - rotationRate
				if t.turretAngle < -math.Pi {
					t.turretAngle += math.Pi * 2
				}
			} else {
				t.turretAngle = t.targetTurretAngle + rotationRate
				if t.turretAngle > localTargetTurretAngle {
					t.turretAngle = localTargetTurretAngle
				}
			}
		}
		if t.turretAngle > localTargetTurretAngle {
			if math.Abs(localTargetTurretAngle-t.turretAngle) > math.Pi {
				t.turretAngle = t.turretAngle + rotationRate
				if t.turretAngle > math.Pi {
					t.turretAngle -= math.Pi * 2
				}
			} else {
				t.turretAngle = t.turretAngle - rotationRate
				if t.turretAngle < localTargetTurretAngle {
					t.turretAngle = localTargetTurretAngle
				}
			}
		}
	}
}

func (t *Tank) GetBulletOffset() attr.Offset {
	of := attr.FromDirection(t.GetBulletAngle(), 36)
	return t.position.Add(of)
}

func (t *Tank) GetBulletAngle() float64 {
	bulletAngle := t.bodyAngle + t.turretAngle
	for bulletAngle > math.Pi {
		bulletAngle -= math.Pi * 2
	}
	for bulletAngle < -math.Pi {
		bulletAngle += math.Pi * 2
	}
	return bulletAngle
}
