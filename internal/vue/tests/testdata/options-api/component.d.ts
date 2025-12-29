interface SubmitPayload {
    /**
     * email of user
     */
    email: string;
    /**
     * password of same user
     */
    password: string;
}
declare const _default: import("vue").DefineComponent2<{
    setup(): {};
    data(): {};
    props: {
        /**
         * Default number
         */
        numberDefault: {
            type: NumberConstructor;
            default: number;
        };
        /**
         * Default function Object
         */
        objectDefault: {
            type: ObjectConstructor;
            default: () => {
                foo: string;
            };
        };
        /**
         * Default function Array
         */
        arrayDefault: {
            type: ArrayConstructor;
            default: () => number[];
        };
        /**
         * Default function more complex
         */
        complexDefault: {
            type: ArrayConstructor;
            default: (props: any) => any[];
        };
    };
    computed: {};
    methods: {};
    mixins: {}[];
    extends: {};
    emits: {
        submit: ({ email, password }: SubmitPayload) => boolean;
    };
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
export default _default;
