export namespace engin {
	
	export class ProgressMsg {
	    name: string;
	    msg: string;
	    percent: number;
	    time: string;
	
	    static createFrom(source: any = {}) {
	        return new ProgressMsg(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.msg = source["msg"];
	        this.percent = source["percent"];
	        this.time = source["time"];
	    }
	}
	export class ProgressResult {
	    name: string;
	    msg: string;
	    percent: number;
	    time: string;
	    err: string;
	    result: string;
	
	    static createFrom(source: any = {}) {
	        return new ProgressResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.msg = source["msg"];
	        this.percent = source["percent"];
	        this.time = source["time"];
	        this.err = source["err"];
	        this.result = source["result"];
	    }
	}

}

