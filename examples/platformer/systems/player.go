package systems

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/examples/platformer/components"
	dresolv "github.com/yohamta/donburi/examples/platformer/resolv"
	"github.com/yohamta/donburi/examples/platformer/tags"
)

func UpdatePlayer(ecs *ecs.ECS) {
	// Now we update the Player's movement. This is the real bread-and-butter of this example, naturally.
	playerEntry, _ := components.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	playerObject := dresolv.GetObject(playerEntry)

	friction := 0.5
	accel := 0.5 + friction
	maxSpeed := 4.0
	jumpSpd := 10.0
	gravity := 0.75

	player.SpeedY += gravity
	if player.WallSliding != nil && player.SpeedY > 1 {
		player.SpeedY = 1
	}

	// Horizontal movement is only possible when not wallsliding.
	if player.WallSliding == nil {
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.GamepadAxis(0, 0) > 0.1 {
			player.SpeedX += accel
			player.FacingRight = true
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.GamepadAxis(0, 0) < -0.1 {
			player.SpeedX -= accel
			player.FacingRight = false
		}
	}

	// Apply friction and horizontal speed limiting.
	if player.SpeedX > friction {
		player.SpeedX -= friction
	} else if player.SpeedX < -friction {
		player.SpeedX += friction
	} else {
		player.SpeedX = 0
	}

	if player.SpeedX > maxSpeed {
		player.SpeedX = maxSpeed
	} else if player.SpeedX < -maxSpeed {
		player.SpeedX = -maxSpeed
	}

	// Check for jumping.
	if inpututil.IsKeyJustPressed(ebiten.KeyX) || ebiten.IsGamepadButtonPressed(0, 0) || ebiten.IsGamepadButtonPressed(1, 0) {

		if (ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.GamepadAxis(0, 1) > 0.1 || ebiten.GamepadAxis(1, 1) > 0.1) && player.OnGround != nil && player.OnGround.HasTags("platform") {

			player.IgnorePlatform = player.OnGround

		} else {

			if player.OnGround != nil {
				player.SpeedY = -jumpSpd
			} else if player.WallSliding != nil {
				// WALLJUMPING
				player.SpeedY = -jumpSpd

				if player.WallSliding.X > playerObject.X {
					player.SpeedX = -4
				} else {
					player.SpeedX = 4
				}

				player.WallSliding = nil

			}

		}

	}

	// We handle horizontal movement separately from vertical movement. This is, conceptually, decomposing movement into two phases / axes.
	// By decomposing movement in this manner, we can handle each case properly (i.e. stop movement horizontally separately from vertical movement, as
	// necesseary). More can be seen on this topic over on this blog post on higherorderfun.com:
	// http://higherorderfun.com/blog/2012/05/20/the-guide-to-implementing-2d-platformers/

	// dx is the horizontal delta movement variable (which is the Player's horizontal speed). If we come into contact with something, then it will
	// be that movement instead.
	dx := player.SpeedX

	// Moving horizontally is done fairly simply; we just check to see if something solid is in front of us. If so, we move into contact with it
	// and stop horizontal movement speed. If not, then we can just move forward.

	if check := playerObject.Check(player.SpeedX, 0, "solid"); check != nil {

		dx = check.ContactWithCell(check.Cells[0]).X()
		player.SpeedX = 0

		// If you're in the air, then colliding with a wall object makes you start wall sliding.
		if player.OnGround == nil {
			player.WallSliding = check.Objects[0]
		}

	}

	// Then we just apply the horizontal movement to the Player's Object. Easy-peasy.
	playerObject.X += dx

	// Now for the vertical movement; it's the most complicated because we can land on different types of objects and need
	// to treat them all differently, but overall, it's not bad.

	// First, we set OnGround to be nil, in case we don't end up standing on anything.
	player.OnGround = nil

	// dy is the delta movement downward, and is the vertical movement by default; similarly to dx, if we come into contact with
	// something, this will be changed to move to contact instead.

	dy := player.SpeedY

	// We want to be sure to lock vertical movement to a maximum of the size of the Cells within the Space
	// so we don't miss any collisions by tunneling through.

	dy = math.Max(math.Min(dy, 16), -16)

	// We're going to check for collision using dy (which is vertical movement speed), but add one when moving downwards to look a bit deeper down
	// into the ground for solid objects to land on, specifically.
	checkDistance := dy
	if dy >= 0 {
		checkDistance++
	}

	// We check for any solid / stand-able objects. In actuality, there aren't any other Objects
	// with other tags in this Space, so we don't -have- to specify any tags, but it's good to be specific for clarity in this example.
	if check := playerObject.Check(0, checkDistance, "solid", "platform", "ramp"); check != nil {

		// So! Firstly, we want to see if we jumped up into something that we can slide around horizontally to avoid bumping the Player's head.

		// Sliding around a misspaced jump is a small thing that makes jumping a bit more forgiving, and is something different polished platformers
		// (like the 2D Mario games) do to make it a smidge more comfortable to play. For a visual example of this, see this excellent devlog post
		// from the extremely impressive indie game, Leilani's Island: https://forums.tigsource.com/index.php?topic=46289.msg1387138#msg1387138

		// To accomplish this sliding, we simply call Collision.SlideAgainstCell() to see if we can slide.
		// We pass the first cell, and tags that we want to avoid when sliding (i.e. we don't want to slide into cells that contain other solid objects).

		slide := check.SlideAgainstCell(check.Cells[0], "solid")

		// We further ensure that we only slide if:
		// 1) We're jumping up into something (dy < 0),
		// 2) If the cell we're bumping up against contains a solid object,
		// 3) If there was, indeed, a valid slide left or right, and
		// 4) If the proposed slide is less than 8 pixels in horizontal distance. (This is a relatively arbitrary number that just so happens to be half the
		// width of a cell. This is to ensure the player doesn't slide too far horizontally.)

		if dy < 0 && check.Cells[0].ContainsTags("solid") && slide != nil && math.Abs(slide.X()) <= 8 {

			// If we are able to slide here, we do so. No contact was made, and vertical speed (dy) is maintained upwards.
			playerObject.X += slide.X()

		} else {

			// If sliding -fails-, that means the Player is jumping directly onto or into something, and we need to do more to see if we need to come into
			// contact with it. Let's press on!

			// First, we check for ramps. For ramps, we can't simply check for collision with Check(), as that's not precise enough. We need to get a bit
			// more information, and so will do so by checking its Shape (a triangular ConvexPolygon, as defined in WorldPlatformer.Init()) against the
			// Player's Shape (which is also a rectangular ConvexPolygon).

			// We get the ramp by simply filtering out Objects with the "ramp" tag out of the objects returned in our broad Check(), and grabbing the first one
			// if there's any at all.
			if ramps := check.ObjectsByTags("ramp"); len(ramps) > 0 {

				ramp := ramps[0]

				// For simplicity, this code assumes we can only stand on one ramp at a time as there is only one ramp in this example.
				// In actuality, if there was a possibility to have a potential collision with multiple ramps (i.e. a ramp that sits on another ramp, and the player running down
				// one onto the other), the collision testing code should probably go with the ramp with the highest confirmed intersection point out of the two.

				// Next, we see if there's been an intersection between the two Shapes using Shape.Intersection. We pass the ramp's shape, and also the movement
				// we're trying to make horizontally, as this makes Intersection return the next y-position while moving, not the one directly
				// underneath the Player. This would keep the player from getting "stuck" when walking up a ramp into the top of a solid block, if there weren't
				// a landing at the top and bottom of the ramp.

				// We use 8 here for the Y-delta so that we can easily see if you're running down the ramp (in which case you're probably in the air as you
				// move faster than you can fall in this example). This way we can maintain contact so you can always jump while running down a ramp. We only
				// continue with coming into contact with the ramp as long as you're not moving upwards (i.e. jumping).

				if contactSet := playerObject.Shape.Intersection(dx, 8, ramp.Shape); dy >= 0 && contactSet != nil {

					// If Intersection() is successful, a ContactSet is returned. A ContactSet contains information regarding where
					// two Shapes intersect, like the individual points of contact, the center of the contacts, and the MTV, or
					// Minimum Translation Vector, to move out of contact.

					// Here, we use ContactSet.TopmostPoint() to get the top-most contact point as an indicator of where
					// we want the player's feet to be. Then we just set that position, and we're done.

					dy = contactSet.TopmostPoint()[1] - playerObject.Bottom() + 0.1
					player.OnGround = ramp
					player.SpeedY = 0

				}

			}

			// Platforms are next; here, we just see if the platform is not being ignored by attempting to drop down,
			// if the player is falling on the platform (as otherwise he would be jumping through platforms), and if the platform is low enough
			// to land on. If so, we stand on it.

			// Because there's a moving floating platform, we use Collision.ContactWithObject() to ensure the player comes into contact
			// with the top of the platform object. An alternative would be to use Collision.ContactWithCell(), but that would be only if the
			// platform didn't move and were aligned with the Spatial cellular grid.

			if platforms := check.ObjectsByTags("platform"); len(platforms) > 0 {

				platform := platforms[0]

				if platform != player.IgnorePlatform && player.SpeedY >= 0 && playerObject.Bottom() < platform.Y+4 {
					dy = check.ContactWithObject(platform).Y()
					player.OnGround = platform
					player.SpeedY = 0
				}

			}

			// Finally, we check for simple solid ground. If we haven't had any success in landing previously, or the solid ground
			// is higher than the existing ground (like if the platform passes underneath the ground, or we're walking off of solid ground
			// onto a ramp), we stand on it instead. We don't check for solid collision first because we want any ramps to override solid
			// ground (so that you can walk onto the ramp, rather than sticking to solid ground).

			// We use ContactWithObject() here because otherwise, we might come into contact with the moving platform's cells (which, naturally,
			// would be selected by a Collision.ContactWithCell() call because the cell is closest to the Player).

			if solids := check.ObjectsByTags("solid"); len(solids) > 0 && (player.OnGround == nil || player.OnGround.Y >= solids[0].Y) {
				dy = check.ContactWithObject(solids[0]).Y()
				player.SpeedY = 0

				// We're only on the ground if we land on it (if the object's Y is greater than the player's).
				if solids[0].Y > playerObject.Y {
					player.OnGround = solids[0]
				}

			}

			if player.OnGround != nil {
				player.WallSliding = nil    // Player's on the ground, so no wallsliding anymore.
				player.IgnorePlatform = nil // Player's on the ground, so reset which platform is being ignored.
			}

		}

	}

	// Move the object on dy.
	playerObject.Y += dy

	wallNext := 1.0
	if !player.FacingRight {
		wallNext = -1
	}

	// If the wall next to the Player runs out, stop wall sliding.
	if c := playerObject.Check(wallNext, 0, "solid"); player.WallSliding != nil && c == nil {
		player.WallSliding = nil
	}
}

func DrawPlayer(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Player.Each(ecs.World, func(e *donburi.Entry) {
		player := components.Player.Get(e)
		o := dresolv.GetObject(e)
		playerColor := color.RGBA{0, 255, 60, 255}
		if player.OnGround == nil {
			// We draw the player as a different color when jumping so we can visually see when he's in the air.
			playerColor = color.RGBA{200, 0, 200, 255}
		}
		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, playerColor)
	})
}
