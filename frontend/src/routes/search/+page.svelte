<script lang="ts">
    import { page } from "$app/state";

    let searchTerm = $derived(page.url.searchParams.get('q') ?? '');
    let newSearchTerm = $state('');
    let responseData = $state<SearchResponse | null>(null);
    let loading = $state(false);
    let error = $state<string | null>(null);
    
    // for first load
    let initialized = false;
    
    $effect(() => {
      if (!initialized && searchTerm) {
        newSearchTerm = searchTerm;
        handleSearch();
        initialized = true;
      }
    });
    
    type Response = {
      chinese: string;
      pinyin: string;
      meanings: string;
    }
    
    type SearchResponse = {
      data: Response[];
      count: number;
      message?: string;
    }
    

    async function handleSearch() {
        console.log('Searching for:', newSearchTerm);
        if(newSearchTerm.trim()){
          loading = true;
          error = null;
          try {
            const response = await fetch("/api/search", {
              method: "POST",
              headers: {
                "Content-Type": "application/json",
              },
              body: JSON.stringify({ word: newSearchTerm }),
            });
            if (response.ok) {
              const data = await response.json();
              responseData = data;
            } else {
              responseData = null;
            }
          } catch (e) {
            error = e instanceof Error ? e.message : String(e);
          } finally {
            loading = false;
          }
        }
    }
    
    function handleKeyPress(event: KeyboardEvent){
      if(event.key === "Enter") {
        handleSearch();
      }
    }
</script>


<div class="min-h-screen flex flex-col gap-10">

	<!-- SEARCH BAR -->
	<div class="sticky top-14 bg-white z-40  px-4 py-3">
		<div class="max-w-2xl mx-auto flex gap-2">
		<input
                        type="text"
                        placeholder="Поиск..."
                        bind:value={newSearchTerm}
                        onkeydown={handleKeyPress}
                        class="w-full px-4 sm:px-5 py-3 sm:py-4 text-base sm:text-lg rounded-xl border border-dict-3 focus:border-accent focus:ring-2 focus:ring-accent/20 outline-none shadow-sm transition"
                        disabled={loading}
                    />
			<button class="px-6 sm:px-8 py-3 sm:py-4 rounded-xl bg-primary text-white font-semibold hover:bg-dict-1 active:scale-95 disabled:opacity-50 transition flex items-center justify-center gap-2" onclick={handleSearch}>Найти</button>
		</div>
	</div>
	

	<!-- Results -->
    <div class="max-w-6xl mx-auto pb-12 p-4">
        {#if responseData}
            <div class="mb-6 flex items-center justify-between flex-wrap gap-3">
                <h2 class="text-xl sm:text-2xl font-bold text-text-primary">
                    Найдено: <span class="text-accent">{responseData.count}</span>
                </h2>
                <button
                    onclick={() => responseData = null}
                    class="px-4 py-2 border rounded-lg text-sm text-dict-2 hover:bg-dict-4"
                >Очистить</button>
            </div>

            <div class="grid gap-4 sm:gap-5 md:grid-cols-2 lg:grid-cols-3">
                {#each responseData.data as word, i}
                    <div class="p-4 sm:p-5 rounded-2xl border border-dict-4 shadow-sm bg-white/50">
                        <div class="flex justify-between items-start mb-3">
                            <div>
                                <div class="text-2xl sm:text-3xl font-bold text-dict-1">
                                    {word.chinese}
                                </div>
                                <div class="text-sm sm:text-base text-accent">
                                    {word.pinyin}
                                </div>
                            </div>
                            <span class="text-xs px-2 py-1 rounded bg-dict-4">#{i+1}</span>
                        </div>

                        <div class="space-y-2">
                            <div class="text-xs uppercase text-dict-2 tracking-wide">Значения</div>

                            <p class="text-dict-1 meaning-text text-sm sm:text-base leading-relaxed whitespace-normal">
                                {word.meanings}
                            </p>
                        </div>
                    </div>
                {/each}
            </div>

            {#if responseData.message === "no results"}
                <div class="text-center py-10">
                    <p class="text-dict-2">Ничего не найдено</p>
                </div>
            {/if}

        {:else if loading}
            <div class="text-center py-16">
                <div class="animate-spin h-10 w-10 border-4 border-dict-3 border-t-accent rounded-full mx-auto mb-4"></div>
                <p class="text-dict-2">Ищем...</p>
            </div>

        {:else if error}
            <div class="text-center py-12">
                <p class="text-error mb-4">{error}</p>
                <button onclick={() => error = null} class="px-4 py-2 bg-primary text-white rounded">
                    Повторить
                </button>
            </div>
        {/if}
    </div>

</div>


<style lang="postcss">
@reference "tailwindcss";
</style>