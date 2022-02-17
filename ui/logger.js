

function Logger() {
    this.server = '0.0.0.0'
    this.location = ''
    this.restOfArgs = {}
    this.send = async function (...args) {
        const data = JSON.stringify(...args)
        if (!this.server) return this
        try {
            await fetch('http://' + this.server + '/send', {
                method: 'POST',
                headers: {},
                body: JSON.stringify({
                    data: data
                })
            })
        } catch (e) {
            console.log('e', e)
        }
        return this
    }
    this.init = function ({ server, location, ...args }) {
        this.server = server
        this.location = location
        this.restOfArgs = { ...args }
    }
    this.logAll = function (name, args) {
        console.log(name, args);
        this.send({ tag: name, location: this.location, ...this.restOfArgs, args: args });
    }
    this.enableLogAll = function (obj) {
        return new Proxy(obj, {
            get(target, p) {
                if (p in target) {
                    return target[p];
                } else if (typeof target.logAll == "function") {
                    return function (...args) {
                        return target.logAll.call(target, p, args);
                    };
                }
            }
        });
    }

    return this.enableLogAll(this);
}

const logger = new Logger();


