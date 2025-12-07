import Dot from "./dot";
import DotContext from "./dot-context";

export default class DotKernel<T> {
    private dotValues: Map<Dot, T>;
    private dotContext: DotContext;

    constructor() {
        this.dotValues = new Map<Dot, T>();
        this.dotContext = new DotContext();
    }

    public getContext(): DotContext {
        return this.dotContext;
    }

    public setContext(context: DotContext): void {
        this.dotContext = context;
    }

    public dotAdd(id: string, value: T): Dot {
        const dot = this.dotContext.makeDot(id);
        this.dotValues.set(dot, value);
        return dot;
    }

    public ddd(id: string, value: T): DotKernel<T> {
        const dot = this.dotAdd(id, value);

        const delta = new DotKernel<T>();
        delta.dotValues.set(dot, value);
        delta.dotContext.insertDot(dot);

        return delta
    }

    public removeDot(dot: Dot): DotKernel<T> {
        const delta = new DotKernel<T>();

        if (this.dotValues.has(dot)) {
            this.dotValues.delete(dot);
            delta.dotContext.insertDot(dot);
        }

        return delta;
    }

    // Removes any dot with a matching value
    public removeValue(value: T): DotKernel<T> {
        const delta = new DotKernel<T>();

        this.dotValues.forEach((dotValue, dot) => {
            if (dotValue === value) {
                this.dotValues.delete(dot);
                delta.dotContext.insertDot(dot, false);
            }
        })

        delta.dotContext.compact();
        return delta;
    }

    public reset(): DotKernel<T> {
        const delta = new DotKernel<T>();

        this.dotValues.forEach((_, dot) => {
            delta.dotContext.insertDot(dot, false);
        });

        delta.dotContext.compact();
        this.dotValues.clear();
        return delta;
    }

    // Merges the kernel with another, preferring the values of other on conflicts
    public join(other: DotKernel<T>): void {
        this.dotValues.forEach((value, dot) => {
            // If dot is not present in other and is known by other, it means it was removed
            if (!other.dotContext.knows(dot) && !other.dotValues.has(dot)) {
                this.dotValues.delete(dot);
            }
        });

        other.dotValues.forEach((value, dot) => {
            // If dot is not present locally, add it
            if (!this.dotValues.has(dot)) {
                this.dotValues.set(dot, value);
            }
        });

        this.dotContext.join(other.dotContext);
    }

    public clone(): DotKernel<T> {
        const clone = new DotKernel<T>();

        clone.dotContext = this.dotContext.clone();
        clone.dotValues = new Map<Dot, T>(this.dotValues);

        return clone;
    }
}