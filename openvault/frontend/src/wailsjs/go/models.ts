export namespace fs {
	
	export class InitOptions {
	    FirstName: string;
	    LastName: string;
	    Email: string;
	    Password: string;
	
	    static createFrom(source: any = {}) {
	        return new InitOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FirstName = source["FirstName"];
	        this.LastName = source["LastName"];
	        this.Email = source["Email"];
	        this.Password = source["Password"];
	    }
	}

}

export namespace main {
	
	export class AccountWithUnlockStatus {
	    id: string;
	    user_email: string;
	    user_first_name: string;
	    user_last_name: string;
	    // Go type: cryptolib
	    secret_key?: any;
	    is_unlocked: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AccountWithUnlockStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.user_email = source["user_email"];
	        this.user_first_name = source["user_first_name"];
	        this.user_last_name = source["user_last_name"];
	        this.secret_key = this.convertValues(source["secret_key"], null);
	        this.is_unlocked = source["is_unlocked"];
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
	export class DecryptedVaultItemDetails {
	    item_id: string;
	    vault_id: string;
	    created_at: string;
	    updated_at: string;
	    // Go type: cryptolib
	    encrypted_details?: any;
	    username: string;
	    password: string;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new DecryptedVaultItemDetails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.item_id = source["item_id"];
	        this.vault_id = source["vault_id"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.encrypted_details = this.convertValues(source["encrypted_details"], null);
	        this.username = source["username"];
	        this.password = source["password"];
	        this.notes = source["notes"];
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
	export class DecryptedVaultItemOverview {
	    item_id: string;
	    vault_id: string;
	    created_at: string;
	    updated_at: string;
	    // Go type: cryptolib
	    encrypted_overview?: any;
	    title: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new DecryptedVaultItemOverview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.item_id = source["item_id"];
	        this.vault_id = source["vault_id"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.encrypted_overview = this.convertValues(source["encrypted_overview"], null);
	        this.title = source["title"];
	        this.url = source["url"];
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

export namespace structs {
	
	export class VaultMetadata {
	    vault_id: string;
	    name: string;
	    description: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new VaultMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vault_id = source["vault_id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}

}

