package core

import (
	"context"

	tgt "github.com/vkidmode/telegram_tree"
)

func (t *Service) generateProfileRoot() tgt.Node {
	return tgt.NewNode(
		tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("Админка"))),
		tgt.WithProc(processorWrap(t.generatePassengerNodes)),
	)
}

func (t *Service) generateProfileNodes(_ context.Context, _ MetaData) ([]tgt.Node, error) {
	return []tgt.Node{
		//tgt.NewNode(
		//	tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("Статистика"))),
		//	tgt.WithProc(processorWrap(t.generateAdminStatisticsNodes)),
		//),
		//tgt.NewNode(
		//	tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("Генерировать промокод"))),
		//	tgt.WithProc(processorWrap(t.generateAdminPromoNodes)),
		//),
		//tgt.NewNode(
		//	tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("Рекламные компании"))),
		//	tgt.WithProc(processorWrap(t.generateAdminAddNodes)),
		//),
		//tgt.NewNode(
		//	tgt.WithTg(
		//		tgt.NewTelegram(
		//			tgt.WithTabTxt("Удалить пользователя"),
		//			tgt.WithSwitchInline(fmt.Sprintf(customMessageFormat, entities.DeleteUserText)),
		//		),
		//	),
		//),
		//tgt.NewNode(
		//	tgt.WithTg(
		//		tgt.NewTelegram(
		//			tgt.WithTabTxt("Удалить сервер"),
		//			tgt.WithSwitchInline(fmt.Sprintf(customMessageFormat, entities.DeleteServerText)),
		//		),
		//	),
		//),
		//tgt.NewNode(
		//	tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("Рассылка"))),
		//	tgt.WithProc(processorWrap(t.generateNotificationsNodes)),
		//),
	}, nil
}
