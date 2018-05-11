package main

import (
	"context"
	"testing"
	"time"

	"github.com/rhcarvalho/tiwe/router"
)

func Test(t *testing.T) {
	r := router.NewTestRouter(40*time.Millisecond, 20*time.Millisecond)

	alice := &Player{
		Name:   "alice",
		Router: r,
	}
	bob := &Player{
		Name:   "bob",
		Router: r,
	}
	carol := &Player{
		Name:   "carol",
		Router: r,
	}

	ctx := context.Background()
	players := []string{alice.Name, bob.Name, carol.Name}

	alice.Start(ctx, players)
	bob.Start(ctx, players)
	carol.Start(ctx, players)

	alice.Init()

	defer alice.Stop(ctx)
	defer bob.Stop(ctx)
	defer carol.Stop(ctx)

	<-alice.Done()
	<-bob.Done()
	<-carol.Done()
}
