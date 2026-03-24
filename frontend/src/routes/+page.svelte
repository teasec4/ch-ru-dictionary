<script lang="ts">
    import { goto } from "$app/navigation";

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
        goto(`/search?q=${searchTerm}`);
        return;
      }
      
    }
</script>

<div class="w-full max-w-2xl text-center">
    <!-- Hero -->
    <div class="px-4">
        <!-- Hero -->
        <h1 class="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-bold text-text-primary mb-4">
                  Китайский <span class="text-dict-1">словарь</span>
              </h1>
              <p class="text-base sm:text-lg md:text-xl text-dict-2 mb-8 max-w-2xl mx-auto">
                  Поиск по иероглифам, пиньину и определениям
              </p>
            
        <!-- Search -->
        <div class="">
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

   
</div>

<style lang="postcss">
@reference "tailwindcss";
</style>