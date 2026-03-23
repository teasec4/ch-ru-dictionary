<script lang="ts">
    let searchTerm = $state("");
    let responseData = $state<SearchResponse | null>(null);
    let loading = $state(false);
    let error = $state<string | null>(null);
    
    $effect(() => {
        console.log(responseData);
    })
    
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
    
    function handleKeyPress(event: KeyboardEvent){
      if(event.key === "Enter") {
        handleClick();
      }
    }

    async function handleClick() {
      if(searchTerm.trim()){
        loading = true;
        error = null;
        try {
          const response = await fetch("/api/search", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ word: searchTerm }),
          });
          if (response.ok) {
            const data = await response.json();
            searchTerm = "";
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
</script>

<div class="min-h-screen px-4 sm:px-6 lg:px-8">
    <!-- Hero -->
    <div class="max-w-6xl mx-auto py-10 sm:py-14 md:py-20 text-center">
        <h1 class="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-bold text-text-primary mb-4">
            Китайский <span class="text-dict-1">словарь</span>
        </h1>
        <p class="text-base sm:text-lg md:text-xl text-dict-2 mb-8 max-w-2xl mx-auto">
            Поиск по иероглифам, пиньину и определениям
        </p>

        <!-- Search -->
        <div class="w-full max-w-2xl mx-auto">
            <div class="flex flex-col sm:flex-row gap-3">
                <div class="relative flex-1">
                    <input
                        type="text"
                        placeholder="Поиск..."
                        bind:value={searchTerm}
                        onkeydown={handleKeyPress}
                        class="w-full px-4 sm:px-5 py-3 sm:py-4 text-base sm:text-lg rounded-xl border border-dict-3 focus:border-accent focus:ring-2 focus:ring-accent/20 outline-none shadow-sm transition"
                        disabled={loading}
                    />

                    {#if searchTerm}
                        <button
                            class="absolute right-3 top-1/2 -translate-y-1/2 text-dict-2 hover:text-dict-1"
                            onclick={() => searchTerm = ""}
                        >✕</button>
                    {/if}
                </div>

                <button
                    onclick={handleClick}
                    disabled={loading || !searchTerm.trim()}
                    class="px-6 sm:px-8 py-3 sm:py-4 rounded-xl bg-primary text-white font-semibold hover:bg-dict-1 active:scale-95 disabled:opacity-50 transition flex items-center justify-center gap-2"
                >
                    {#if loading}
                        <span class="animate-spin h-5 w-5 border-2 border-white border-t-transparent rounded-full"></span>
                        Поиск...
                    {:else}
                        Найти
                    {/if}
                </button>
            </div>
        </div>
    </div>

    <!-- Results -->
    <div class="max-w-6xl mx-auto pb-12">
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