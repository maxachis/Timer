mod timer;

use std::sync::Mutex;
use serde::{Deserialize, Serialize};
use tauri::Manager;
use tauri_plugin_store::StoreExt;
use timer::TimerCollection;

const SETTINGS_STORE: &str = "settings.json";

#[derive(Serialize, Deserialize, Clone)]
pub struct AppSettings {
    default_duration_secs: u64,
    default_increment_secs: u64,
    secondary_increment_secs: u64,
}

impl Default for AppSettings {
    fn default() -> Self {
        Self {
            default_duration_secs: 300,
            default_increment_secs: 300,
            secondary_increment_secs: 60,
        }
    }
}

const DURATION_MIN: u64 = 60;
const DURATION_MAX: u64 = 10800;
const INCREMENT_MIN: u64 = 60;
const INCREMENT_MAX: u64 = 3600;

fn validate_settings(s: &AppSettings) -> Result<(), String> {
    if !(DURATION_MIN..=DURATION_MAX).contains(&s.default_duration_secs) {
        return Err("default_duration_secs must be between 60 and 10800".into());
    }
    if !(INCREMENT_MIN..=INCREMENT_MAX).contains(&s.default_increment_secs) {
        return Err("default_increment_secs must be between 60 and 3600".into());
    }
    if !(INCREMENT_MIN..=INCREMENT_MAX).contains(&s.secondary_increment_secs) {
        return Err("secondary_increment_secs must be between 60 and 3600".into());
    }
    Ok(())
}

fn parse_settings_from<F>(get: F) -> AppSettings
where
    F: Fn(&str) -> Option<serde_json::Value>,
{
    let defaults = AppSettings::default();
    let duration = get("default_duration_secs")
        .and_then(|v| v.as_u64())
        .filter(|&v| (DURATION_MIN..=DURATION_MAX).contains(&v))
        .unwrap_or(defaults.default_duration_secs);
    let increment = get("default_increment_secs")
        .and_then(|v| v.as_u64())
        .filter(|&v| (INCREMENT_MIN..=INCREMENT_MAX).contains(&v))
        .unwrap_or(defaults.default_increment_secs);
    let secondary_increment = get("secondary_increment_secs")
        .and_then(|v| v.as_u64())
        .filter(|&v| (INCREMENT_MIN..=INCREMENT_MAX).contains(&v))
        .unwrap_or(defaults.secondary_increment_secs);
    AppSettings {
        default_duration_secs: duration,
        default_increment_secs: increment,
        secondary_increment_secs: secondary_increment,
    }
}

fn parse_timer_names_from<F>(get: F) -> Vec<String>
where
    F: Fn(&str) -> Option<serde_json::Value>,
{
    get("timer_names")
        .and_then(|v| serde_json::from_value::<Vec<String>>(v).ok())
        .unwrap_or_default()
}

fn load_settings(app: &tauri::AppHandle) -> AppSettings {
    let Some(store) = app.store(SETTINGS_STORE).ok() else {
        return AppSettings::default();
    };
    parse_settings_from(|k| store.get(k))
}

fn load_timer_names(app: &tauri::AppHandle) -> Vec<String> {
    let Some(store) = app.store(SETTINGS_STORE).ok() else {
        return vec![];
    };
    parse_timer_names_from(|k| store.get(k))
}

fn save_timer_names(app: &tauri::AppHandle, names: &[String]) {
    if let Ok(store) = app.store(SETTINGS_STORE) {
        store.set("timer_names", serde_json::to_value(names).unwrap());
    }
}

#[derive(Serialize)]
pub struct TimerStatus {
    remaining_secs: f64,
    state: String,
    is_finished: bool,
    active_index: usize,
    active_name: String,
    timer_count: usize,
}

#[derive(Serialize)]
pub struct TimerListEntry {
    index: usize,
    name: String,
    state: String,
    remaining_secs: f64,
    is_active: bool,
}

#[tauri::command]
fn get_timer_status(timers: tauri::State<Mutex<TimerCollection>>) -> TimerStatus {
    let col = timers.lock().unwrap();
    let t = col.active();
    TimerStatus {
        remaining_secs: t.remaining().as_secs_f64(),
        state: t.state_name().to_string(),
        is_finished: t.is_finished(),
        active_index: col.active_index(),
        active_name: col.active_name().to_string(),
        timer_count: col.count(),
    }
}

