export namespace history {
	
	export class Entry {
	    // Go type: time
	    at: any;
	    kind: string;
	    timer_name: string;
	
	    static createFrom(source: any = {}) {
	        return new Entry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.at = this.convertValues(source["at"], null);
	        this.kind = source["kind"];
	        this.timer_name = source["timer_name"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace main {
	
	export class TimerStatus {
	    remaining_secs: number;
	    state: string;
	    is_finished: boolean;
	    active_index: number;
	    active_name: string;
	    timer_count: number;
	    on_break: boolean;
	    break_index: number;
	
	    static createFrom(source: any = {}) {
	        return new TimerStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.remaining_secs = source["remaining_secs"];
	        this.state = source["state"];
	        this.is_finished = source["is_finished"];
	        this.active_index = source["active_index"];
	        this.active_name = source["active_name"];
	        this.timer_count = source["timer_count"];
	        this.on_break = source["on_break"];
	        this.break_index = source["break_index"];
	    }
	}

}

export namespace settings {
	
	export class AppSettings {
	    default_duration_secs: number;
	    default_increment_secs: number;
	    secondary_increment_secs: number;
	    tertiary_increment_secs: number;
	    mini_window_width: number;
	    mini_window_height: number;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.default_duration_secs = source["default_duration_secs"];
	        this.default_increment_secs = source["default_increment_secs"];
	        this.secondary_increment_secs = source["secondary_increment_secs"];
	        this.tertiary_increment_secs = source["tertiary_increment_secs"];
	        this.mini_window_width = source["mini_window_width"];
	        this.mini_window_height = source["mini_window_height"];
	    }
	}

}

export namespace timer {
	
	export class TimerInfo {
	    index: number;
	    name: string;
	    state: string;
	    remaining_secs: number;
	    is_active: boolean;
	    tint: string;
	    original_duration_secs: number;
	    is_break_timer: boolean;
	    break_accrued_secs: number;
	
	    static createFrom(source: any = {}) {
	        return new TimerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.index = source["index"];
	        this.name = source["name"];
	        this.state = source["state"];
	        this.remaining_secs = source["remaining_secs"];
	        this.is_active = source["is_active"];
	        this.tint = source["tint"];
	        this.original_duration_secs = source["original_duration_secs"];
	        this.is_break_timer = source["is_break_timer"];
	        this.break_accrued_secs = source["break_accrued_secs"];
	    }
	}

}

