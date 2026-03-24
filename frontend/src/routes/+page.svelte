<script lang="ts">
    import { goto } from "$app/navigation";
    import { Search, Bookmark, Sparkles } from "lucide-svelte";
    let searchTerm = $state("");
  
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
                        
                    />
                    
                    <button
                        class="absolute right-3 top-1/2 -translate-y-1/2  px-6 sm:px-8rounded-xl bg-primary hover:bg-dict-1 text-white font-semibold rounded-xl py-2 active:scale-95 disabled:opacity-50 transition"
                        onclick={handleClick}
                        disabled={!searchTerm.trim()}
                    >Поиск
                    </button>
                </div>

            </div>
        </div>
        
        <div class="pt-8 flex justify-center">
          <ul class="space-y-3 text-dict-2 text-sm sm:text-base">
        
            <li class="flex items-center justify-center gap-3 border border-dict-4 rounded-2xl p-1">
              <Search class="w-5 h-5 text-accent" />
              <span>Быстрый поиск</span>
            </li>
        
            <li class="flex items-center justify-center gap-3 border border-dict-4 rounded-2xl p-1">
              <Bookmark class="w-5 h-5 text-accent" />
              <span>Сохранение слов</span>
            </li>
        
            <li class="flex items-center justify-center gap-3 border border-dict-4 rounded-2xl p-1">
              <Sparkles class="w-5 h-5 text-accent" />
              <span>AI генерация примеров</span>
            </li>
        
          </ul>
        </div>
        
    </div>

   
</div>

<style lang="postcss">
@reference "tailwindcss";
</style>