export interface MyNestedProps {
    /**
     * nested prop documentation
     */
    nestedProp: string;
}
export interface MyIgnoredNestedProps {
    nestedProp: string;
}
export interface MyNestedRecursiveProps {
    recursive: MyNestedRecursiveProps;
}
declare enum MyEnum {
    Small = 0,
    Medium = 1,
    Large = 2
}
declare namespace MyNamespace {
    type MyType = {};
}
declare const categories: readonly ["Uncategorized", "Content", "Interaction", "Display", "Forms", "Addons"];
type MyCategories = typeof categories[number];
export interface MyProps {
    /**
     * string foo
     *
     * @default "rounded"
     * @since v1.0.0
     * @see https://vuejs.org/
     * @example
     * ```vue
     * <template>
     *   <component foo="straight" />
     * </template>
     * ```
     */
    foo: string;
    /**
     * optional number bar
     */
    bar?: number;
    /**
     * string array baz
     */
    baz?: string[];
    /**
     * required union type
     */
    union: string | number;
    /**
     * optional union type
     */
    unionOptional?: string | number;
    /**
     * required nested object
     */
    nested: MyNestedProps;
    /**
     * required nested object with intersection
     */
    nestedIntersection: MyNestedProps & {
        /**
         * required additional property
         */
        additionalProp: string;
    };
    /**
     * optional nested object
     */
    nestedOptional?: MyNestedProps | MyIgnoredNestedProps;
    /**
     * required array object
     */
    array: MyNestedProps[];
    /**
     * optional array object
     */
    arrayOptional?: MyNestedProps[];
    /**
     * enum value
     */
    enumValue: MyEnum;
    /**
     * namespace type
     */
    namespaceType: MyNamespace.MyType;
    /**
     * literal type alias that require context
     */
    literalFromContext: MyCategories;
    inlined: {
        foo: string;
    };
    recursive: MyNestedRecursiveProps;
}
export declare const StringRequired: {
    readonly type: StringConstructor;
    readonly required: true;
};
export declare const StringEmpty: {
    readonly type: StringConstructor;
    readonly default: "";
};
export declare const StringUndefined: {
    readonly type: StringConstructor;
};
export {};
