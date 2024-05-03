package langserver

import (
	"errors"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/exception"
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
	// capabilities.SemanticTokensProvider = &protocol.SemanticTokensOptions{
	// 	Legend: protocol.SemanticTokensLegend{
	// 		TokenTypes:     []string{"number", "string", "operator"},
	// 		TokenModifiers: []string{},
	// 	},
	// 	Full: protocol.True,
	// }

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

func TextDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	var completionItems []protocol.CompletionItem

	operator := protocol.CompletionItemKindOperator
	for _, f := range ast.FixtureFunctions {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: f.Name,
			Kind:  &operator,
		})
	}

	return completionItems, nil
}

func TextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	text := params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole).Text
	var codeErr *exception.CodeError

	tokens, err := parser.Lex(text)
	if errors.As(err, &codeErr) {
		SendCodeError(context, params.TextDocument.URI, *codeErr)
		return nil
	}

	expr, err := parser.Parse(tokens)
	if errors.As(err, &codeErr) {
		SendCodeError(context, params.TextDocument.URI, *codeErr)
		return nil
	}

	_, err = expr.Evaluate()
	if errors.As(err, &codeErr) {
		SendCodeError(context, params.TextDocument.URI, *codeErr)
		return nil
	}
	SendNoError(context, params.TextDocument.URI)

	return nil
}

func SendCodeError(context *glsp.Context, uri string, codeError exception.CodeError) {
	context.Notify("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
		URI: uri,
		Diagnostics: []protocol.Diagnostic{
			{
				Range: protocol.Range{
					Start: protocol.Position{Line: uint32(codeError.Range.Start.Line), Character: uint32(codeError.Range.Start.Column)},
					End:   protocol.Position{Line: uint32(codeError.Range.End.Line), Character: uint32(codeError.Range.End.Column)},
				},
				Message: codeError.Message,
			},
		},
	})
}

func SendNoError(context *glsp.Context, uri string) {
	context.Notify("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: []protocol.Diagnostic{},
	})
}
