import * as React from 'react'

const LogContext = React.createContext()

function Logger() {
    this.failTime = 0
    this.fails = 0
    this.failLimit = 5 //number of times post fails before blocking 
    this.secsTillRetry = 5 * 60000 //seconds till retry after blocking
    this.server = '0.0.0.0'
    this.location = ''
    this.restOfArgs = {}
    this.send = async function (...args) {
        const data = JSON.stringify(...args)
        const date = new Date()
        if (this.fails >= this.failLimit && this.failTime + this.secsTillRetry > date.getTime()) return this
        if (!this.server) return this
        try {
            await fetch('http://' + this.server + '/send', {
                method: 'POST',
                headers: {},
                body: JSON.stringify({
                    data: data
                })
            })
            this.fails = 0
            this.failTime = 0
        } catch (e) {
            console.log('failed to send to log-server:', e)
            this.fails += 1
            this.failTime = date.getTime()
        }
        return this
    }
    this.init = function ({ server, location, ...args }) {
        this.server = server
        this.location = location
        this.restOfArgs = { ...args }
    }
    this.logAll = function (name, args) {
        // console.log(name, args);
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

function LogProvider({ children, ...args }) {
    const logger = new Logger();
    logger.init(args.configs)
    return <LogContext.Provider value={logger}>{children}</LogContext.Provider>
}

function useLog() {
    const context = React.useContext(LogContext)
    if (context === undefined) {
        throw new Error('useLog must be used within a LogProvider')
    }
    return context
}

export { LogProvider, useLog }