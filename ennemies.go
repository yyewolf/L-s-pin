package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnnemyState struct {
	PosX   int
	PosY   int
	SpeedX int
	SpeedY int
	Angle  float64
}

type Ennemy struct {
	CurrentState  EnnemyState
	PreviousState EnnemyState
}

func addEnnemy() {
	State := EnnemyState{}
	speed := rand.Intn(5) + 1
	rng := rand.Float64()
	if rng < 0.25 {
		State = EnnemyState{
			PosX:   screenWidth / 2,
			PosY:   0,
			SpeedX: 0,
			SpeedY: speed,
			Angle:  3 * math.Pi / 2,
		}
	} else if rng < 0.5 {
		State = EnnemyState{
			PosX:   screenWidth / 2,
			PosY:   screenHeight,
			SpeedX: 0,
			SpeedY: -speed,
			Angle:  math.Pi / 2,
		}
	} else if rng < 0.75 {
		State = EnnemyState{
			PosX:   0,
			PosY:   screenHeight / 2,
			SpeedX: speed,
			SpeedY: 0,
			Angle:  7,
		}
	} else {
		State = EnnemyState{
			PosX:   screenWidth,
			PosY:   screenHeight / 2,
			SpeedX: -speed,
			SpeedY: 0,
			Angle:  0,
		}
	}

	newEnnemy := &Ennemy{
		CurrentState: State,
	}
	newEnnemy.PreviousState = newEnnemy.CurrentState
	spawnedEnnemies = append(spawnedEnnemies, newEnnemy)
}

func (g *Game) pickEnnemy() {
	if rand.Float64() < float64(g.count)/1000 && time.Since(lastEnnemy) > 600*time.Millisecond {
		lastEnnemy = time.Now()
		addEnnemy()
	}
}

func (g *Game) drawAllEnnemies(screen *ebiten.Image) {
	w, h := ennemies.Size()
	op := &ebiten.DrawImageOptions{}

	var ennemiesAlive []*Ennemy
	for _, e := range spawnedEnnemies {
		// Calculates new pos

		e.CurrentState.PosX = e.PreviousState.PosX + e.PreviousState.SpeedX
		e.CurrentState.PosY = e.PreviousState.PosY + e.PreviousState.SpeedY
		e.PreviousState = e.CurrentState

		// Check hit box
		X, Y := e.CurrentState.PosX, e.CurrentState.PosY
		if !(X > (screenWidth/2)-70 && X < (screenWidth/2)+70 && Y > (screenHeight/2)-70 && Y < (screenHeight/2)+70) {
			ennemiesAlive = append(ennemiesAlive, e)
		} else {
			// Check for success or loss
			if checkAngles(playerPos, e.CurrentState.Angle) {
				playerScore += 1
				p := g.audioContext.NewPlayerFromBytes(hitSound)
				p.SetVolume(0.1)
				p.Play()
			} else {
				playerScore -= 2
			}
			if playerScore < 0 {
				playerScore = 0
			}
		}

		// Handle ennemies images
		op.GeoM.Reset()
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)

		if e.CurrentState.Angle == 7 {
			op.GeoM.Scale(-1, 1)
		} else {
			op.GeoM.Rotate(e.CurrentState.Angle)
		}
		op.GeoM.Translate(float64(e.CurrentState.PosX), float64(e.CurrentState.PosY))
		op.ColorM.RotateHue(e.CurrentState.Angle)
		screen.DrawImage(ennemies, op)
	}

	spawnedEnnemies = ennemiesAlive
}

func checkAngles(pAng, eAng float64) bool {
	switch eAng {
	case 3 * math.Pi / 2:
		if pAng == 0 {
			return true
		}
		return false
	case math.Pi / 2:
		if pAng == math.Pi {
			return true
		}
		return false
	case 0:
		if pAng == math.Pi/2 {
			return true
		}
		return false
	case 7:
		if pAng == 3*math.Pi/2 {
			return true
		}
		return false
	}
	return false
}
