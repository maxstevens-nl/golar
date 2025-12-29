declare const __VLS_export: import("vue").DefineComponent2<{
    setup(): {};
    data(): {};
    props: {
        foo: {
            type: StringConstructor;
            required: true;
        };
        bar: {
            type: StringConstructor;
            default: string;
        };
        baz: {
            type: StringConstructor;
        };
        xfoo: {
            readonly type: StringConstructor;
            readonly required: true;
        };
        xbar: {
            readonly type: StringConstructor;
            readonly default: "";
        };
        xbaz: {
            readonly type: StringConstructor;
        };
        /**
         * The hello property.
         *
         * @since v1.0.0
         */
        hello: {
            type: StringConstructor;
            default: string;
        };
        numberOrStringProp: {
            type: (StringConstructor | NumberConstructor)[];
            default: number;
        };
        arrayProps: {
            type: ArrayConstructor;
            default: () => number[];
        };
    };
    computed: {};
    methods: {};
    mixins: {}[];
    extends: {};
    emits: string[];
    slots: {};
    inject: {};
    components: {};
    directives: {};
    provide: {};
    expose: string;
    __typeProps: unknown;
    __typeEmits: unknown;
    __typeRefs: {};
    __typeEl: any;
    __defaults: unknown;
}>;
declare const _default: typeof __VLS_export;
export default _default;
