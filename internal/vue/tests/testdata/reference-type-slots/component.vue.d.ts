declare var __VLS_1: {}, __VLS_3: {
    num: number;
}, __VLS_5: {
    str: string;
}, __VLS_7: {
    num: number;
    str: string;
};
type __VLS_Slots = {} & {
    'no-bind'?: (props: typeof __VLS_1) => any;
} & {
    default?: (props: typeof __VLS_3) => any;
} & {
    'named-slot'?: (props: typeof __VLS_5) => any;
} & {
    vbind?: (props: typeof __VLS_7) => any;
};
declare const __VLS_base: import("vue").DefineComponent2<{
    setup(): {};
    data(): {};
    props: {};
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
declare const __VLS_export: __VLS_WithSlots<typeof __VLS_base, __VLS_Slots>;
declare const _default: typeof __VLS_export;
export default _default;
type __VLS_WithSlots<T, S> = T & {
    new (): {
        $slots: S;
    };
};
