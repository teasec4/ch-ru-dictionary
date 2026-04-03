<script lang="ts">
    import { goto } from "$app/navigation";
    import { page } from "$app/state";
    import { Search, Bookmark, Sparkles, Loader2 } from "lucide-svelte";
    let searchTerm = $state("");
    let loading = $state(false);

    let { data } = $props();
    const translations = data.translations;

    function handleKeyPress(event: KeyboardEvent) {
        if (event.key === "Enter") {
            handleSearch();
        }
    }

    async function handleSearch() {
        if (!searchTerm.trim() || loading) return;
        
        const lang = page.params.lang;
        loading = true;
        
        try {
            await goto(`/${lang}/search/${searchTerm}`);
        } finally {
            loading = false;
        }
    }
</script>

<div class="w-full max-w-2xl text-center">
    <!-- Hero -->
    <div class="px-4">
        <!-- Hero -->
        <h1
            class="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-bold text-text-primary mb-4"
        >
            {translations.heroTitle}
        </h1>
        <p
            class="text-base sm:text-lg md:text-xl text-dict-2 mb-8 max-w-2xl mx-auto"
        >
            {translations.heroSubtitle}
        </p>

        <!-- Search -->
        <div class="">
            <div class="flex flex-col sm:flex-row gap-3">
                <div class="relative flex-1">
                    <input
                        type="text"
                        placeholder={translations.searchPlaceholder}
                        bind:value={searchTerm}
                        onkeydown={handleKeyPress}
                        class="w-full px-4 sm:px-5 py-3 sm:py-4 text-base sm:text-lg rounded-xl border border-dict-3 focus:border-accent focus:ring-2 focus:ring-accent/20 outline-none shadow-sm transition"
                        disabled={loading}
                    />

                    <button
                        class="absolute right-3 top-1/2 -translate-y-1/2 px-6 sm:px-8 rounded-xl bg-primary hover:bg-dict-1 text-white font-semibold py-2 active:scale-95 disabled:opacity-50 transition flex items-center gap-2"
                        onclick={handleSearch}
                        disabled={!searchTerm.trim() || loading}
                    >
                        {#if loading}
                            <Loader2 class="w-5 h-5 animate-spin" />
                        {:else}
                            {translations.searchButton}
                        {/if}
                    </button>
                </div>
            </div>
        </div>

        <div class="pt-8 flex justify-center">
            <ul class="space-y-3 text-dict-2 text-sm sm:text-base">
                <li
                    class="flex items-center justify-center gap-3 border border-dict-4 rounded-2xl p-1"
                >
                    <Search class="w-5 h-5 text-accent" />
                    <span>{translations.fastSearch}</span>
                </li>

                <li
                    class="flex items-center justify-center gap-3 border border-dict-4 rounded-2xl p-1"
                >
                    <Bookmark class="w-5 h-5 text-accent" />
                    <span>{translations.saveWords}</span>
                </li>

                <li
                    class="flex items-center justify-center gap-3 border border-dict-4 rounded-2xl p-1"
                >
                    <Sparkles class="w-5 h-5 text-accent" />
                    <span>{translations.aiExamples}</span>
                </li>
            </ul>
        </div>
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
