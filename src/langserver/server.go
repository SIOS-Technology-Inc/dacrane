package langserver

import (
	"errors"
	"fmt"

	"github.com/SIOS-Technology-Inc/dacrane/v0/src/ast"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/exception"
	"github.com/SIOS-Technology-Inc/dacrane/v0/src/parser"
	"github.com/macrat/simplexer"
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

var Files = map[string]string{}

func Start() {
	// This increases logging verbosity (optional)
	commonlog.Configure(2, nil)

	handler = protocol.Handler{
		Initialize:                     initialize,
		Initialized:                    initialized,
		Shutdown:                       shutdown,
		SetTrace:                       setTrace,
		TextDocumentCompletion:         TextDocumentCompletion,
		TextDocumentSemanticTokensFull: TextDocumentSemanticTokensFull,
		TextDocumentDidOpen:            TextDocumentDidOpen,
		TextDocumentDidChange:          TextDocumentDidChange,
		TextDocumentDidSave:            TextDocumentDidSave,
	}

	server := server.NewServer(&handler, lsName, true)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	commonlog.NewInfoMessage(0, "Initializing server...")

	capabilities := handler.CreateServerCapabilities()

	capabilities.TextDocumentSync = protocol.TextDocumentSyncKindFull
	capabilities.CompletionProvider = &protocol.CompletionOptions{}
	capabilities.SemanticTokensProvider = &protocol.SemanticTokensOptions{
		Legend: protocol.SemanticTokensLegend{
			TokenTypes:     tokenTypes,
			TokenModifiers: []string{},
		},
		Full: protocol.True,
	}

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

func TextDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	text := params.TextDocument.Text
	uri := params.TextDocument.URI
	Files[uri] = text
	SendNotifications(context, uri, text)
	return nil
}

func TextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	text := params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole).Text
	uri := params.TextDocument.URI
	Files[uri] = text
	SendNotifications(context, uri, text)
	return nil
}

func TextDocumentDidSave(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	return nil
}

func SendNotifications(context *glsp.Context, uri string, text string) {
	var codeErr *exception.CodeError
	codeErrors := []exception.CodeError{}

	tokens, err := parser.Lex(text)
	if errors.As(err, &codeErr) {
		codeErrors = append(codeErrors, *codeErr)
		SendCodeError(context, uri, codeErrors)
		return
	}

	m, err := parser.Parse(tokens)
	if errors.As(err, &codeErr) {
		codeErrors = append(codeErrors, *codeErr)
		SendCodeError(context, uri, codeErrors)
		return
	}

	for _, v := range m.Vars {
		_, err := v.Expr.Infer(m.Vars)
		if errors.As(err, &codeErr) {
			codeErrors = append(codeErrors, *codeErr)
		}
	}

	SendCodeError(context, uri, codeErrors)
}

func TextDocumentSemanticTokensFull(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	text := Files[params.TextDocument.URI]

	tokens, err := parser.Lex(text)
	if err != nil {
		return nil, nil
	}

	return CreateSemanticTokens(tokens)
}

var tokenTypes = []string{"number", "string", "operator", "variable"}

func TokenKindNumber(token *simplexer.Token) (uint32, error) {
	switch token.Type.GetID() {
	case parser.INTEGER:
		return 0, nil
	case parser.STRING:
		return 1, nil
	case parser.ASSIGN:
		return 2, nil
	case parser.ADD:
		return 2, nil
	case parser.LBRACKET:
		return 2, nil
	case parser.RBRACKET:
		return 2, nil
	case parser.IDENTIFIER:
		return 3, nil
	default:
		return 99, fmt.Errorf("cannot mapping token kind: token id (%d)", token.Type.GetID())
	}
}

func CreateSemanticTokens(tokens []*simplexer.Token) (*protocol.SemanticTokens, error) {
	data := []uint32{}
	if len(tokens) == 0 {
		return &protocol.SemanticTokens{
			Data: data,
		}, nil
	}
	previousToken := tokens[0]
	for _, t := range tokens {
		tokenKindNumber, err := TokenKindNumber(t)
		if err != nil {
			return nil, err
		}
		dColumn := 0
		if t.Position.Line == previousToken.Position.Line {
			dColumn = t.Position.Column - previousToken.Position.Column
		}

		data = append(data,
			uint32(t.Position.Line-previousToken.Position.Line),
			uint32(dColumn),
			uint32(len(t.Literal)),
			tokenKindNumber,
			0, //
		)
		previousToken = t
	}
	return &protocol.SemanticTokens{
		Data: data,
	}, nil
}

func SendCodeError(context *glsp.Context, uri string, codeErrors []exception.CodeError) {
	diagnostics := []protocol.Diagnostic{}
	for _, codeError := range codeErrors {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(codeError.Range.Start.Line), Character: uint32(codeError.Range.Start.Column)},
				End:   protocol.Position{Line: uint32(codeError.Range.End.Line), Character: uint32(codeError.Range.End.Column)},
			},
			Message: codeError.Message,
		})
	}

	context.Notify("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
}

func SendNoError(context *glsp.Context, uri string) {
	context.Notify("textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: []protocol.Diagnostic{},
	})
}