#[tauri::command]
fn get_timer_list(timers: tauri::State<Mutex<TimerCollection>>) -> Vec<TimerListEntry> {
    let col = timers.lock().unwrap();
    col.timer_list()
        .into_iter()
        .map(|info| TimerListEntry {
            index: info.index,
            name: info.name,
            state: info.state,
            remaining_secs: info.remaining_secs,
            is_active: info.is_active,
        })
        .collect()
}

#[tauri::command]
fn start_timer(timers: tauri::State<Mutex<TimerCollection>>) {
    timers.lock().unwrap().active_mut().start();
}

#[tauri::command]
fn pause_timer(timers: tauri::State<Mutex<TimerCollection>>) {
    timers.lock().unwrap().active_mut().pause();
}

#[tauri::command]
fn resume_timer(timers: tauri::State<Mutex<TimerCollection>>) {
    timers.lock().unwrap().active_mut().resume();
}

#[tauri::command]
fn reset_timer(timers: tauri::State<Mutex<TimerCollection>>) {
    timers.lock().unwrap().active_mut().reset();
}

#[tauri::command]
fn add_time(timers: tauri::State<Mutex<TimerCollection>>, settings: tauri::State<Mutex<AppSettings>>) {
    let increment = settings.lock().unwrap().default_increment_secs;
    timers.lock().unwrap().active_mut().add_time(Some(std::time::Duration::from_secs(increment)));
}

#[tauri::command]
fn remove_time(timers: tauri::State<Mutex<TimerCollection>>, settings: tauri::State<Mutex<AppSettings>>) {
    let increment = settings.lock().unwrap().default_increment_secs;
    timers.lock().unwrap().active_mut().remove_time(Some(std::time::Duration::from_secs(increment)));
}

#[tauri::command]
fn add_time_secondary(timers: tauri::State<Mutex<TimerCollection>>, settings: tauri::State<Mutex<AppSettings>>) {
    let increment = settings.lock().unwrap().secondary_increment_secs;
    timers.lock().unwrap().active_mut().add_time(Some(std::time::Duration::from_secs(increment)));
}

#[tauri::command]
fn remove_time_secondary(timers: tauri::State<Mutex<TimerCollection>>, settings: tauri::State<Mutex<AppSettings>>) {
    let increment = settings.lock().unwrap().secondary_increment_secs;
    timers.lock().unwrap().active_mut().remove_time(Some(std::time::Duration::from_secs(increment)));
}

#[tauri::command]
fn add_time_custom(timers: tauri::State<Mutex<TimerCollection>>, seconds: u64) {
    timers.lock().unwrap().active_mut().add_time(Some(std::time::Duration::from_secs(seconds)));
}

#[tauri::command]
fn remove_time_custom(timers: tauri::State<Mutex<TimerCollection>>, seconds: u64) {
    timers.lock().unwrap().active_mut().remove_time(Some(std::time::Duration::from_secs(seconds)));
}

#[tauri::command]
fn create_timer(duration_secs: u64, timers: tauri::State<Mutex<TimerCollection>>) {
    let mut col = timers.lock().unwrap();
    *col.active_mut() = timer::CountdownTimer::new(duration_secs);
}

#[tauri::command]
fn add_new_timer(
    name: String,
    timers: tauri::State<Mutex<TimerCollection>>,
    settings: tauri::State<Mutex<AppSettings>>,
    app: tauri::AppHandle,
) -> Result<usize, String> {
    let duration = settings.lock().unwrap().default_duration_secs;
    let mut col = timers.lock().unwrap();
    let idx = col.add_timer(name, duration).map_err(|e| e.to_string())?;
    save_timer_names(&app, &col.timer_names());
    Ok(idx)
}

#[tauri::command]
fn remove_existing_timer(
    index: usize,
    timers: tauri::State<Mutex<TimerCollection>>,
    app: tauri::AppHandle,
) -> Result<(), String> {
    let mut col = timers.lock().unwrap();
    col.remove_timer(index).map_err(|e| e.to_string())?;
    save_timer_names(&app, &col.timer_names());
    Ok(())
}

#[tauri::command]
fn switch_timer(
    index: usize,
    timers: tauri::State<Mutex<TimerCollection>>,
) -> Result<(), String> {
    timers.lock().unwrap().switch_to(index).map_err(|e| e.to_string())
}

#[tauri::command]
fn switch_timer_next(timers: tauri::State<Mutex<TimerCollection>>) {
    timers.lock().unwrap().switch_next();
}

#[tauri::command]
fn switch_timer_prev(timers: tauri::State<Mutex<TimerCollection>>) {
    timers.lock().unwrap().switch_prev();
}

