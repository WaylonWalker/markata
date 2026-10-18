package plugins

import "github.com/WaylonWalker/markata/pkg/lifecycle"

// CriticalStageErrors marks publication failures as build-fatal. Other
// Write-stage plugins keep the lifecycle's default warning-only behavior.
func (p *PublishHTMLPlugin) CriticalStageErrors(stage lifecycle.Stage) bool {
	return stage == lifecycle.StageWrite
}

var _ lifecycle.CriticalStageErrorsPlugin = (*PublishHTMLPlugin)(nil)
