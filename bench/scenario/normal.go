package scenario

import (
	"context"
	"time"

	"github.com/chibiegg/isucon9-final/bench/internal/bencherror"
	"github.com/chibiegg/isucon9-final/bench/internal/config"
	"github.com/chibiegg/isucon9-final/bench/internal/xrandom"
	"github.com/chibiegg/isucon9-final/bench/isutrain"
	"github.com/chibiegg/isucon9-final/bench/payment"
)

// NormalScenario は基本的な予約フローのシナリオです
func NormalScenario(ctx context.Context) error {
	client, err := isutrain.NewClient()
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewCriticalError(err, "Isutrainクライアントが作成できません. 運営に確認をお願いいたします"))
	}

	paymentClient, err := payment.NewClient()
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewCriticalError(err, "Paymentクライアントが作成できません. 運営に確認をお願いいたします"))
	}

	if config.Debug {
		client.ReplaceMockTransport()
	}

	err = client.Signup(ctx, "hoge", "hoge", nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "ユーザ登録ができません"))
	}

	err = client.Login(ctx, "hoge", "hoge", nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "ユーザログインができません"))
	}

	_, err = client.ListStations(ctx, nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "駅一覧を取得できません"))
	}

	var (
		departure = xrandom.GetRandomStations()
		arrival   = xrandom.GetRandomStations()
	)
	_, err = client.SearchTrains(ctx, time.Now().AddDate(1, 0, 0), departure, arrival, nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "列車検索ができません"))
	}

	_, err = client.ListTrainSeats(ctx, time.Now().AddDate(1, 0, 0), "こだま", "96号", 1, "東京", "大阪", nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "列車の座席座席列挙できません"))
	}

	reservation, err := client.Reserve(ctx, "こだま", "69号", "premium", isutrain.TrainSeats{}, "東京", "名古屋", time.Now().AddDate(1, 0, 0), 1, 1, 1, "isle", nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "予約ができません"))
	}

	cardToken, err := paymentClient.RegistCard(ctx, "11111111", "222", "10/50")
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewCriticalError(err, "クレジットカードの登録ができません. 運営にご確認をお願いいたします"))
	}

	err = client.CommitReservation(ctx, reservation.ReservationID, cardToken, nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "予約を確定できませんでした"))
	}

	_, err = client.ListReservations(ctx, nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "予約を列挙できませんでした"))
	}

	reservation2, err := client.ShowReservation(ctx, reservation.ReservationID, nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "予約詳細を取得できませんでした"))
	}

	if reservation.ReservationID != reservation2.ID {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "正しい予約詳細を取得できませんでした"))
	}

	if err := client.Logout(ctx, nil); err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "ログアウトできません"))
	}

	return nil
}

// 予約キャンセル含む(Commit後にキャンセル)
func NormalCancelScenario(ctx context.Context) error {
	client, err := isutrain.NewClient()
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewCriticalError(err, "Isutrainクライアントが作成できません. 運営に確認をお願いいたします"))
	}

	paymentClient, err := payment.NewClient()
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewCriticalError(err, "Paymentクライアントが作成できません. 運営に確認をお願いいたします"))
	}

	if config.Debug {
		client.ReplaceMockTransport()
	}

	err = client.Signup(ctx, "hoge", "hoge", nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "ユーザ登録ができません"))
	}

	err = client.Login(ctx, "hoge", "hoge", nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "ユーザログインができません"))
	}

	_, err = client.ListStations(ctx, nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "駅一覧を取得できません"))
	}

	var (
		departure = xrandom.GetRandomStations()
		arrival   = xrandom.GetRandomStations()
	)
	_, err = client.SearchTrains(ctx, time.Now().AddDate(1, 0, 0), departure, arrival, nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "列車検索ができません"))
	}

	_, err = client.ListTrainSeats(ctx, time.Now().AddDate(1, 0, 0), "こだま", "96号", 1, "東京", "大阪", nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "列車の座席座席列挙できません"))
	}

	reservation, err := client.Reserve(ctx, "こだま", "69号", "premium", isutrain.TrainSeats{}, "東京", "名古屋", time.Now().AddDate(1, 0, 0), 1, 1, 1, "isle", nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "予約ができません"))
	}

	cardToken, err := paymentClient.RegistCard(ctx, "11111111", "222", "20/20")
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewCriticalError(err, "クレジットカードの登録ができませんでした. 運営にご確認をお願いいたします"))
	}

	err = client.CommitReservation(ctx, reservation.ReservationID, cardToken, nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "予約を確定できませんでした"))
	}

	_, err = client.ListReservations(ctx, nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "予約を列挙できませんでした"))
	}

	err = client.CancelReservation(ctx, reservation.ReservationID, nil)
	if err != nil {
		return bencherror.BenchmarkErrs.AddError(bencherror.NewApplicationError(err, "予約をキャンセルできませんでした"))
	}

	return nil
}

// 曖昧検索シナリオ
func NormalAmbigiousSearchScenario(ctx context.Context) error {
	return nil
}

// オリンピック開催期間に負荷をあげるシナリオ
// FIXME: initializeで指定された日数に応じ、負荷レベルを変化させる
func NormalOlympicParticipantsScenario(ctx context.Context) error {

	// NOTE: webappから見て、明らかに負荷が上がったと感じるレベルに持ってく必要がある
	// NOTE: 指定できる最大の日数で負荷をかける際、飽和しないようにする
	// NOTE:

	return nil
}