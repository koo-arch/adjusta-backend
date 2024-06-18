package scheduler

import (
	"context"
	"log"
	
	"github.com/robfig/cron/v3"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/auth"
)

type Scheduler struct {
	client *ent.Client
	cron  *cron.Cron
}

func NewScheduler(client *ent.Client) *Scheduler {
	return &Scheduler{
		client: client,
		cron: cron.New(),
	}
}

func (s *Scheduler) SetupJobs(ctx context.Context) {
	_, err := s.cron.AddFunc("@daily", func() {
		// キーの生成
		keyManager := auth.NewKeyManager(s.client)
		err := keyManager.GenerateJWTKey(ctx, "access")
		if err != nil {
			log.Printf("failed to generate JWT key: %v", err)
		}
		err = keyManager.GenerateJWTKey(ctx, "refresh")
		if err != nil {
			log.Printf("failed to generate JWT key: %v", err)
		}

		// キーの削除
		err = keyManager.DeleteJWTKeys(ctx, "access")
		if err != nil {
			log.Printf("failed to delete access JWT key: %v", err)
		}
		err = keyManager.DeleteJWTKeys(ctx, "refresh")
		if err != nil {
			log.Printf("failed to delete refresh JWT key: %v", err)
		}

	})
	if err != nil {
		log.Printf("failed to add cron job: %v", err)
	}
}

func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("scheduler started")
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("scheduler stopped")
}