/**
 * Emitted when the foo property is changed.
 */
type MyEventsFoo = 'foo';
export interface MyEvents {
    (event: MyEventsFoo, data?: {
        foo: string;
    }): void;
    (event: 'bar', value: {
        arg1: number;
        arg2?: any;
    }): void;
    (e: 'baz'): void;
}
export {};
