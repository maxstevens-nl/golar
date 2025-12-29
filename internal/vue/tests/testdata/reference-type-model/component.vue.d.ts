type __VLS_ModelProps = {
    /**
     * required number modelValue
     */
    modelValue: number;
    /**
     * optional boolean foo with default false
     */
    'foo'?: boolean;
    /**
     * optional string bar with lazy and trim modifiers
     */
    'bar'?: string;
    'barModifiers'?: Partial<Record<'lazy' | 'trim', true>>;
};
type __VLS_ModelEmit = {
    'update:modelValue': [value: number];
    'update:foo': [value: boolean];
    'update:bar': [value: string | undefined];
};
declare const __VLS_export: import("vue").DefineComponent2<{
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
    __typeProps: __VLS_ModelProps;
    __typeEmits: __VLS_ModelEmit;
    __typeRefs: {};
    __typeEl: any;
    __defaults: unknown;
}>;
declare const _default: typeof __VLS_export;
export default _default;
