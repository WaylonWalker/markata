use std::{
    fs,
    path::{Path, PathBuf},
    time::{Duration, Instant},
};

use excalidraw_render::{ColorMode, OutputFormat, RenderRequest, Renderer, SceneDocument};

const RUNS: usize = 6;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let fixture_dir = PathBuf::from(
        std::env::args()
            .nth(1)
            .ok_or("usage: renderer-spike FIXTURE_DIR OUTPUT_DIR")?,
    );
    let output_dir = PathBuf::from(
        std::env::args()
            .nth(2)
            .ok_or("usage: renderer-spike FIXTURE_DIR OUTPUT_DIR")?,
    );
    let only_fixture = std::env::var("EXCALIDRAW_SPIKE_FIXTURE").ok();
    fs::create_dir_all(&output_dir)?;
    let renderer = Renderer::new();
    let request = RenderRequest::builder(OutputFormat::Svg)
        .color_mode(ColorMode::Dark)
        .padding(20.0)
        .build()?;

    let mut files = fs::read_dir(&fixture_dir)?
        .filter_map(Result::ok)
        .map(|entry| entry.path())
        .filter(|path| {
            path.extension().is_some_and(|ext| ext == "excalidraw")
                || path
                    .file_name()
                    .is_some_and(|name| name == "malformed.json")
        })
        .filter(|path| {
            only_fixture
                .as_deref()
                .is_none_or(|name| path.file_name().and_then(|name| name.to_str()) == Some(name))
        })
        .collect::<Vec<_>>();
    files.sort();

    for file in files {
        render_fixture(&renderer, &request, &file, &output_dir)?;
    }
    println!("process_max_rss_kb={}", peak_rss_kb().unwrap_or(0));
    Ok(())
}

fn render_fixture(
    renderer: &Renderer,
    request: &RenderRequest,
    file: &Path,
    output_dir: &Path,
) -> Result<(), Box<dyn std::error::Error>> {
    let name = file
        .file_name()
        .and_then(|name| name.to_str())
        .unwrap_or("unknown");
    let source = fs::read(file)?;
    let scene = match SceneDocument::parse(&source) {
        Ok(scene) => scene,
        Err(error) => {
            println!("fixture={name} parse_error={error}");
            return Ok(());
        }
    };

    let mut durations = Vec::with_capacity(RUNS);
    let mut previous = None;
    let mut deterministic = true;
    let mut diagnostics = 0;
    let mut diagnostic_summary = Vec::new();
    let mut svg = Vec::new();
    for _ in 0..RUNS {
        let start = Instant::now();
        match renderer.render(&scene, request) {
            Ok(outcome) => {
                durations.push(start.elapsed());
                svg = outcome.artifact().bytes().to_vec();
                diagnostics = outcome.report().diagnostics().len();
                diagnostic_summary = outcome
                    .report()
                    .diagnostics()
                    .iter()
                    .map(|diagnostic| {
                        format!(
                            "{}@{}:{}",
                            diagnostic.code(),
                            diagnostic.json_pointer(),
                            diagnostic.message()
                        )
                    })
                    .collect();
                if let Some(prior) = previous.as_ref() {
                    deterministic &= prior == &svg;
                }
                previous = Some(svg.clone());
            }
            Err(error) => {
                println!("fixture={name} render_error={error}");
                return Ok(());
            }
        }
    }

    fs::write(output_dir.join(format!("{name}.svg")), &svg)?;
    let first = durations.first().copied().unwrap_or_default();
    let warm_median = median(&durations[1..]);
    println!(
        "fixture={name} first_ms={:.3} warm_median_ms={:.3} svg_bytes={} diagnostics={} diagnostic_summary={:?} deterministic={deterministic}",
        millis(first),
        millis(warm_median),
        svg.len(),
        diagnostics,
        diagnostic_summary
    );
    Ok(())
}

fn median(durations: &[Duration]) -> Duration {
    let mut sorted = durations.to_vec();
    sorted.sort_unstable();
    sorted[sorted.len() / 2]
}

fn millis(duration: Duration) -> f64 {
    duration.as_secs_f64() * 1000.0
}

fn peak_rss_kb() -> Option<u64> {
    let status = fs::read_to_string("/proc/self/status").ok()?;
    status.lines().find_map(|line| {
        line.strip_prefix("VmHWM:")
            .and_then(|value| value.split_whitespace().next())
            .and_then(|value| value.parse().ok())
    })
}
