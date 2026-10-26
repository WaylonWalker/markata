package themes

import (
	"context"

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
