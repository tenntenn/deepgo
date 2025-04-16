package style

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/gostaticanalysis/skeleton/v2/skeleton"
	"github.com/josharian/txtarfs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tenntenn/goversion"
	"golang.org/x/mod/module"
	"golang.org/x/tools/txtar"

	"github.com/tenntenn/deepgo/toolutil"
)

const skeletonDescription = `
The tool generates skeleton code of linter using golang.org/x/tools/go/analysis.Analyzer with github.com/gostaticanalysis/skeleton.
The tool responses generated codes as txtar format.
The dst parameter represents destination of files which must be absoluted path. When the dst parameter is empty string, the tool return generated code as txtar format string.

You MUST read the notes on developing your linters with the Skeleton below.

[NOTE]
There are several levels of static analysis in Go.
You should choose the parameter kind which must be matched the level.
Parsing and inspecting abstract syntax trees is the simplest (kind=inspect).
The ASTs are passed through analysis.Pass in Analyzer.Run.

However, a node of an AST has no type information.
The object corresponding to an identifier cannot be obtained from a node of the AST.
This information can be obtained through the TypesInfo field and the Pkg field of Analyzer.Pass.

If you want to check the code flow, you can use the Static Single Assign (SSA) form (kind=ssa).
Because the variable value can change through multiple flow paths, it is difficult to get the value of the variable without a control flow graph.
You can use the Static Single Assign (SSA) form with golang.org/x/tools/go/analysis/passes/buildssa.

You should choose the best static analysis method to reduce false positives and false negatives.

If you want to modify skeleton codes to complete your linter, at first you should consider test data in testdata directory.
You should apply edge cases as much as possible to your test data.
`

func NewSkeleton() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("style/skeleton",
			mcp.WithDescription(strings.TrimSpace(skeletonDescription)),
			mcp.WithString("kind",
				mcp.Description("The parameter kind (inspect,ssa,codegen,packages) is kind of skeleton code."),
			),
			mcp.WithString("module",
				mcp.Required(),
				mcp.Description("The parameter module is a module path of skeleton code."),
			),
			mcp.WithString("dst",
				mcp.Required(),
				mcp.Description("The dst parameter represents destination of files which must be absoluted path. When the dst parameter is empty string, the tool return generated code as txtar format string"),
			),
		),
		Handler: handleSkeleton,
	}
}

func handleSkeleton(ctx context.Context, request mcp.CallToolRequest) (result *mcp.CallToolResult, rerr error) {

	var kind skeleton.Kind
	kindStr, ok := request.Params.Arguments["kind"].(string)
	if ok {
		if err := kind.Set(kindStr); err != nil {
			return mcp.NewToolResultError("invalid kind (inspect,ssa,codegen,packages)"), nil
		}
	} else {
		kind = skeleton.KindInspect
	}

	modulepath, ok := request.Params.Arguments["module"].(string)
	if !ok {
		return mcp.NewToolResultError("The parameter module must be specified."), nil
	}

	if err := module.CheckPath(modulepath); err != nil {
		return mcp.NewToolResultError("The parameter module is invalid module path in Go.: module.CheckPath returns error: " + err.Error()), nil
	}

	latestGo, err := goversion.FetchLatest(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the latest Go version: %w", err)
	}

	info := &skeleton.Info{
		Kind:      kind,
		Checker:   skeleton.CheckerUnit,
		Path:      modulepath,
		Pkg:       path.Base(modulepath),
		Cmd:       true,
		GoMod:     true,
		GoVersion: latestGo.Version,
	}

	g := &skeleton.Generator{
		Template: skeleton.DefaultTemplate,
	}

	fsys, err := g.Run(info)
	if err != nil {
		return nil, fmt.Errorf("failed to generate skeleton code: %w", err)
	}

	ar, err := txtarfs.From(fsys)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to txtar format: %w", err)
	}

	dst, ok := request.Params.Arguments["dst"].(string)
	if ok && dst != "" {
		if err := toolutil.CopyTxtar(dst, string(txtar.Format(ar))); err != nil {
			return nil, fmt.Errorf("failed to copy files from txtar format generated skeleton code")
		}
		msg := fmt.Sprintf("The skeleton code has been generated in  %s.", dst)
		return mcp.NewToolResultText(msg), nil
	}

	return mcp.NewToolResultText(string(txtar.Format(ar))), nil
}
