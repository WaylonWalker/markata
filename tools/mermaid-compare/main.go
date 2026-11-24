// A standalone dependency-isolated benchmark, NOT a Markata production renderer.
// Reads JSONL and compares two Go SVG engines in the same warm process.
package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "time"

    mmdg "github.com/bvolpato/mermaid-go-renderer"
    gomermaid "github.com/zkrebbekx/go-mermaid"
)

type request struct {
    Source string `json:"source"`
    Iterations int `json:"iterations"`
    Font string `json:"font"`
    FontSize float64 `json:"font_size"`
    Background string `json:"background"`
    Fill string `json:"fill"`
    Stroke string `json:"stroke"`
    Text string `json:"text"`
    Edge string `json:"edge"`
}
type candidate struct {
    Backend string `json:"backend"`
    SVG string `json:"svg,omitempty"`
    NS []int64 `json:"ns"`
    Error string `json:"error,omitempty"`
}
type response struct {
    Results []candidate `json:"results"`
}

func (r request) normalized() request {
    if r.Iterations < 1 { r.Iterations = 1 }
    if r.Iterations > 100 { r.Iterations = 100 }
    if r.Font == "" { r.Font = "sans-serif" }
    if r.FontSize <= 0 { r.FontSize = 16 }
    if r.Background == "" { r.Background = "#ffffff" }
    if r.Fill == "" { r.Fill = "#f9fafb" }
    if r.Stroke == "" { r.Stroke = "#237b5e" }
    if r.Text == "" { r.Text = "#1f2937" }
    if r.Edge == "" { r.Edge = r.Text }
    return r
}
func renderMMDG(r request) (string, error) {
    opts := mmdg.DefaultRenderOptions()
    theme := opts.Theme
    theme.Background = r.Background
    theme.PrimaryColor = r.Fill
    theme.PrimaryBorderColor = r.Stroke
    theme.PrimaryTextColor = r.Text
    theme.LineColor = r.Edge
    theme.TextColor = r.Text
    theme.EdgeLabelBackground = r.Background
    theme.ClusterBorder = r.Stroke
    theme.FontFamily = r.Font
    theme.FontSize = r.FontSize
    opts.Theme = theme
    return mmdg.RenderWithOptions(r.Source, opts)
}
func renderGoMermaid(r request) (string, error) {
    pal := gomermaid.Palette{Background:r.Background, NodeFill:r.Fill,
        NodeStroke:r.Stroke, Text:r.Text, Edge:r.Edge}
    svg, err := gomermaid.Render(r.Source,
        gomermaid.WithCustomTheme("markata-comparison", pal),
        gomermaid.WithFont(r.Font, r.FontSize),
    )
    return string(svg), err
}
func measure(name string, r request, fn func(request)(string,error)) candidate {
    out := candidate{Backend:name, NS:make([]int64,0,r.Iterations)}
    // Warmup is excluded. No subprocess, network call or JSON included in timings.
    for i:=0;i<2;i++ {
        if _,err:=fn(r);err!=nil {out.Error=err.Error();return out}
    }
    for i:=0;i<r.Iterations;i++ {
        t:=time.Now()
        svg,err:=fn(r)
        out.NS=append(out.NS,time.Since(t).Nanoseconds())
        if err!=nil {out.Error=err.Error();return out}
        if !strings.Contains(svg,"<svg") {out.Error="output is not SVG";return out}
        if i==0 {out.SVG=svg}
    }
    return out
}
func main() {
    scan:=bufio.NewScanner(os.Stdin)
    scan.Buffer(make([]byte,4096),2<<20)
    out:=bufio.NewWriter(os.Stdout)
    defer out.Flush()
    for scan.Scan() {
        var req request
        if err:=json.Unmarshal(scan.Bytes(),&req);err!=nil {
            fmt.Fprintln(out,`{"error":"invalid input JSON"}`);out.Flush();continue
        }
        req=req.normalized()
        resp:=response{Results:[]candidate{
            measure("mmdg",req,renderMMDG),
            measure("go-mermaid",req,renderGoMermaid),
        }}
        if err:=json.NewEncoder(out).Encode(resp);err!=nil {fmt.Fprintln(os.Stderr,err);os.Exit(1)}
        out.Flush()
    }
    if err:=scan.Err();err!=nil {fmt.Fprintln(os.Stderr,err);os.Exit(1)}
}
