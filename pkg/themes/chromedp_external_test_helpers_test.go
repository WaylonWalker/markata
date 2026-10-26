package themes_test

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

func evaluateInto[T any](expression string, result *T) chromedp.Action[chromedp.Void] {
	evaluate := chromedp.Evaluate[T](expression, chromedp.EvalAwaitPromise)
	return chromedp.Func(func(ctx context.Context, target *chromedp.Target) error {
		value, err := evaluate(ctx, target)
		if err != nil {
			return err
		}
		*result = value
		return nil
	})
}

func evaluateVoid(expression string) chromedp.Action[chromedp.Void] {
	return chromedp.Evaluate[chromedp.Void](expression, chromedp.EvalAwaitPromise)
}

func pollVoid(predicate string, opts ...chromedp.PollOption) chromedp.Action[chromedp.Void] {
	return chromedp.Poll[chromedp.Void](predicate, opts...)
}

func setDeviceMetricsOverride(width, height int64, scale float64, mobile bool) chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, target *chromedp.Target) error {
		_, err := cdp.Call(ctx, target, emulation.SetDeviceMetricsOverride, emulation.SetDeviceMetricsOverrideParams{
			Width:             width,
			Height:            height,
			DeviceScaleFactor: scale,
			Mobile:            mobile,
		})
		return err
	})
}
