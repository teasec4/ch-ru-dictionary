<script lang="ts">
    import { Loader2 } from "lucide-svelte";

    let { word, translations } = $props();
    
    let loading = $state(false);
    let results = $state<any[]>([]);
    let loaded = $state(false);
    let error = $state<string | null>(null);
    let clicked = $state(false);

    async function loadExtended() {
        if (loading || loaded) return;
        
        clicked = true;
        loading = true;
        error = null;
        
        try {
            const res = await fetch(`/api/search/${word}?mode=extended`);
            const data = await res.json();
            results = data.data || [];
            loaded = true;
        } catch (e) {
            error = translations?.error || "Error";
        } finally {
            loading = false;
        }
    }
</script>

{#if !clicked}
    <div class="text-center mt-6">
        <button 
            onclick={loadExtended}
            class="px-6 py-2 bg-dict-4 text-dict-1 rounded-lg hover:bg-dict-3 transition"
        >
            {translations?.showMore || "Показать где встречается"}
        </button>
    </div>
{:else if loaded && results.length > 0}
    <div class="mt-8">
        <h3 class="text-lg font-semibold text-text-primary mb-4">
            {translations?.whereAppears || "Где встречается"}:
        </h3>
        <div class="grid gap-3 sm:gap-4 md:grid-cols-2 lg:grid-cols-3">
            {#each results as item, i}
                <div class="p-3 sm:p-4 rounded-xl border border-dict-4 shadow-sm bg-white/30">
                    <div class="flex justify-between items-start">
                        <div>
                            <div class="text-xl sm:text-2xl font-bold text-dict-1">
                                {item.chinese}
                            </div>
                            <div class="text-sm text-accent">
                                {item.pinyin}
                            </div>
                        </div>
                    </div>
                    {#if item.meanings}
                        <p class="text-xs sm:text-sm text-dict-2 mt-2 line-clamp-2">
                            {item.meanings}
                        </p>
                    {/if}
                </div>
            {/each}
        </div>
    </div>
{:else if loading}
    <div class="text-center py-6">
        <Loader2 class="w-6 h-6 animate-spin mx-auto mb-2 text-accent" />
        <p class="text-sm text-dict-2">{translations?.loading || "Загрузка..."}</p>
    </div>
{:else if error}
    <div class="text-center py-6">
        <p class="text-sm text-red-500">{error}</p>
    </div>
{/if}