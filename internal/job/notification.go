package job

import (
	"context"
	"time"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/notification"
)

type mailSenderJob struct {
	ctx        context.Context
	ticker     *time.Ticker
	cancel     context.CancelFunc
	db         *repository.Database
	mailSender notification.ISender
	templates  map[string]string
}

func NewMailSender(db *repository.Database, mailSender notification.ISender, templates map[string]string) Job {
	ctx, cancel := context.WithCancel(context.Background())

	return &mailSenderJob{
		ticker:     time.NewTicker(10 * time.Minute),
		ctx:        ctx,
		cancel:     cancel,
		db:         db,
		mailSender: mailSender,
		templates:  templates,
	}
}

func (p *mailSenderJob) Start() {
	p.sendVoucherNotifications()
	go p.runTicker()
}

func (p *mailSenderJob) sendVoucherNotifications() {
	notif, err := p.db.GetPendingMailNotification()
	if err != nil {
		log.Error("failed to get pending mail from db", "err", err)

		return
	}

	tmpl, err := notification.LoadMailTemplate(p.templates["voucher"])
	if err != nil {
		log.Fatal("failed to load mail template", "err", err)
	}

	err = p.mailSender.SendTemplateMail(
		notification.NotificationProviderZapToMail,
		"pagu@pactus.org", []string{notif.Recipient}, tmpl, notif.Data)
	if err != nil {
		log.Error("failed to send mail notification", "err", err)
		err = p.db.UpdateNotificationStatus(notif.ID, entity.NotificationStatusFail)
		if err != nil {
			log.Error("failed to update status of sent mail", "err", err)
		}
	} else {
		err = p.db.UpdateNotificationStatus(notif.ID, entity.NotificationStatusDone)
		if err != nil {
			log.Error("failed to update status of sent mail", "err", err)
		}
	}
}

func (p *mailSenderJob) runTicker() {
	for {
		select {
		case <-p.ctx.Done():
			return

		case <-p.ticker.C:
			p.sendVoucherNotifications()
		}
	}
}

func (p *mailSenderJob) Stop() {
	p.ticker.Stop()
}