#[tauri::command]
fn rename_timer(
    index: usize,
    name: String,
    timers: tauri::State<Mutex<TimerCollection>>,
    app: tauri::AppHandle,
) -> Result<(), String> {
    let mut col = timers.lock().unwrap();
    col.rename_timer(index, name).map_err(|e| e.to_string())?;
    save_timer_names(&app, &col.timer_names());
    Ok(())
}

#[tauri::command]
fn get_settings(settings: tauri::State<Mutex<AppSettings>>) -> AppSettings {
    settings.lock().unwrap().clone()
}

#[tauri::command]
fn save_settings(
    new_settings: AppSettings,
    settings: tauri::State<Mutex<AppSettings>>,
    app: tauri::AppHandle,
) -> Result<(), String> {
    validate_settings(&new_settings)?;

    let store = app.store(SETTINGS_STORE).map_err(|e| e.to_string())?;
    store.set("default_duration_secs", new_settings.default_duration_secs);
    store.set("default_increment_secs", new_settings.default_increment_secs);
    store.set("secondary_increment_secs", new_settings.secondary_increment_secs);

    let mut s = settings.lock().unwrap();
    *s = new_settings;

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::{json, Value};
    use std::collections::HashMap;

    fn make_get(map: HashMap<&'static str, Value>) -> impl Fn(&str) -> Option<Value> {
        move |k: &str| map.get(k).cloned()
    }

    #[test]
    fn default_settings_values() {
        let d = AppSettings::default();
        assert_eq!(d.default_duration_secs, 300);
        assert_eq!(d.default_increment_secs, 300);
        assert_eq!(d.secondary_increment_secs, 60);
    }

    #[test]
    fn validate_accepts_boundaries() {
        assert!(validate_settings(&AppSettings {
            default_duration_secs: 60,
            default_increment_secs: 60,
            secondary_increment_secs: 60,
        }).is_ok());
        assert!(validate_settings(&AppSettings {
            default_duration_secs: 10800,
            default_increment_secs: 3600,
            secondary_increment_secs: 3600,
        }).is_ok());
    }

    #[test]
    fn validate_rejects_duration_below_min() {
        let r = validate_settings(&AppSettings {
            default_duration_secs: 59,
            default_increment_secs: 300,
            secondary_increment_secs: 60,
        });
        assert!(r.is_err());
        assert!(r.unwrap_err().contains("default_duration_secs"));
    }

    #[test]
    fn validate_rejects_duration_above_max() {
        let r = validate_settings(&AppSettings {
            default_duration_secs: 10801,
            default_increment_secs: 300,
            secondary_increment_secs: 60,
        });
        assert!(r.is_err());
    }

    #[test]
    fn validate_rejects_increment_out_of_range() {
        let r = validate_settings(&AppSettings {
            default_duration_secs: 300,
            default_increment_secs: 3601,
            secondary_increment_secs: 60,
        });
        assert!(r.is_err());
        assert!(r.unwrap_err().contains("default_increment_secs"));
    }

    #[test]
    fn validate_rejects_secondary_increment_out_of_range() {
        let r = validate_settings(&AppSettings {
            default_duration_secs: 300,
            default_increment_secs: 300,
            secondary_increment_secs: 59,
        });
        assert!(r.is_err());
        assert!(r.unwrap_err().contains("secondary_increment_secs"));
    }

    #[test]
    fn parse_settings_empty_returns_defaults() {
        let s = parse_settings_from(make_get(HashMap::new()));
        let d = AppSettings::default();
        assert_eq!(s.default_duration_secs, d.default_duration_secs);
        assert_eq!(s.default_increment_secs, d.default_increment_secs);
        assert_eq!(s.secondary_increment_secs, d.secondary_increment_secs);
    }

    #[test]
    fn parse_settings_reads_valid_stored_values() {
        let mut m = HashMap::new();
        m.insert("default_duration_secs", json!(600));
        m.insert("default_increment_secs", json!(120));
        m.insert("secondary_increment_secs", json!(90));
        let s = parse_settings_from(make_get(m));
        assert_eq!(s.default_duration_secs, 600);
        assert_eq!(s.default_increment_secs, 120);
        assert_eq!(s.secondary_increment_secs, 90);
    }

    #[test]
    fn parse_settings_falls_back_when_out_of_range() {
        let mut m = HashMap::new();
        m.insert("default_duration_secs", json!(1));
        m.insert("default_increment_secs", json!(99999));
        m.insert("secondary_increment_secs", json!(0));
        let s = parse_settings_from(make_get(m));
        let d = AppSettings::default();
        assert_eq!(s.default_duration_secs, d.default_duration_secs);
        assert_eq!(s.default_increment_secs, d.default_increment_secs);
        assert_eq!(s.secondary_increment_secs, d.secondary_increment_secs);
    }

    #[test]
    fn parse_settings_falls_back_on_wrong_type() {
        let mut m = HashMap::new();
        m.insert("default_duration_secs", json!("not a number"));
        let s = parse_settings_from(make_get(m));
        assert_eq!(s.default_duration_secs, AppSettings::default().default_duration_secs);
    }

    #[test]
    fn parse_timer_names_empty() {
        let names = parse_timer_names_from(make_get(HashMap::new()));
        assert!(names.is_empty());
    }

    #[test]
    fn parse_timer_names_valid() {
        let mut m = HashMap::new();
        m.insert("timer_names", json!(["A", "B", "C"]));
        let names = parse_timer_names_from(make_get(m));
        assert_eq!(names, vec!["A", "B", "C"]);
    }

    #[test]
    fn parse_timer_names_rejects_wrong_type() {
        let mut m = HashMap::new();
        m.insert("timer_names", json!("not an array"));
        let names = parse_timer_names_from(make_get(m));
        assert!(names.is_empty());
    }

    #[test]
    fn app_settings_serializes_with_expected_field_names() {
        let s = AppSettings::default();
        let v = serde_json::to_value(&s).unwrap();
        assert!(v.get("default_duration_secs").is_some());
        assert!(v.get("default_increment_secs").is_some());
        assert!(v.get("secondary_increment_secs").is_some());
    }

    #[test]
    fn timer_status_serializes_with_expected_shape() {
        let s = TimerStatus {
            remaining_secs: 12.5,
            state: "Running".into(),
            is_finished: false,
            active_index: 0,
            active_name: "Timer 1".into(),
            timer_count: 1,
        };
        let v = serde_json::to_value(&s).unwrap();
        assert_eq!(v["remaining_secs"], json!(12.5));
        assert_eq!(v["state"], json!("Running"));
        assert_eq!(v["is_finished"], json!(false));
        assert_eq!(v["active_index"], json!(0));
        assert_eq!(v["active_name"], json!("Timer 1"));
        assert_eq!(v["timer_count"], json!(1));
    }

    #[test]
    fn timer_list_entry_serializes_with_expected_shape() {
        let e = TimerListEntry {
            index: 2,
            name: "Work".into(),
            state: "Paused".into(),
            remaining_secs: 60.0,
            is_active: true,
        };
        let v = serde_json::to_value(&e).unwrap();
        assert_eq!(v["index"], json!(2));
        assert_eq!(v["name"], json!("Work"));
        assert_eq!(v["state"], json!("Paused"));
        assert_eq!(v["remaining_secs"], json!(60.0));
        assert_eq!(v["is_active"], json!(true));
    }

    #[test]
    fn app_settings_roundtrip() {
        let s = AppSettings {
            default_duration_secs: 600,
            default_increment_secs: 120,
            secondary_increment_secs: 90,
        };
        let v = serde_json::to_value(&s).unwrap();
        let back: AppSettings = serde_json::from_value(v).unwrap();
        assert_eq!(back.default_duration_secs, 600);
        assert_eq!(back.default_increment_secs, 120);
        assert_eq!(back.secondary_increment_secs, 90);
    }

    #[test]
    fn settings_store_key_is_stable() {
        assert_eq!(SETTINGS_STORE, "settings.json");
    }
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_notification::init())
        .plugin(tauri_plugin_store::Builder::default().build())
        .setup(|app| {
            let settings = load_settings(app.handle());
            let timer_names = load_timer_names(app.handle());
            let duration = settings.default_duration_secs;
            app.manage(Mutex::new(settings));
            app.manage(Mutex::new(TimerCollection::from_names(&timer_names, duration)));
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            get_timer_status,
            get_timer_list,
            start_timer,
            pause_timer,
            resume_timer,
            reset_timer,
            add_time,
            remove_time,
            add_time_secondary,
            remove_time_secondary,
            add_time_custom,
            remove_time_custom,
            create_timer,
            add_new_timer,
            remove_existing_timer,
            switch_timer,
            switch_timer_next,
            switch_timer_prev,
            rename_timer,
            get_settings,
            save_settings,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
