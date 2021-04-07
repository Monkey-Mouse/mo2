export class LazyExecutor {
    private i = 0;
    private f: () => void;
    private delay = 0;
    constructor(f?: () => void, delay: number = 200) {
        this.f = f;
        this.delay = delay;
    }
    /**
     * Execute
     */
    public Execute(f?: () => void,) {
        this.i++;
        const num = this.i;
        setTimeout(() => {
            if (num === this.i) {
                if (f) {
                    f()
                } else this.f();
            }
        }, 200);
    }
}
export class SlowExecutor {
    private i = 0;
    private f: () => void;
    private delay = 0;
    private exe = false;
    constructor(f?: () => void, delay: number = 200) {
        this.f = f;
        this.delay = delay;
        setInterval(() => {
            if (this.exe) {
                this.f()
                this.exe = false;
            }
        }, delay)
    }
    /**
     * Execute
     */
    public Execute(f?: () => void,) {
        if (f) {
            this.f = f;
        }
        this.exe = true;
    }
}