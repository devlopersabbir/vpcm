export namespace config {
	
	export class APIConfig {
	    Enabled: boolean;
	    Host: string;
	    Port: number;
	    Mode: string;
	    Token: string;
	    GlobalURL: string;
	
	    static createFrom(source: any = {}) {
	        return new APIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Enabled = source["Enabled"];
	        this.Host = source["Host"];
	        this.Port = source["Port"];
	        this.Mode = source["Mode"];
	        this.Token = source["Token"];
	        this.GlobalURL = source["GlobalURL"];
	    }
	}
	export class CollectorConfig {
	    Workers: number;
	
	    static createFrom(source: any = {}) {
	        return new CollectorConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Workers = source["Workers"];
	    }
	}
	export class PluginsConfig {
	    Enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PluginsConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Enabled = source["Enabled"];
	    }
	}
	export class LoggingConfig {
	    Level: string;
	    Format: string;
	
	    static createFrom(source: any = {}) {
	        return new LoggingConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Level = source["Level"];
	        this.Format = source["Format"];
	    }
	}
	export class SSHConfig {
	    Timeout: number;
	
	    static createFrom(source: any = {}) {
	        return new SSHConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Timeout = source["Timeout"];
	    }
	}
	export class DatabaseConfig {
	    Driver: string;
	    Path: string;
	    URI: string;
	    Name: string;
	
	    static createFrom(source: any = {}) {
	        return new DatabaseConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Driver = source["Driver"];
	        this.Path = source["Path"];
	        this.URI = source["URI"];
	        this.Name = source["Name"];
	    }
	}
	export class Config {
	    Database: DatabaseConfig;
	    API: APIConfig;
	    SSH: SSHConfig;
	    Logging: LoggingConfig;
	    Collector: CollectorConfig;
	    Plugins: PluginsConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Database = this.convertValues(source["Database"], DatabaseConfig);
	        this.API = this.convertValues(source["API"], APIConfig);
	        this.SSH = this.convertValues(source["SSH"], SSHConfig);
	        this.Logging = this.convertValues(source["Logging"], LoggingConfig);
	        this.Collector = this.convertValues(source["Collector"], CollectorConfig);
	        this.Plugins = this.convertValues(source["Plugins"], PluginsConfig);
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

export namespace inventory {
	
	export class ConnectionLog {
	    id: number;
	    server_id: number;
	    server_name: string;
	    username: string;
	    host: string;
	    // Go type: time
	    logged_in_at: any;
	    // Go type: time
	    logged_out_at?: any;
	    duration?: string;
	    status: string;
	    error_message?: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionLog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.server_id = source["server_id"];
	        this.server_name = source["server_name"];
	        this.username = source["username"];
	        this.host = source["host"];
	        this.logged_in_at = this.convertValues(source["logged_in_at"], null);
	        this.logged_out_at = this.convertValues(source["logged_out_at"], null);
	        this.duration = source["duration"];
	        this.status = source["status"];
	        this.error_message = source["error_message"];
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
	export class ServerHardware {
	    id: number;
	    server_id: number;
	    cpu_model: string;
	    cpu_cores: number;
	    ram_total: string;
	    swap_total: string;
	    disk_total: string;
	    virtualization: string;
	    instance_type: string;
	    serial_number: string;
	    bios_version: string;
	    uptime: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerHardware(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.server_id = source["server_id"];
	        this.cpu_model = source["cpu_model"];
	        this.cpu_cores = source["cpu_cores"];
	        this.ram_total = source["ram_total"];
	        this.swap_total = source["swap_total"];
	        this.disk_total = source["disk_total"];
	        this.virtualization = source["virtualization"];
	        this.instance_type = source["instance_type"];
	        this.serial_number = source["serial_number"];
	        this.bios_version = source["bios_version"];
	        this.uptime = source["uptime"];
	    }
	}
	export class ServerNetwork {
	    id: number;
	    server_id: number;
	    hostname: string;
	    public_ip: string;
	    private_ip: string;
	    mac_address: string;
	    region: string;
	    availability_zone: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerNetwork(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.server_id = source["server_id"];
	        this.hostname = source["hostname"];
	        this.public_ip = source["public_ip"];
	        this.private_ip = source["private_ip"];
	        this.mac_address = source["mac_address"];
	        this.region = source["region"];
	        this.availability_zone = source["availability_zone"];
	    }
	}
	export class ServerOS {
	    id: number;
	    server_id: number;
	    os_family: string;
	    os_version: string;
	    kernel_version: string;
	    architecture: string;
	    init_system: string;
	    timezone: string;
	    locale: string;
	    package_manager: string;
	
	    static createFrom(source: any = {}) {
	        return new ServerOS(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.server_id = source["server_id"];
	        this.os_family = source["os_family"];
	        this.os_version = source["os_version"];
	        this.kernel_version = source["kernel_version"];
	        this.architecture = source["architecture"];
	        this.init_system = source["init_system"];
	        this.timezone = source["timezone"];
	        this.locale = source["locale"];
	        this.package_manager = source["package_manager"];
	    }
	}
	export class Software {
	    id: number;
	    server_id: number;
	    name: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new Software(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.server_id = source["server_id"];
	        this.name = source["name"];
	        this.version = source["version"];
	    }
	}
	export class Tag {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Tag(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class ServerView {
	    id: number;
	    uuid: string;
	    name: string;
	    host: string;
	    port: number;
	    username: string;
	    auth_type: string;
	    auth_secret?: string;
	    provider: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    // Go type: time
	    last_seen?: any;
	    is_favorite: boolean;
	    tags: Tag[];
	    network?: ServerNetwork;
	    hardware?: ServerHardware;
	    os?: ServerOS;
	    software: Software[];
	
	    static createFrom(source: any = {}) {
	        return new ServerView(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.uuid = source["uuid"];
	        this.name = source["name"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.auth_type = source["auth_type"];
	        this.auth_secret = source["auth_secret"];
	        this.provider = source["provider"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.last_seen = this.convertValues(source["last_seen"], null);
	        this.is_favorite = source["is_favorite"];
	        this.tags = this.convertValues(source["tags"], Tag);
	        this.network = this.convertValues(source["network"], ServerNetwork);
	        this.hardware = this.convertValues(source["hardware"], ServerHardware);
	        this.os = this.convertValues(source["os"], ServerOS);
	        this.software = this.convertValues(source["software"], Software);
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
	
	
	export class TerminalPreference {
	    id: number;
	    font_size: number;
	    font_family: string;
	    background: string;
	    foreground: string;
	    opacity: number;
	    blur: number;
	    cursor_style: string;
	    cursor_blink: boolean;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new TerminalPreference(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.font_size = source["font_size"];
	        this.font_family = source["font_family"];
	        this.background = source["background"];
	        this.foreground = source["foreground"];
	        this.opacity = source["opacity"];
	        this.blur = source["blur"];
	        this.cursor_style = source["cursor_style"];
	        this.cursor_blink = source["cursor_blink"];
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	
	export class SSHConnectionParams {
	    host: string;
	    port: number;
	    username: string;
	    auth_type: string;
	    auth_secret: string;
	    rows: number;
	    cols: number;
	
	    static createFrom(source: any = {}) {
	        return new SSHConnectionParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.host = source["host"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.auth_type = source["auth_type"];
	        this.auth_secret = source["auth_secret"];
	        this.rows = source["rows"];
	        this.cols = source["cols"];
	    }
	}

}

