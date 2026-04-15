export namespace main {
	
	export class TimerStatus {
	    remaining_secs: number;
	    state: string;
	    is_finished: boolean;
	    active_index: number;
	    active_name: string;
	    timer_count: number;
	
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
	    }
	}

}

export namespace settings {
	
	export class AppSettings {
	    default_duration_secs: number;
	    default_increment_secs: number;
	    secondary_increment_secs: number;
	    tertiary_increment_secs: number;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.default_duration_secs = source["default_duration_secs"];
	        this.default_increment_secs = source["default_increment_secs"];
	        this.secondary_increment_secs = source["secondary_increment_secs"];
	        this.tertiary_increment_secs = source["tertiary_increment_secs"];
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
	    }
	}

}

