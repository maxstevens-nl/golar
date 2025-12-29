export declare const Foo: import("vue").DefineSetupFnComponent<{
    foo: string;
}, string[], {}, {
    foo: string;
} & {
    [x: `on${Capitalize<string>}`]: (...args: any[]) => any;
}, import("vue").PublicProps>;
export declare const Bar: import("vue").DefineSetupFnComponent<{
    bar?: number;
}, string[], {}, {
    bar?: number;
} & {
    [x: `on${Capitalize<string>}`]: (...args: any[]) => any;
}, import("vue").PublicProps>;
