<script lang="ts">
    import { goto } from "$app/navigation";
    import { navigating } from "$app/stores";
    import { Loader2 } from "lucide-svelte";
    import { page } from "$app/state";
    let { data } = $props();
    const translations = data.translations;

    let searchTerm = $state("");
    let error = $state<string | null>(null);


    async function handleSearch() {
        console.log('handleSearch called', searchTerm, $navigating);
        const lang = page.params.lang;
        if (!searchTerm.trim() || !!$navigating) return;

        error = null;
        goto(`/${lang}/search/${encodeURIComponent(searchTerm.trim())}`);
    }

    function handleKeyPress(event: KeyboardEvent) {
        if (event.key === "Enter") {
            handleSearch();
        }
    }

    // Удаляем старые переменные и кнопку
</script>

<div class="min-h-screen w-full flex flex-col">
    <!-- SEARCH BAR -->
    <div class="sticky top-12 bg-white z-40 px-4 py-3 shrink-0">
        <div class="max-w-2xl mx-auto flex gap-2">
            <div class="relative flex-1">
                <input
                    type="text"
                    placeholder={translations.searchPlaceholder}
                    bind:value={searchTerm}
                    onkeydown={handleKeyPress}
                    class="w-full px-4 sm:px-5 py-3 sm:py-4 text-base sm:text-lg rounded-xl border border-dict-3 focus:border-accent focus:ring-2 focus:ring-accent/20 outline-none shadow-sm transition"
                    disabled={!!$navigating}
                />
                {#if searchTerm.trim()}
                    <button
                        class="absolute right-5 top-1/2 -translate-y-1/2 w-6 h-6 flex items-center justify-center text-dict-2 hover:text-dict-1 text-sm"
                        onclick={() => (searchTerm = "")}
                        aria-label={translations.clearSearch || "Очистить"}
                    >
                        ×
                    </button>
                {/if}
            </div>
            <button
                disabled={!searchTerm.trim() || !!$navigating}
                class="px-6 sm:px-8 py-3 sm:py-4 rounded-xl bg-primary text-white font-semibold hover:bg-dict-1 active:scale-95 disabled:opacity-50 transition flex items-center justify-center gap-2"
                onclick={handleSearch}
            >
                {#if !!$navigating}
                    <Loader2 class="w-5 h-5 animate-spin" />
                {:else}
                    {translations.find}
                {/if}
            </button>
        </div>
    </div>

    <!-- Индикатор загрузки под поиском -->

    {#if error}
        <div class="max-w-2xl mx-auto mt-3">
            <div class="bg-red-50 border border-red-200 rounded-lg p-3">
                <p class="text-red-600 text-sm">
                    {translations.error}: {error}
                </p>
            </div>
        </div>
    {/if}

    <!-- Results -->
    <div class="max-w-6xl mx-auto pb-12 p-4 pt-20">
        {#if !!!$navigating && data}

            {#if data.data.count > 0}
                <div class="grid gap-4 sm:gap-5 md:grid-cols-2 lg:grid-cols-3">
                    {#each data.data.data as word, i}
                        <div
                            class="p-4 sm:p-5 rounded-2xl border border-dict-4 shadow-sm bg-white/50"
                        >
                            <div class="flex justify-between items-start mb-3">
                                <div>
                                    <div
                                        class="text-2xl sm:text-3xl font-bold text-dict-1"
                                    >
                                        {word.hanzi}
                                    </div>
                                    <div class="text-sm sm:text-base text-accent">
                                        {word.pinyin}
                                    </div>
                                </div>
                                <span class="text-xs px-2 py-1 rounded bg-dict-4"
                                    >#{i + 1}</span
                                >
                            </div>

                            <div class="space-y-2">
                                <div
                                    class="text-xs uppercase text-dict-2 tracking-wide"
                                >
                                    {translations.meanings}
                                </div>

                                <p
                                    class="text-dict-1 meaning-text text-sm sm:text-base leading-relaxed whitespace-normal"
                                >
                                    {#each word.meanings as meaning}
                                        <span class="meaning-line">
                                            <span class="meaning-index"
                                                >{meaning.index})</span
                                            >
                                            {meaning.text}
                                            {#if meaning.refs && meaning.refs.length > 0}
                                                <span class="refs-list ml-1">
                                                    (
                                                    {#each meaning.refs as ref, ri}
                                                        <a
                                                            href="/{page.params.lang}/search/{ref}"
                                                            class="text-accent hover:underline font-medium ref-link"
                                                            >{ref}</a
                                                        >{ri <
                                                            meaning.refs.length -
                                                                1
                                                            ? ", "
                                                            : ""}
                                                    {/each}
                                                    )
                                                </span>
                                            {/if}
                                        </span>
                                    {/each}
                                </p>
                            </div>
                        </div>
                    {/each}
                </div>
            {:else}
                <div class="text-center py-10">
                    <p class="text-dict-2">{translations.noResults}</p>
                </div>
            {/if}

        {:else if !!$navigating}
            <div class="text-center py-16">
                <Loader2 class="w-10 h-10 animate-spin mx-auto mb-4 text-accent" />
                <p class="text-dict-2">{translations.loading}...</p>
            </div>
        {:else if error}
            <div class="text-center py-12">
                <p class="text-error mb-4">{error}</p>
                <button
                    onclick={() => (error = null)}
                    class="px-4 py-2 bg-primary text-white rounded"
                >
                    {translations.retry}
                </button>
            </div>
        {/if}
    </div>
</div>

<style lang="postcss">
    @reference "tailwindcss";

    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }
    .animate-spin {
        animation: spin 1s linear infinite;
    }
</style>
