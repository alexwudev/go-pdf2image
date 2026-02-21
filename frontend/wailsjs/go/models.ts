export namespace app {

	export class ConvertConfig {
	    dpi: number;
	    quality: number;
	    format: string;
	    pages: string;
	    output_dir: string;
	    workers: number;
	    zip_output: boolean;

	    static createFrom(source: any = {}) {
	        return new ConvertConfig(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dpi = source["dpi"];
	        this.quality = source["quality"];
	        this.format = source["format"];
	        this.pages = source["pages"];
	        this.output_dir = source["output_dir"];
	        this.workers = source["workers"];
	        this.zip_output = source["zip_output"];
	    }
	}
	export class ConvertResult {
	    output_files: string[];
	    error: string;

	    static createFrom(source: any = {}) {
	        return new ConvertResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.output_files = source["output_files"];
	        this.error = source["error"];
	    }
	}
	export class PDFInfo {
	    page_count: number;
	    error: string;

	    static createFrom(source: any = {}) {
	        return new PDFInfo(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page_count = source["page_count"];
	        this.error = source["error"];
	    }
	}

}

