import type { MySlots } from './my-slots';
type __VLS_Slots = MySlots;
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
