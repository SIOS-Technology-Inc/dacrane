package langserver

import (
	"errors"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/parser"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	_ "github.com/tliron/commonlog/simple"
)

const lsName = "Dacrane Language Server"

var (
	version string = "0.0.1"
	handler protocol.Handler
)

func Start() {
	// This increases logging verbosity (optional)
	commonlog.Configure(2, nil)

	handler = protocol.Handler{
		Initialize:             initialize,
		Initialized:            initialized,
		Shutdown:               shutdown,
		SetTrace:               setTrace,
		TextDocumentCompletion: TextDocumentCompletion,
		TextDocumentDidChange:  TextDocumentDidChange,
	}

	server := server.NewServer(&handler, lsName, true)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	commonlog.NewInfoMessage(0, "Initializing server...")

	capabilities := handler.CreateServerCapabilities()

	capabilities.TextDocumentSync = protocol.TextDocumentSyncKindFull
	capabilities.CompletionProvider = &protocol.CompletionOptions{}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

var Completions = map[string]string{
	"a": "A",
	"b": "B",
}

func TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	var completionItems []protocol.CompletionItem

	for word, emoji := range Completions {
		emojiCopy := emoji // Create a copy of emoji
		completionItems = append(completionItems, protocol.CompletionItem{
			Label:      word,
			Detail:     &emojiCopy,
			InsertText: &emojiCopy,
		})
	}

	return completionItems, nil
}

func TextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	text := params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole).Text
	expr := parser.Parse(text)
	_, err := expr.Evaluate()
	var evalErr *ast.EvalError
	if errors.As(err, &evalErr) {
		context.Notify("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
			URI: params.TextDocument.URI,
			Diagnostics: []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{Line: uint32(evalErr.Position.Line), Character: uint32(evalErr.Position.Column)},
						End:   protocol.Position{Line: uint32(evalErr.Position.Line), Character: uint32(evalErr.Position.Column + 1)},
					},
					Message: evalErr.Message,
				},
			},
		})
	} else {
		context.Notify("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: []protocol.Diagnostic{},
		})
	}
	return nil
}
