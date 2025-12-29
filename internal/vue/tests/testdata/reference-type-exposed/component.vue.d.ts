import { type VNode } from 'vue';
import type { MyExposed } from './my-exposed';
type __VLS_Props = {
    label: string;
};
type __VLS_Emit = {
    click: [event: Event];
};
type __VLS_Slots = {
    default: () => VNode[];
};
declare const __VLS_base: import("vue").DefineComponent2<{
    setup(): MyExposed;
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
    __typeProps: __VLS_Props;
    __typeEmits: __VLS_Emit;
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
