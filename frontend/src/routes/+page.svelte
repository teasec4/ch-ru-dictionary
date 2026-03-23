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
      meanings: string[];
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

<div class="min-h-screen bg-bg-primary">
    <!-- Hero Section -->
    <div class="relative overflow-hidden">
        <div class="absolute inset-0 opacity-5 bg-gradient-primary"></div>
        
        <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 sm:py-6 md:py-24">
            <div class="text-center">
                <h1 class="text-4xl md:text-5xl lg:text-6xl font-bold text-text-primary mb-6">
                    Китайский <span class="text-dict-1">словарь</span>
                </h1>
                <p class="text-xl md:text-2xl text-dict-2 mb-10 max-w-3xl mx-auto">
                    Поиск по иероглифам, пиньину и определениям
                </p>
                
                <!-- Поиск -->
                <div class="max-w-2xl mx-auto">
                    <div class="relative">
                        <div class="flex flex-col sm:flex-row gap-4 items-center justify-center">
                            <div class="relative flex-1 w-full sm:w-auto">
                                <input
                                    type="text"
                                    placeholder="Поиск..."
                                    bind:value={searchTerm}
                                    onkeydown={handleKeyPress}
                                    class="w-full px-6 py-4 text-lg rounded-xl border-2 border-dict-3 focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/20 shadow-md transition-all duration-200"
                                    disabled={loading}
                                />
                                {#if searchTerm}
                                    <button
                                        class="absolute right-3 top-1/2 transform -translate-y-1/2 text-dict-2 hover:text-dict-1"
                                        onclick={() => searchTerm = ""}
                                    >
                                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                                        </svg>
                                    </button>
                                {/if}

                            </div>
                            
                            <button
                                onclick={handleClick}
                                disabled={loading || !searchTerm.trim()}
                                class="px-8 py-4 rounded-xl bg-primary text-white font-semibold text-lg hover:bg-dict-1 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-lg hover:shadow-xl flex items-center gap-2 min-w-[140px] justify-center"
                            >
                                {#if loading}
                                    <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                    </svg>
                                    <span>Поиск...</span>
                                {:else}
                                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                                    </svg>
                                    <span>Найти</span>
                                {/if}
                            </button>
                            
                        </div>
                        
                        <div class="mt-4 text-sm text-dict-2 flex flex-wrap justify-center gap-3">
                            <span class="flex items-center gap-1">
                                <span class="w-2 h-2 rounded-full bg-success"></span>
                                Иероглифы
                            </span>
                            <span class="flex items-center gap-1">
                                <span class="w-2 h-2 rounded-full bg-accent"></span>
                                Пиньин
                            </span>
                            <span class="flex items-center gap-1">
                                <span class="w-2 h-2 rounded-full bg-warning"></span>
                                Перевод
                            </span>
                        </div>

                        
                    </div>
                </div>
                
                
            </div>
        </div>
    </div>
    
    <!-- Результаты поиска --> 
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 md:py-12">
        {#if responseData}
            <div class="mb-8 fade-in">
                <div class="flex items-center justify-between mb-6">
                    <h2 class="text-2xl md:text-3xl font-bold text-text-primary">
                        Найдено: <span class="text-accent">{responseData.count}</span> результатов
                    </h2>
                    <button
                        onclick={() => responseData = null}
                        class="px-4 py-2 rounded-lg border border-dict-3 text-dict-2 hover:bg-dict-4 hover:text-dict-1 transition-colors"
                    >
                        Очистить
                    </button>
                </div>

                <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                    {#each responseData.data as word, i}
                        <div class="dict-card p-6 fade-in" style="animation-delay: {i * 0.1}s">
                            <div class="flex items-start justify-between mb-4">
                                <div>
                                    <div class="chinese-text text-3xl font-bold text-dict-1 mb-1">
                                        {word.chinese}
                                    </div>
                                    <div class="pinyin-text text-lg text-accent font-medium">
                                        {word.pinyin}
                                    </div>
                                </div>
                                <div class="text-sm px-3 py-1 rounded-full bg-dict-4 text-dict-2">
                                    #{i + 1}
                                </div>
                            </div>

                            <div class="space-y-3">
                                <div class="text-sm font-semibold text-dict-2 uppercase tracking-wide">
                                    Значения:
                                </div>
                                <div class="space-y-2">
                                    {#each word.meanings as meaning}
                                        <p class="text-dict-1">{meaning}</p>
                                    {/each}
                                </div>
                            </div>

                            <div class="mt-6 pt-4 border-t border-dict-4 flex justify-end">
                                <button
                                    class="text-sm text-accent hover:text-primary font-medium flex items-center gap-1"
                                    onclick={() => {/* TODO: Добавить функционал копирования */}}
                                >
                                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
                                    </svg>
                                    Копировать
                                </button>
                            </div>
                        </div>
                    {/each}
                </div>

                {#if responseData.message === "no results"}
                    <div class="text-center py-12 fade-in">
                        <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-dict-4 mb-4">
                            <svg class="w-8 h-8 text-dict-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                            </svg>
                        </div>
                        <h3 class="text-xl font-semibold text-dict-1 mb-2">Ничего не найдено</h3>
                        <p class="text-dict-2 max-w-md mx-auto">
                            Попробуйте изменить запрос или использовать другой формат поиска
                        </p>
                    </div>
                {/if}

                
            </div>
        {:else if loading}
        <div class="text-center py-16">
            <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-dict-3 border-t-accent mb-4"></div>
            <p class="text-lg text-dict-2">Ищем в словаре...</p>
        </div>
        {:else if error}
        <div class="max-w-md mx-auto text-center py-12 fade-in">
            <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-error/10 mb-4">
                <svg class="w-8 h-8 text-error" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
            </div>
            <h3 class="text-xl font-semibold text-error mb-2">Произошла ошибка</h3>
            <p class="text-dict-2 mb-6">{error}</p>
            <button
                onclick={() => error = null}
                class="px-6 py-2 rounded-lg bg-primary text-white hover:bg-dict-1 transition-colors"
            >
                Попробовать снова
            </button>
        </div>

        {/if}
   </div>
    
    
</div>

<style lang="postcss">
    @reference "tailwindcss";

</style>