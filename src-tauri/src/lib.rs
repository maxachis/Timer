mod timer;

use std::sync::Mutex;
use serde::{Deserialize, Serialize};
use tauri::Manager;
use tauri_plugin_store::StoreExt;
use timer::CountdownTimer;

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

#[derive(Serialize)]
pub struct TimerStatus {
    remaining_secs: f64,
    state: String,
    is_finished: bool,
}

#[tauri::command]
fn get_timer_status(timer: tauri::State<Mutex<CountdownTimer>>) -> TimerStatus {
    let t = timer.lock().unwrap();
    TimerStatus {
        remaining_secs: t.remaining().as_secs_f64(),
        state: t.state_name().to_string(),
        is_finished: t.is_finished(),
    }
}

#[tauri::command]
fn start_timer(timer: tauri::State<Mutex<CountdownTimer>>) {
    timer.lock().unwrap().start();
}

#[tauri::command]
fn pause_timer(timer: tauri::State<Mutex<CountdownTimer>>) {
    timer.lock().unwrap().pause();
}

#[tauri::command]
fn resume_timer(timer: tauri::State<Mutex<CountdownTimer>>) {
    timer.lock().unwrap().resume();
}

#[tauri::command]
fn reset_timer(timer: tauri::State<Mutex<CountdownTimer>>) {
    timer.lock().unwrap().reset();
}

#[tauri::command]
fn add_time(timer: tauri::State<Mutex<CountdownTimer>>, settings: tauri::State<Mutex<AppSettings>>) {
    let increment = settings.lock().unwrap().default_increment_secs;
    timer.lock().unwrap().add_time(Some(std::time::Duration::from_secs(increment)));
}

#[tauri::command]
fn remove_time(timer: tauri::State<Mutex<CountdownTimer>>, settings: tauri::State<Mutex<AppSettings>>) {
    let increment = settings.lock().unwrap().default_increment_secs;
    timer.lock().unwrap().remove_time(Some(std::time::Duration::from_secs(increment)));
}

#[tauri::command]
fn add_time_secondary(timer: tauri::State<Mutex<CountdownTimer>>, settings: tauri::State<Mutex<AppSettings>>) {
    let increment = settings.lock().unwrap().secondary_increment_secs;
    timer.lock().unwrap().add_time(Some(std::time::Duration::from_secs(increment)));
}

#[tauri::command]
fn remove_time_secondary(timer: tauri::State<Mutex<CountdownTimer>>, settings: tauri::State<Mutex<AppSettings>>) {
    let increment = settings.lock().unwrap().secondary_increment_secs;
    timer.lock().unwrap().remove_time(Some(std::time::Duration::from_secs(increment)));
}

#[tauri::command]
fn add_time_custom(timer: tauri::State<Mutex<CountdownTimer>>, seconds: u64) {
    timer.lock().unwrap().add_time(Some(std::time::Duration::from_secs(seconds)));
}

#[tauri::command]
fn remove_time_custom(timer: tauri::State<Mutex<CountdownTimer>>, seconds: u64) {
    timer.lock().unwrap().remove_time(Some(std::time::Duration::from_secs(seconds)));
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

#[tauri::command]
fn create_timer(duration_secs: u64, timer: tauri::State<Mutex<CountdownTimer>>) {
    let mut t = timer.lock().unwrap();
    *t = CountdownTimer::new(duration_secs);
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_notification::init())
        .plugin(tauri_plugin_store::Builder::default().build())
        .setup(|app| {
            let settings = load_settings(app.handle());
            let duration = settings.default_duration_secs;
            app.manage(Mutex::new(settings));
            app.manage(Mutex::new(CountdownTimer::new(duration)));
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            get_timer_status,
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
            get_settings,
            save_settings,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
