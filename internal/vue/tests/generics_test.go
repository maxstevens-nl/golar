package vue_tests

import (
	"testing"

	"github.com/microsoft/typescript-go/shim/fourslash"
	"github.com/microsoft/typescript-go/shim/lsp/lsproto"
	"github.com/microsoft/typescript-go/shim/testutil"
)

// Tests for Vue 3.3+ generic components.
// The `generic` attribute on <script setup> allows defining type parameters
// that can be used in props, emits, and throughout the component.
//
// See: https://vuejs.org/api/sfc-script-setup.html#generics

func TestGenericComponentBasic(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := withVueNodeModules(t, `// @filename: file.vue
<script setup lang="ts" generic="T">
	const props = defineProps<{
		items: T[]
		selected: T
	}>()
	const items/*1*/ = props.items
	const selected/*2*/ = props.selected
</script>

<template>
	<div v-for="item/*3*/ in items/*4*/">
		{{ item/*5*/ }}
	</div>
</template>`)
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.VerifyQuickInfoAt(t, "1", "const items: T[]", "")
	f.VerifyQuickInfoAt(t, "2", "const selected: T", "")
	f.VerifyQuickInfoAt(t, "3", "const item: T", "")
	f.VerifyQuickInfoAt(t, "4", "(property) items: T[]", "")
	f.VerifyQuickInfoAt(t, "5", "const item: T", "")

	f.VerifyNonSuggestionDiagnostics(t, []*lsproto.Diagnostic{})
}

func TestGenericComponentWithConstraint(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := withVueNodeModules(t, `// @filename: file.vue
<script setup lang="ts" generic="T extends { id: number }">
	const props = defineProps<{
		item: T
	}>()
	const id/*1*/ = props.item.id
</script>

<template>
	<div>{{ props.item.id/*2*/ }}</div>
</template>`)
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.VerifyQuickInfoAt(t, "1", "const id: number", "")
	f.VerifyQuickInfoAt(t, "2", "(property) id: number", "")

	f.VerifyNonSuggestionDiagnostics(t, []*lsproto.Diagnostic{})
}

func TestGenericComponentMultipleTypeParams(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := withVueNodeModules(t, `// @filename: file.vue
<script setup lang="ts" generic="K, V">
	const props = defineProps<{
		entries: Map<K, V>
	}>()
	const entries/*1*/ = props.entries
</script>

<template>
	<div>{{ entries/*2*/ }}</div>
</template>`)
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.VerifyQuickInfoAt(t, "1", "const entries: Map<K, V>", "")
	f.VerifyQuickInfoAt(t, "2", "(property) entries: Map<K, V>", "")

	f.VerifyNonSuggestionDiagnostics(t, []*lsproto.Diagnostic{})
}

func TestGenericComponentWithDefault(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := withVueNodeModules(t, `// @filename: file.vue
<script setup lang="ts" generic="T = string">
	const props = defineProps<{
		value: T
	}>()
	const value/*1*/ = props.value
</script>

<template>
	<div>{{ value/*2*/ }}</div>
</template>`)
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.VerifyQuickInfoAt(t, "1", "const value: T", "")
	f.VerifyQuickInfoAt(t, "2", "(property) value: T", "")

	f.VerifyNonSuggestionDiagnostics(t, []*lsproto.Diagnostic{})
}

func TestGenericComponentWithConstraintAndDefault(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := withVueNodeModules(t, `// @filename: file.vue
<script setup lang="ts" generic="T extends string | number = string">
	const props = defineProps<{
		value: T
	}>()
	const value/*1*/ = props.value
</script>

<template>
	<div>{{ value/*2*/ }}</div>
</template>`)
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.VerifyQuickInfoAt(t, "1", "const value: T", "")
	f.VerifyQuickInfoAt(t, "2", "(property) value: T", "")

	f.VerifyNonSuggestionDiagnostics(t, []*lsproto.Diagnostic{})
}

func TestGenericComponentEmit(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := withVueNodeModules(t, `// @filename: file.vue
<script setup lang="ts" generic="T">
	const props = defineProps<{
		items: T[]
	}>()
	const emit = defineEmits<{
		(e: 'select', item: T): void
	}>()
	function handleSelect(item/*1*/: T) {
		emit('select', item/*2*/)
	}
</script>

<template>
	<div v-for="item/*3*/ in props.items">
		{{ item/*4*/ }}
	</div>
</template>`)
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.VerifyQuickInfoAt(t, "1", "(parameter) item: T", "")
	f.VerifyQuickInfoAt(t, "2", "(parameter) item: T", "")
	f.VerifyQuickInfoAt(t, "3", "const item: T", "")
	f.VerifyQuickInfoAt(t, "4", "const item: T", "")

	f.VerifyNonSuggestionDiagnostics(t, []*lsproto.Diagnostic{})
}

func TestGenericComponentConstraintViolation(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := withVueNodeModules(t, `// @filename: file.vue
<script setup lang="ts" generic="T extends { id: number }">
	const props = defineProps<{
		item: T
	}>()
	// This should error - T might not have 'name' property
	const name/*1*/ = props.item.[|name|]
</script>`)
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.VerifyNonSuggestionDiagnostics(t, []*lsproto.Diagnostic{
		{
			Code:    &lsproto.IntegerOrString{Integer: ptrTo[int32](2339)},
			Message: "Property 'name' does not exist on type 'T'.",
		},
	})
}

func TestGenericComponentKeyofConstraint(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	content := withVueNodeModules(t, `// @filename: file.vue
<script setup lang="ts" generic="T, K extends keyof T">
	const props = defineProps<{
		obj: T
		key: K
	}>()
	const value/*1*/ = props.obj[props.key]
</script>

<template>
	<div>{{ value/*2*/ }}</div>
</template>`)
	f, done := fourslash.NewFourslash(t, nil, content)
	defer done()

	f.VerifyQuickInfoAt(t, "1", "const value: T[K]", "")
	f.VerifyQuickInfoAt(t, "2", "(property) value: T[K]", "")

	f.VerifyNonSuggestionDiagnostics(t, []*lsproto.Diagnostic{})
}
