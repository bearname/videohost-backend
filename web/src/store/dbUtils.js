const DB_NAME = "authDB";
const DB_VERSION = 1;
let DB;

export default {
    async getDb() {
        return new Promise((resolve, reject) => {
            if (DB) {
                return resolve(DB);
            }
            let request = window.indexedDB.open(DB_NAME, DB_VERSION);

            request.onerror = e => {
                console.log(e)
                reject("Error");
            };

            request.onsuccess = e => {
                DB = e.target.result;
                resolve(DB);
            };

            request.onupgradeneeded = e => {
                let db = e.target.result;
                db.createObjectStore("userAuth", {
                    autoIncrement: true,
                    keyPath: "username"
                });
            };
        });
    },
    async addUser(user) {
        let db = await this.getDb();

        return new Promise(resolve => {
            let trans = db.transaction(["userAuth"], "readwrite");
            trans.oncomplete = () => {
                resolve();
            };

            let store = trans.objectStore("userAuth");
            store.openCursor().onsuccess = e => {
                let cursor = e.target.result;
                if (cursor) {
                    if (cursor.value.username !== user.username) {
                        store.delete(cursor.value.username);
                    }
                    cursor.continue();
                }
            };
            store.put(user);
        });
    },
    async removeUser(user) {
        if (!user || !user.username) {
            return;
        }
        let db = await this.getDb();

        return new Promise(resolve => {
            let trans = db.transaction(["userAuth"], "readwrite");
            trans.oncomplete = () => {
                resolve();
            };

            let store = trans.objectStore("userAuth");
            store.delete(user.username);
        });
    },
    async getUser() {
        let db = await this.getDb();

        return new Promise(resolve => {
            let trans = db.transaction(["userAuth"], "readwrite");
            trans.oncomplete = () => {
                resolve(user);
            };
            let user = null;
            let store = trans.objectStore("userAuth");
            store.openCursor().onsuccess = e => {
                let cursor = e.target.result;

                if (!cursor) return;
                user = cursor.value;
            };
        });
    }
};