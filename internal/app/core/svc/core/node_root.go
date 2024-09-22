package core

import (
	"context"
	"database/sql"
	"errors"

	tgt "github.com/vkidmode/telegram_tree"

	"github.com/tarmalonchik/falaemae/internal/entities"
	"github.com/tarmalonchik/falaemae/internal/pkg/trace"
)

var symbolsToNum = map[string]int{
	entities.AdminRoot:        0,
	entities.ServerRoot:       1,
	entities.PaymentRoot:      2,
	"d":                       3,
	"e":                       4,
	entities.InstructionsRoot: 5,
	"g":                       6,
	entities.ReviewsRoot:      7,
	entities.AddDemoRoot:      8,
	"j":                       9,
	"k":                       10,
	"l":                       11,
	"m":                       12,
	"n":                       13,
	"o":                       14,
	"p":                       15,
	"q":                       16,
	"r":                       17,
	"s":                       18,
	"t":                       19,
	"u":                       20,
	"v":                       21,
	"w":                       22,
	"x":                       23,
	"y":                       24,
	"z":                       25,
}

const (
	DaysKey                         = "."
	PromoKey                        = ","
	MonthKey                        = ":"
	CompaniesKey                    = ";"
	PaymentKey                      = "!"
	AddDemoKey                      = "`"
	ServersKey                      = "#"
	ReviewWithForHaveSubscription   = "?"
	ReviewWithForHaveNoSubscription = "%"
	messageSubscriptionExpired      = "У вас нет активной подписки. Чтобы подписаться, воспользуйтесь кнопками:"
	messageHaveActiveSubscription   = "Ваша подписка активна до: <b>%s</b> \nДля продления подписки воспользуйтесь кнопками:"
)

type processorFunc func(ctx context.Context, meta MetaData) ([]tgt.Node, error)

type MetaData struct {
	ChatID    int64
	MessageID int64
	Callback  string
	IsMiddle  bool
	Message   string
}

func (m *MetaData) GetCallback() string     { return m.Callback }
func (m *MetaData) SetupCallback(in string) { m.Callback = in }
func (m *MetaData) SetIsMiddle(middle bool) { m.IsMiddle = middle }
func (m *MetaData) GetIsMiddle() bool       { return m.IsMiddle }

func isCloseOrIgnoreCallback(in string) bool {
	if in == tgt.CallBackClose || in == tgt.CallBackIgnore {
		return true
	}
	return false
}

func isCloseCallback(in string) bool {
	if in == tgt.CallBackClose {
		return true
	}
	return false
}

func (t *Service) generateRootNodes(ctx context.Context, chatID int64) []tgt.Node {
	return []tgt.Node{
		//t.generateAdminRoot(),                   // 0 AdminRoot
		//t.generateServerRoot(),                  // 1 ServerRoot
		//t.generatePaymentRoot(ctx, chatID),      // 2 PaymentRoot
		t.reserved2Root(),    // 3 not reserved
		t.reservedNodeRoot(), // 4 not reserved
		//t.generateInstructionsRoot(ctx, chatID), // 5 InstructionsRoot
		//t.generatePaymentRoot(ctx, chatID),      // 6 Deprecated NewUsersRoot
		//t.generateReviewsRoot(),                 // 7 ReviewsRoot
		//t.generateAddDemoRoot(),                 // 8 AddDemoRoot
	}
}

func (*Service) reservedNodeRoot() tgt.Node {
	return tgt.NewNode(
		tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("Settings"), tgt.DeleteMsg())),
	)
}

func (*Service) reserved2Root() tgt.Node {
	return tgt.NewNode(
		tgt.WithTg(tgt.NewTelegram(tgt.WithTabTxt("s"))),
	)
}

func (t *Service) genServersBranch() commandsProcessorFunc {
	return func(ctx context.Context, meta MetaData) error {
		return t.ProcessCommandAsCallback(ctx, &MetaData{
			Callback:  entities.ServerRoot,
			ChatID:    meta.ChatID,
			MessageID: meta.MessageID,
		})
	}
}

func (t *Service) genAdminBranch() commandsProcessorFunc {
	return func(ctx context.Context, meta MetaData) error {
		return t.ProcessCommandAsCallback(ctx, &MetaData{
			Callback:  entities.AdminRoot,
			ChatID:    meta.ChatID,
			MessageID: meta.MessageID,
		})
	}
}

func (t *Service) genPaymentBranch() commandsProcessorFunc {
	return func(ctx context.Context, meta MetaData) error {
		return t.ProcessCommandAsCallback(ctx, &MetaData{
			Callback:  entities.PaymentRoot,
			ChatID:    meta.ChatID,
			MessageID: meta.MessageID,
		})
	}
}

func (t *Service) genInstructionsBranch() commandsProcessorFunc {
	return func(ctx context.Context, meta MetaData) error {
		return t.ProcessCommandAsCallback(ctx, &MetaData{
			Callback:  entities.InstructionsRoot,
			ChatID:    meta.ChatID,
			MessageID: meta.MessageID,
		})
	}
}

func (t *Service) genSendSupportBranch() commandsProcessorFunc {
	return func(ctx context.Context, info MetaData) error {
		message := "Напишите сюда, мы вам обязательно поможем 😊 @vpnchik_support 👨‍💻"
		_ = t.telegramClient.SendMessage(sql.NullInt64{Valid: true, Int64: info.ChatID}, message)
		return nil
	}
}

func processorWrap(in processorFunc) tgt.ProcessorFunc {
	return func(ctx context.Context, treeMeta tgt.Meta) ([]tgt.Node, error) {
		infoMeta, ok := treeMeta.(*MetaData)
		if !ok {
			return nil, trace.FuncNameWithError(errors.New("interface conversion"))
		}
		return in(ctx, *infoMeta)
	}
}

func (t *Service) ProcessCommandAsCallback(ctx context.Context, info *MetaData) error {
	if err := t.commonProcessor(ctx, info); err != nil {
		return trace.FuncNameWithErrorMsg(err, "processing command as callback")
	}
	return nil
}
