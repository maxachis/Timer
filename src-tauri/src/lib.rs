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

fn load_settings(app: &tauri::AppHandle) -> AppSettings {
    let store = app.store(SETTINGS_STORE).ok();
    let defaults = AppSettings::default();

    let Some(store) = store else {
        return defaults;
    };

    let duration = store
        .get("default_duration_secs")
        .and_then(|v| v.as_u64())
        .filter(|&v| (60..=10800).contains(&v))
        .unwrap_or(defaults.default_duration_secs);

    let increment = store
        .get("default_increment_secs")
        .and_then(|v| v.as_u64())
        .filter(|&v| (60..=3600).contains(&v))
        .unwrap_or(defaults.default_increment_secs);

    let secondary_increment = store
        .get("secondary_increment_secs")
        .and_then(|v| v.as_u64())
        .filter(|&v| (60..=3600).contains(&v))
        .unwrap_or(defaults.secondary_increment_secs);

    AppSettings {
        default_duration_secs: duration,
        default_increment_secs: increment,
        secondary_increment_secs: secondary_increment,
    }
}

fn load_timer_names(app: &tauri::AppHandle) -> Vec<String> {
    let store = app.store(SETTINGS_STORE).ok();
    let Some(store) = store else {
        return vec![];
    };
    store
        .get("timer_names")
        .and_then(|v| serde_json::from_value::<Vec<String>>(v.clone()).ok())
        .unwrap_or_default()
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
    if !(60..=10800).contains(&new_settings.default_duration_secs) {
        return Err("default_duration_secs must be between 60 and 10800".into());
    }
    if !(60..=3600).contains(&new_settings.default_increment_secs) {
        return Err("default_increment_secs must be between 60 and 3600".into());
    }
    if !(60..=3600).contains(&new_settings.secondary_increment_secs) {
        return Err("secondary_increment_secs must be between 60 and 3600".into());
    }

    let store = app.store(SETTINGS_STORE).map_err(|e| e.to_string())?;
    store.set("default_duration_secs", new_settings.default_duration_secs);
    store.set("default_increment_secs", new_settings.default_increment_secs);
    store.set("secondary_increment_secs", new_settings.secondary_increment_secs);

    let mut s = settings.lock().unwrap();
    *s = new_settings;

    Ok(())
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
